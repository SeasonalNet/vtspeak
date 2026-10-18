set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-repeat-2006/legacy-neighbor-links-hardware-gdb.log
set logging overwrite on
set logging enabled on
set $legacy_dp_calls = 0

hbreak *0x1001cff0
commands
  silent
  set $legacy_dp_calls = $legacy_dp_calls + 1
  set $index = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  thbreak *$return
  commands
    silent
    set $context = $state + 0xbb154 + $index * 0xfc
    set $count = *(unsigned short *)($context + 0xf4)
    set $unit_ids = $context + 0x7c
    set $predecessors = $state + 0xaa7b4 + $index * 0x3c
    printf "LEGACY_NEIGHBOR_RETURN call=%u position=%u retained=%u\n", $legacy_dp_calls, $index, $count
    set $i = 0
    while $i < $count && $i < 30
      set $unit = *(unsigned int *)($unit_ids + $i * 4)
      set $pred_index = *(unsigned short *)($predecessors + $i * 2)
      if $index > 0
        set $previous_context = $state + 0xbb154 + ($index - 1) * 0xfc
        set $pred_unit = *(unsigned int *)($previous_context + 0x7c + $pred_index * 4)
        printf "LEGACY_NEIGHBOR_LINK position=%u candidate=%u predecessor_index=%u predecessor_unit=%u\n", $index, $unit, $pred_index, $pred_unit
      else
        printf "LEGACY_NEIGHBOR_LINK position=%u candidate=%u predecessor_index=%u\n", $index, $unit, $pred_index
      end
      set $i = $i + 1
    end
    continue
  end
  continue
end

hbreak *0x1001c860
commands
  silent
  set $legacy_score_unit = *(unsigned int *)($esp + 4)
  if ($legacy_score_unit == 272822) || ($legacy_score_unit == 273369) || ($legacy_score_unit == 272823) || ($legacy_score_unit == 273370)
    set $legacy_score_state = *(unsigned int *)($esp + 8)
    set $legacy_score_target = *(unsigned int *)($esp + 12)
    set $legacy_score_features = *(unsigned int *)($esp + 16)
    set $legacy_score_return = *(unsigned int *)$esp
    set $legacy_score_row = ($legacy_score_target - ($legacy_score_state + 0xb99e4)) / 6
    printf "LEGACY_REPEAT_SCORE_INPUT unit=%u row=%u target_context=%#x features=%#x\n", $legacy_score_unit, $legacy_score_row, $legacy_score_target, $legacy_score_features
    tbreak *$legacy_score_return
    commands
      silent
      printf "LEGACY_REPEAT_SCORE_RETURN unit=%u row=%u cost=%g\n", $legacy_score_unit, $legacy_score_row, $st0
      continue
    end
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "LEGACY_REPEAT_SELECTED_UNIT id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_NEIGHBOR_TRACE_READY\n"
  continue
end

continue
