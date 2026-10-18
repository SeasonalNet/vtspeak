set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-selector-byte0-old-candidate-scores-2026-09-29/apple-old-candidate-scores-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001c860
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if ($unit == 100180 || $unit == 59558 || $unit == 59559 || $unit == 64719 || $unit == 2554 || $unit == 2555 || $unit == 82036 || $unit == 82037)
    set $state = *(unsigned int *)($esp + 8)
    set $target = *(unsigned int *)($esp + 12)
    set $target_features = *(unsigned int *)($esp + 16)
    set $model = *(unsigned int *)($esp + 24)
    set $candidate_signature = *(unsigned int *)($model + 100) + $unit * 5
    set $return = *(unsigned int *)$esp
    set $row = ($target - ($state + 0xb99e4)) / 6
    printf "APPLE_OLD_CANDIDATE_SCORE_INPUT unit=%u target_row=%u target_context=%#x target_features=%#x signature:", $unit, $row, $target, $target_features
    x/5ub $candidate_signature
    thbreak *$return
    commands
      silent
      printf "APPLE_OLD_CANDIDATE_SCORE_RETURN unit=%u target_row=%u score=%g\n", $unit, $row, $st0
      continue
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_CANDIDATE_SCORE_TRACE_READY\n"
  continue
end

continue
