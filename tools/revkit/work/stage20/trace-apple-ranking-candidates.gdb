set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $selected_calls = 0
set $score_calls = 0

hbreak *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if $unit == 100180 || $unit == 103871 || $unit == 142230 || $unit == 128863 || $unit == 266884 || $unit == 266885 || $unit == 128862 || $unit == 142231 || $unit == 266883 || $unit == 128864 || $unit == 271296 || $unit == 96551
    set $score_calls = $score_calls + 1
    set $target = *(unsigned int *)($esp + 12)
    set $view = *(unsigned int *)($esp + 16)
    set $model = *(unsigned int *)($esp + 24)
    set $candidate = *(unsigned int *)($model + 0x64) + $unit * 7
    set $return = *(unsigned int *)$esp
    printf "APPLE_RANK_SCORE_ENTRY n=%u unit=%u target:", $score_calls, $unit
    x/7ub $target
    printf "APPLE_RANK_SCORE_VIEW:"
    x/10ub $view
    printf "APPLE_RANK_CANDIDATE_SIGNATURE:"
    x/7ub $candidate
    tbreak *$return
    commands
      silent
      printf "APPLE_RANK_SCORE_RETURN unit=%u score=%g\n", $unit, $st0
      continue
    end
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "APPLE_RANK_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_RANK_TRACE_READY\n"
  continue
end

continue
