set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-legacy-feature-crosswalk-candidate-path/adapted-query-selector-gdb.log
set logging overwrite on
set logging enabled on

set $query_calls = 0
set $split_calls = 0
set $fallback_calls = 0
set $selected_calls = 0
set $last_sum = -1

hbreak *0x10023dc0
commands
  silent
  set $query_calls = $query_calls + 1
  set $signature = *(unsigned int *)($ebp + 8)
  set $return = *(unsigned int *)($ebp + 4)
  set $key0 = *(unsigned char *)(0x1007b7ec + *(unsigned char *)($signature + 1))
  set $key1 = *(unsigned char *)(0x1007b788 + *(unsigned char *)($signature + 2))
  set $key2 = *(unsigned char *)(0x1007b850 + *(unsigned char *)($signature + 3))
  set $key3 = *(unsigned char *)($signature + 5)
  set $key4 = *(unsigned char *)($signature + 6) & 0x20
  tbreak *$return
  commands
    silent
    set $count = $eax & 0xffff
    printf "CROSSWALK_OLD_CUTOFF_QUERY call=%u key=%02x,%02x,%02x,%02x,%02x classes=%u\n", $query_calls, $key0, $key1, $key2, $key3, $key4, $count
    continue
  end
  continue
end

break *0x10024060
commands
  silent
  set $split_calls = $split_calls + 1
  set $last_sum = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "CROSSWALK_OLD_CUTOFF_SPLIT_ENTER call=%u slot=%u\n", $split_calls, $slot
  tbreak *$return
  commands
    silent
    set $native = $al
    if $native == 0 && $last_sum > 2
      printf "CROSSWALK_OLD_CUTOFF_OVERRIDE slot=%u sum=%d native=%u overridden=1\n", $slot, $last_sum, $native
      set $eax = 1
    else
      printf "CROSSWALK_OLD_CUTOFF_SPLIT_RETURN slot=%u sum=%d accepted=%u\n", $slot, $last_sum, $native
    end
    continue
  end
  continue
end

break *0x10023060
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $ids = *(unsigned int *)($esp + 4)
    set $count = *(unsigned short *)($esp + 8)
    set $context = *(unsigned int *)($esp + 12)
    set $weights = *(unsigned int *)($context + 0x8c)
    set $return = $caller
    printf "CROSSWALK_OLD_CUTOFF_METRIC_ENTRY count=%u ids_weights:", $count
    set $i = 0
    while $i < $count && $i < 96
      set $id = *(unsigned int *)($ids + $i * 4)
      printf " %u:%u", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      printf "CROSSWALK_OLD_CUTOFF_METRIC_RETURN sum=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x100242a0
commands
  silent
  set $fallback_calls = $fallback_calls + 1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "CROSSWALK_OLD_CUTOFF_FALLBACK_ENTER call=%u slot=%u\n", $fallback_calls, $slot
  tbreak *$return
  commands
    silent
    printf "CROSSWALK_OLD_CUTOFF_FALLBACK_RETURN slot=%u positions=%u\n", $slot, $eax
    continue
  end
  continue
end

break *0x10024680
commands
  silent
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    printf "CROSSWALK_OLD_CUTOFF_POSITION_BUILDER_RETURN positions=%u\n", $eax
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 32
    printf "CROSSWALK_OLD_CUTOFF_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  if $position < 4
    set $count = *(short *)($state + 0xae988 + $position * 0xfc)
    printf "CROSSWALK_FINAL_CANDIDATES position=%d count=%d", $position, $count
    set $index = 0
    while $index < $count && $index < 100
      set $candidate_node = *(unsigned int *)($state + 0x477ac + $index * 4)
      printf " %u", *(unsigned int *)($candidate_node + 8)
      set $index = $index + 1
    end
    printf "\n"
  end
  continue
end

break *0x1002389a
commands
  silent
  set $score_position = *(int *)($ebp + 8)
  set $score_node = *(unsigned int *)$esi
  set $score_unit = *(unsigned int *)($score_node + 8)
  if (($score_unit == 272822) || ($score_unit == 272823) || ($score_unit == 272824) || ($score_unit == 272825) || ($score_unit == 52542) || ($score_unit == 264072) || ($score_unit == 264073) || ($score_unit == 264074))
    printf "CROSSWALK_LOCAL_SCORE position=%d unit=%u base=%g span=%d weighted=%d local=%g\n", $score_position, $score_unit, *(float *)($score_node + 4), *(short *)($score_node + 0xe), *(short *)($score_node + 0x10), $st0
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "CROSSWALK_OLD_CUTOFF_TRACE_READY\n"
  continue
end

continue
