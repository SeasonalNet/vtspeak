set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-plain-crosswalk-old-candidate-scores/legacy-unit-score-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001c860
commands
  silent
  set $old_unit = *(unsigned int *)($esp + 4)
  set $new_unit = 0
  if $old_unit == 272822
    set $new_unit = 52542
  else
    if $old_unit == 272823
      set $new_unit = 273370
    else
      if $old_unit == 272824
        set $new_unit = 273371
      else
        if $old_unit == 272825
          set $new_unit = 255282
        end
      end
    end
  end
  if $new_unit != 0
    set $arg_address = $esp + 4
    set $state = *(unsigned int *)($esp + 8)
    set $target = *(unsigned int *)($esp + 12)
    set $target_features = *(unsigned int *)($esp + 16)
    set $model = *(unsigned int *)($esp + 24)
    set $candidate_signature = *(unsigned int *)($model + 100) + $new_unit * 5
    set $return = *(unsigned int *)$esp
    set $row = ($target - ($state + 0xb99e4)) / 6
    set {unsigned int}($arg_address) = $new_unit
    printf "OLD_SCORE_COUNTERFACTUAL_INPUT old_unit=%u scored_unit=%u target_row=%u target_context=%#x target_features=%#x signature:", $old_unit, $new_unit, $row, $target, $target_features
    x/5ub $candidate_signature
    printf "OLD_SCORE_COUNTERFACTUAL_TARGET_CONTEXT:"
    x/6ub $target
    printf "OLD_SCORE_COUNTERFACTUAL_TARGET_FEATURES:"
    x/8ub $target_features
    thbreak *$return
    commands
      silent
      printf "OLD_SCORE_COUNTERFACTUAL_RETURN old_unit=%u scored_unit=%u target_row=%u score=%g\n", $old_unit, $new_unit, $row, $st0
      set {unsigned int}($esp + 4) = $old_unit
      continue
    end
  end
  continue
end

break *0x408187
commands
  silent
  printf "OLD_SCORE_COUNTERFACTUAL_TRACE_READY\n"
  continue
end

continue
