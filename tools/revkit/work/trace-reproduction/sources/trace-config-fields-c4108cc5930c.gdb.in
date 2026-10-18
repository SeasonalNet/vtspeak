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
  disable 1
  set $out = (int *)malloc(20)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD baseline return=%d\n", $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-baseline.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(200, -1, -1, -1, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD pitch_max get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-pitch-max.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(50, -1, -1, -1, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD pitch_min get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-pitch-min.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-1, 400, -1, -1, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD speed_max get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-speed-max.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-1, 50, -1, -1, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD speed_min get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-speed-min.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-1, -1, 500, -1, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD volume_max get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-volume-max.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-1, -1, 0, -1, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD volume_min get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-volume-min.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-1, -1, -1, 65535, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD sent_pause_max get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-sentence-pause-max.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(-1, -1, -1, 0, 1)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD sent_pause_min get=%d pitch=%d speed=%d volume=%d pause=%d synth=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3), $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-sentence-pause-min.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int))0x100281b0)(200, 1)
  call ((void (*)(int, int))0x100281b0)(65535, 1)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($out, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD comma_max get=%d value=%d synth=%d\n", $comma_result, *$out, $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-comma-max.wav

  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, 1)
  call ((void (*)(int, int))0x100281b0)(200, 1)
  call ((void (*)(int, int))0x100281b0)(0, 1)
  set $comma_result = ((int (*)(int *, int))0x100281f0)($out, 1)
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, $texttype)
  printf "CONFIG_FIELD comma_min get=%d value=%d synth=%d\n", $comma_result, *$out, $result
  shell cp /work/stage5/output.wav /work/stage16/config-field-comma-min.wav
  continue
end

continue
