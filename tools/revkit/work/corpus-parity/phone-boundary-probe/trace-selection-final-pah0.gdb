set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /probe/score-signatures-gdb.log
set logging overwrite on
set logging enabled on
set $score_calls = 0
set $selected_calls = 0

break *0x100182e0
commands
  silent
  set $score_calls = $score_calls + 1
  if $score_calls <= 12
    set $unit = *(unsigned int *)($esp + 4)
    set $state = *(unsigned int *)($esp + 8)
    set $target = *(unsigned int *)($esp + 12)
    set $view = *(unsigned int *)($esp + 16)
    set $scale = *(float *)($esp + 20)
    set $model = *(unsigned int *)($esp + 24)
    set $expected = *(unsigned int *)($state + 0x4c) + *(unsigned char *)$target * 0x3c0 + *(unsigned char *)($target + 2) * 7 + 0x6e2
    set $candidate = *(unsigned int *)($model + 100) + $unit * 7
    set $ret = *(unsigned int *)$esp
    printf "SCORE_ENTRY call=%d unit=%u scale=%g target_desc=", $score_calls, $unit, $scale
    x/7bx $target
    printf "SCORE_EXPECTED_SIGNATURE="
    x/7bx $expected
    printf "SCORE_UNIT_SIGNATURE="
    x/7bx $candidate
    printf "SCORE_BYTE1 target=%u unit=%u category_target=%u category_unit=%u lookup_target=%u lookup_unit=%u\n", *(unsigned char *)($expected + 1), *(unsigned char *)($candidate + 1), *(unsigned char *)(0x1007b6c0 + *(unsigned char *)($expected + 1)), *(unsigned char *)(0x1007b6c0 + *(unsigned char *)($candidate + 1)), *(unsigned char *)(0x1007b97c + *(unsigned char *)($expected + 1)), *(unsigned char *)(0x1007b97c + *(unsigned char *)($candidate + 1))
    printf "SCORE_BYTE3 target=%u unit=%u category_target=%u category_unit=%u lookup_target=%u lookup_unit=%u\n", *(unsigned char *)($expected + 3), *(unsigned char *)($candidate + 3), *(unsigned char *)(0x1007b6c0 + *(unsigned char *)($expected + 3)), *(unsigned char *)(0x1007b6c0 + *(unsigned char *)($candidate + 3)), *(unsigned char *)(0x1007b97c + *(unsigned char *)($expected + 3)), *(unsigned char *)(0x1007b97c + *(unsigned char *)($candidate + 3))
    printf "SCORE_FEATURE_VIEW="
    x/10bx $view
    tbreak *$ret
    commands
      silent
      printf "SCORE_RETURN call=%d value=%g\n", $score_calls, $st0
      continue
    end
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

continue
exit
