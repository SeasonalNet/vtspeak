set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-legacy-feature-crosswalk-old-score/legacy-unit-score-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001c860
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if ($unit == 272822 || $unit == 272823 || $unit == 272824 || $unit == 272825 || $unit == 273369 || $unit == 273370 || $unit == 273371 || $unit == 255282 || $unit == 264072 || $unit == 264073 || $unit == 264074)
    set $state = *(unsigned int *)($esp + 8)
    set $target = *(unsigned int *)($esp + 12)
    set $target_features = *(unsigned int *)($esp + 16)
    set $model = *(unsigned int *)($esp + 24)
    set $candidate_signature = *(unsigned int *)($model + 100) + $unit * 5
    set $return = *(unsigned int *)$esp
    set $row = ($target - ($state + 0xb99e4)) / 6
    printf "LEGACY_SCORE_INPUT unit=%u target_row=%u target_context=%#x target_features=%#x signature:", $unit, $row, $target, $target_features
    x/5ub $candidate_signature
    printf "LEGACY_SCORE_TARGET_CONTEXT:"
    x/6ub $target
    printf "LEGACY_SCORE_TARGET_FEATURES:"
    x/8ub $target_features
    thbreak *$return
    commands
      silent
      printf "LEGACY_SCORE_RETURN unit=%u target_row=%u score=%g\n", $unit, $row, $st0
      continue
    end
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_SCORE_HARDWARE_TRACE_READY\n"
  continue
end

hbreak *0x1001ed7f
commands
  silent
  set $selected_index = *(short *)$edi
  printf "OLD_CROSSWALK_SELECTED index=%d unit=%u\n", $selected_index, $edx
  continue
end

continue
