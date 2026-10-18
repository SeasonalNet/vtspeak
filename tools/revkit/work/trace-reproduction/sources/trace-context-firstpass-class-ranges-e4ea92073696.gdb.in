set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $range_pending = 0

# Capture only the successful P/B synthesized signatures observed in the
# first-pass trace, then record their five-byte class keys and returned ranges.
break *0x10023de1
commands
  silent
  set $range_input = *(unsigned int *)($ebp + 8)
  set $range_match = 1
  if (int)*(unsigned char *)$range_input != 90
    set $range_match = 0
  end
  if (int)*(unsigned char *)($range_input + 1) != 57
    set $range_match = 0
  end
  if (int)*(unsigned char *)($range_input + 2) != 7
    set $range_match = 0
  end
  if (int)*(unsigned char *)($range_input + 3) != 53 && (int)*(unsigned char *)($range_input + 3) != 19 && (int)*(unsigned char *)($range_input + 3) != 32 && (int)*(unsigned char *)($range_input + 3) != 65
    set $range_match = 0
  end
  if (int)*(unsigned char *)($range_input + 4) != 7
    set $range_match = 0
  end
  if (int)*(unsigned char *)($range_input + 5) != 16 && (int)*(unsigned char *)($range_input + 5) != 32
    set $range_match = 0
  end
  if (int)*(unsigned char *)($range_input + 6) != 0
    set $range_match = 0
  end
  if $range_match
    set $range_pending = 1
    set $range_frame = $ebp
    printf "CLASS_QUERY signature="
    x/7ub $range_input
    printf "CLASS_KEY key="
    x/5ub $ebp-0xc
  end
  continue
end

break *0x10023dfb
commands
  silent
  if $range_pending
  if (int)$ebp == (int)$range_frame
    set $range_table = *(unsigned int *)($ebp + 0x14) + 0x88
    set $range_keys = *(unsigned int *)($range_table + 8)
    set $range_first = *(unsigned int *)($ebp + 8)
    printf "CLASS_RANGE count=%d first_index=%u key_count=%u\n", $eax, $range_first, *(unsigned int *)$range_table
    set $range_i = 0
    while $range_i < $eax
      printf "CLASS_RECORD index=%u bytes=", $range_first + $range_i
      x/5ub $range_keys + ($range_first + $range_i) * 5
      set $range_i = $range_i + 1
    end
    set $range_pending = 0
  end
  end
  continue
end

continue
