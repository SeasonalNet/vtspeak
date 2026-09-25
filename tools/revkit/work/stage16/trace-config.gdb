set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $fmt = *(int *)($esp + 4)
  set $text = *(unsigned int *)($esp + 8)
  set $path = *(unsigned int *)($esp + 12)
  set $speaker = *(int *)($esp + 16)
  set $pitch = *(int *)($esp + 20)
  set $speed = *(int *)($esp + 24)
  set $volume = *(int *)($esp + 28)
  set $pause = *(int *)($esp + 32)
  set $dict = *(int *)($esp + 36)
  set $texttype = *(int *)($esp + 40)
  printf "CONFIG_FILE_ENTRY fmt=%d speaker=%d pitch=%d speed=%d volume=%d pause=%d\n", $fmt, $speaker, $pitch, $speed, $volume, $pause
  disable 1
  set $out = (int *)malloc(20)
  set {int}($out) = 0
  set {int}($out + 1) = 0
  set {int}($out + 2) = 0
  set {int}($out + 3) = 0
  set {int}($out + 4) = 0
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  printf "CONFIG_INITIAL_GET ret=%d pitch=%d speed=%d volume=%d pause=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($out + 4, 1)
  printf "CONFIG_INITIAL_COMMA ret=%d value=%d\n", $comma_result, *($out + 4)
  set $baseline_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FILE_BASELINE_RETURN ret=%d\n", $baseline_result
  shell cp /work/stage5/output.wav /work/stage16/config-baseline.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(999, 999, 999, 70000, 1)
  call ((void (*)(int, int))0x100281b0)(70000, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  printf "CONFIG_CLAMPED_GET ret=%d pitch=%d speed=%d volume=%d pause=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($out + 4, 1)
  printf "CONFIG_CLAMPED_COMMA ret=%d value=%d\n", $comma_result, *($out + 4)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 0)
  printf "CONFIG_SLOT0_GET ret=%d pitch=%d speed=%d volume=%d pause=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($out + 4, 0)
  printf "CONFIG_SLOT0_COMMA ret=%d value=%d\n", $comma_result, *($out + 4)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FILE_RETURN ret=%d\n", $result
  shell cp /work/stage5/output.wav /work/stage16/config-output.wav
  continue
end

continue
