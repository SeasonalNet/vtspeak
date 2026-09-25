set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $half_pending = 0

# FUN_10019570 is the half-key lookup path. Keep only the controlled
# context-1 signatures used by the P/B/F/V/K/G left/right-phone fixtures.
break *0x10019570
commands
  silent
  set $range_out = *(unsigned int *)($esp + 4)
  set $label = *(unsigned int *)($esp + 8)
  set $half = *(unsigned int *)($esp + 12)
  set $query_match = 1
  if (int)*(unsigned char *)$label != 90 && (int)*(unsigned char *)$label != 7
    set $query_match = 0
  end
  if (int)*(unsigned char *)($label + 1) != 57 && (int)*(unsigned char *)($label + 1) != 53 && (int)*(unsigned char *)($label + 1) != 19
    set $query_match = 0
  end
  if (int)*(unsigned char *)($label + 2) != 7
    set $query_match = 0
  end
  if (int)*(unsigned char *)($label + 3) != 53 && (int)*(unsigned char *)($label + 3) != 19 && (int)*(unsigned char *)($label + 3) != 32 && (int)*(unsigned char *)($label + 3) != 65 && (int)*(unsigned char *)($label + 3) != 42 && (int)*(unsigned char *)($label + 3) != 33 && (int)*(unsigned char *)($label + 3) != 57
    set $query_match = 0
  end
  if (int)*(unsigned char *)($label + 4) != 7
    set $query_match = 0
  end
  if (int)*(unsigned char *)($label + 5) != 97 && (int)*(unsigned char *)($label + 5) != 81
    set $query_match = 0
  end
  if (int)*(unsigned char *)($label + 6) != 0
    set $query_match = 0
  end
  if $query_match && ($half == 1 || $half == 2)
    set $half_pending = 1
    set $pending_half = $half
    set $pending_label = $label
    set $pending_range_out = $range_out
    set $return_address = *(unsigned int *)$esp
    printf "HALF_QUERY_ENTRY half=%d label=", $half
    x/7ub $label
    tbreak *$return_address
    commands
      silent
      if $half_pending
        printf "HALF_QUERY_RETURN half=%d range_words=", $pending_half
        x/2uw $pending_range_out
        set $half_pending = 0
      end
      continue
    end
  end
  continue
end

# FUN_10019570 expands the five-byte class key to a ten-byte key immediately
# before calling this range search. The pending flag pairs the key with the
# target label and left/right half recorded above.
break *0x1002df50
commands
  silent
  if $half_pending
    set $half_key = *(unsigned int *)($esp + 8)
    printf "HALF_KEY_SEARCH half=%d key=", $pending_half
    x/10ub $half_key
  end
  continue
end

# FUN_10023c70 ranks the class range for one half target. Record the resulting
# class list and its model populations so each lookup is tied to candidates.
break *0x10023c70
commands
  silent
  set $candidate_list = *(unsigned int *)($esp + 12)
  set $candidate_label = *(unsigned int *)($esp + 16)
  set $candidate_half = *(unsigned short *)($esp + 20)
  set $candidate_model = *(unsigned int *)($esp + 32)
  set $candidate_match = 1
  if $candidate_half != 1 && $candidate_half != 2
    set $candidate_match = 0
  end
  if ((int)*(unsigned char *)$candidate_label != 90 && (int)*(unsigned char *)$candidate_label != 7) || ((int)*(unsigned char *)($candidate_label + 1) != 57 && (int)*(unsigned char *)($candidate_label + 1) != 53 && (int)*(unsigned char *)($candidate_label + 1) != 19) || (int)*(unsigned char *)($candidate_label + 2) != 7 || (int)*(unsigned char *)($candidate_label + 4) != 7 || ((int)*(unsigned char *)($candidate_label + 5) != 97 && (int)*(unsigned char *)($candidate_label + 5) != 81) || (int)*(unsigned char *)($candidate_label + 6) != 0
    set $candidate_match = 0
  end
  if (int)*(unsigned char *)($candidate_label + 3) != 53 && (int)*(unsigned char *)($candidate_label + 3) != 19 && (int)*(unsigned char *)($candidate_label + 3) != 32 && (int)*(unsigned char *)($candidate_label + 3) != 65 && (int)*(unsigned char *)($candidate_label + 3) != 42 && (int)*(unsigned char *)($candidate_label + 3) != 33 && (int)*(unsigned char *)($candidate_label + 3) != 57
    set $candidate_match = 0
  end
  if $candidate_match
    set $candidate_ret = *(unsigned int *)$esp
    set $candidate_half_saved = $candidate_half
    set $candidate_list_saved = $candidate_list
    set $candidate_model_saved = $candidate_model
    set $candidate_label_saved = $candidate_label
    tbreak *$candidate_ret
    commands
      silent
      set $candidate_count = $eax & 0xffff
      set $candidate_total_units = 0
      set $candidate_index = 0
      printf "HALF_CANDIDATES half=%d label=", $candidate_half_saved
      x/7ub $candidate_label_saved
      printf " count=%d ids=", $candidate_count
      while $candidate_index < $candidate_count && $candidate_index < 30
        set $class_id = *(unsigned int *)($candidate_list_saved + $candidate_index * 4)
        set $class_units = *(unsigned short *)(*(unsigned int *)($candidate_model_saved + 0x8c) + $class_id * 2)
        printf "%u(%u) ", $class_id, $class_units
        set $candidate_total_units = $candidate_total_units + $class_units
        set $candidate_index = $candidate_index + 1
      end
      printf "total_units=%d\n", $candidate_total_units
      continue
    end
  end
  continue
end

continue
