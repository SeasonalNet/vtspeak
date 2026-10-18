set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-continuity-bit7/old/apple-old-unit-scores-hardware-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001c860
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if ($unit == 100180 || $unit == 128863 || $unit == 177774 || $unit == 142230)
    set $state = *(unsigned int *)($esp + 8)
    set $target = *(unsigned int *)($esp + 12)
    set $target_features = *(unsigned int *)($esp + 16)
    set $model = *(unsigned int *)($esp + 24)
    set $candidate_signature = *(unsigned int *)($model + 100) + $unit * 5
    set $return = *(unsigned int *)$esp
    set $row = ($target - ($state + 0xb99e4)) / 6
    printf "APPLE_OLD_SCORE_INPUT unit=%u target_row=%u target_context=%#x target_features=%#x signature:", $unit, $row, $target, $target_features
    x/5ub $candidate_signature
    printf "APPLE_OLD_SCORE_TARGET_CONTEXT:"
    x/6ub $target
    printf "APPLE_OLD_SCORE_TARGET_FEATURES:"
    x/8ub $target_features
    thbreak *$return
    commands
      silent
      printf "APPLE_OLD_SCORE_RETURN unit=%u target_row=%u score=%g\n", $unit, $row, $st0
      continue
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_SCORE_TRACE_READY\n"
  continue
end

continue
