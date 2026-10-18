set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $speaker = *(int *)($esp + 16)
  disable 1
  set $out = (int *)malloc(16)
  set $comma_out = (int *)malloc(4)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, 345, 456, $speaker)
  call ((void (*)(int, int))0x100281b0)(567, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE stage=baseline getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-2147483648, 234, 345, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=pitch input=-2147483648 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-2, 234, 345, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=pitch input=-2 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-1, 234, 345, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=pitch input=-1 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, -2147483648, 345, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=speed input=-2147483648 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, -2, 345, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=speed input=-2 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, -1, 345, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=speed input=-1 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, -2147483648, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=volume input=-2147483648 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, -2, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=volume input=-2 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, -1, 456, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=volume input=-1 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, 345, -2147483648, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=sentence input=-2147483648 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, 345, -2, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=sentence input=-2 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, 345, -1, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=sentence input=-1 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int))0x100281b0)(-2147483648, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=comma input=-2147483648 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int))0x100281b0)(-2, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=comma input=-2 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int))0x100281b0)(-1, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($comma_out, $speaker)
  printf "NEG_PAUSE field=comma input=-1 getter_ret=%d pitch=%d speed=%d volume=%d sentence=%d comma_ret=%d comma=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $comma_result, *$comma_out
  call ((void (*)(int, int, int, int, int))0x10027fe0)(123, 234, 345, 456, $speaker)
  call ((void (*)(int, int))0x100281b0)(567, $speaker)
  kill
  quit
end

continue
