set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $fmt = *(int *)($esp + 4)
  set $path = *(char **)($esp + 12)
  set $speaker = *(int *)($esp + 16)
  disable 1
  set $out = (int *)malloc(4)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(0, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=0 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_0_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_0_space + 0) = 72
  set {char}($text_comma_0_space + 1) = 101
  set {char}($text_comma_0_space + 2) = 108
  set {char}($text_comma_0_space + 3) = 108
  set {char}($text_comma_0_space + 4) = 111
  set {char}($text_comma_0_space + 5) = 44
  set {char}($text_comma_0_space + 6) = 32
  set {char}($text_comma_0_space + 7) = 119
  set {char}($text_comma_0_space + 8) = 111
  set {char}($text_comma_0_space + 9) = 114
  set {char}($text_comma_0_space + 10) = 108
  set {char}($text_comma_0_space + 11) = 100
  set {char}($text_comma_0_space + 12) = 46
  set {char}($text_comma_0_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_0_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=0 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-0-space.wav
  set $text_comma_0_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_0_adjacent + 0) = 72
  set {char}($text_comma_0_adjacent + 1) = 101
  set {char}($text_comma_0_adjacent + 2) = 108
  set {char}($text_comma_0_adjacent + 3) = 108
  set {char}($text_comma_0_adjacent + 4) = 111
  set {char}($text_comma_0_adjacent + 5) = 44
  set {char}($text_comma_0_adjacent + 6) = 119
  set {char}($text_comma_0_adjacent + 7) = 111
  set {char}($text_comma_0_adjacent + 8) = 114
  set {char}($text_comma_0_adjacent + 9) = 108
  set {char}($text_comma_0_adjacent + 10) = 100
  set {char}($text_comma_0_adjacent + 11) = 46
  set {char}($text_comma_0_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_0_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=0 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-0-adjacent.wav
  set $text_comma_0_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_0_double + 0) = 72
  set {char}($text_comma_0_double + 1) = 101
  set {char}($text_comma_0_double + 2) = 108
  set {char}($text_comma_0_double + 3) = 108
  set {char}($text_comma_0_double + 4) = 111
  set {char}($text_comma_0_double + 5) = 44
  set {char}($text_comma_0_double + 6) = 32
  set {char}($text_comma_0_double + 7) = 32
  set {char}($text_comma_0_double + 8) = 119
  set {char}($text_comma_0_double + 9) = 111
  set {char}($text_comma_0_double + 10) = 114
  set {char}($text_comma_0_double + 11) = 108
  set {char}($text_comma_0_double + 12) = 100
  set {char}($text_comma_0_double + 13) = 46
  set {char}($text_comma_0_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_0_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=0 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-0-double.wav
  set $text_comma_0_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_0_tab + 0) = 72
  set {char}($text_comma_0_tab + 1) = 101
  set {char}($text_comma_0_tab + 2) = 108
  set {char}($text_comma_0_tab + 3) = 108
  set {char}($text_comma_0_tab + 4) = 111
  set {char}($text_comma_0_tab + 5) = 44
  set {char}($text_comma_0_tab + 6) = 9
  set {char}($text_comma_0_tab + 7) = 119
  set {char}($text_comma_0_tab + 8) = 111
  set {char}($text_comma_0_tab + 9) = 114
  set {char}($text_comma_0_tab + 10) = 108
  set {char}($text_comma_0_tab + 11) = 100
  set {char}($text_comma_0_tab + 12) = 46
  set {char}($text_comma_0_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_0_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=0 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-0-tab.wav
  set $text_comma_0_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_0_no_comma + 0) = 72
  set {char}($text_comma_0_no_comma + 1) = 101
  set {char}($text_comma_0_no_comma + 2) = 108
  set {char}($text_comma_0_no_comma + 3) = 108
  set {char}($text_comma_0_no_comma + 4) = 111
  set {char}($text_comma_0_no_comma + 5) = 32
  set {char}($text_comma_0_no_comma + 6) = 119
  set {char}($text_comma_0_no_comma + 7) = 111
  set {char}($text_comma_0_no_comma + 8) = 114
  set {char}($text_comma_0_no_comma + 9) = 108
  set {char}($text_comma_0_no_comma + 10) = 100
  set {char}($text_comma_0_no_comma + 11) = 46
  set {char}($text_comma_0_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_0_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=0 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-0-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(1, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=1 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_1_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_1_space + 0) = 72
  set {char}($text_comma_1_space + 1) = 101
  set {char}($text_comma_1_space + 2) = 108
  set {char}($text_comma_1_space + 3) = 108
  set {char}($text_comma_1_space + 4) = 111
  set {char}($text_comma_1_space + 5) = 44
  set {char}($text_comma_1_space + 6) = 32
  set {char}($text_comma_1_space + 7) = 119
  set {char}($text_comma_1_space + 8) = 111
  set {char}($text_comma_1_space + 9) = 114
  set {char}($text_comma_1_space + 10) = 108
  set {char}($text_comma_1_space + 11) = 100
  set {char}($text_comma_1_space + 12) = 46
  set {char}($text_comma_1_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_1_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=1 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-1-space.wav
  set $text_comma_1_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_1_adjacent + 0) = 72
  set {char}($text_comma_1_adjacent + 1) = 101
  set {char}($text_comma_1_adjacent + 2) = 108
  set {char}($text_comma_1_adjacent + 3) = 108
  set {char}($text_comma_1_adjacent + 4) = 111
  set {char}($text_comma_1_adjacent + 5) = 44
  set {char}($text_comma_1_adjacent + 6) = 119
  set {char}($text_comma_1_adjacent + 7) = 111
  set {char}($text_comma_1_adjacent + 8) = 114
  set {char}($text_comma_1_adjacent + 9) = 108
  set {char}($text_comma_1_adjacent + 10) = 100
  set {char}($text_comma_1_adjacent + 11) = 46
  set {char}($text_comma_1_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_1_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=1 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-1-adjacent.wav
  set $text_comma_1_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_1_double + 0) = 72
  set {char}($text_comma_1_double + 1) = 101
  set {char}($text_comma_1_double + 2) = 108
  set {char}($text_comma_1_double + 3) = 108
  set {char}($text_comma_1_double + 4) = 111
  set {char}($text_comma_1_double + 5) = 44
  set {char}($text_comma_1_double + 6) = 32
  set {char}($text_comma_1_double + 7) = 32
  set {char}($text_comma_1_double + 8) = 119
  set {char}($text_comma_1_double + 9) = 111
  set {char}($text_comma_1_double + 10) = 114
  set {char}($text_comma_1_double + 11) = 108
  set {char}($text_comma_1_double + 12) = 100
  set {char}($text_comma_1_double + 13) = 46
  set {char}($text_comma_1_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_1_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=1 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-1-double.wav
  set $text_comma_1_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_1_tab + 0) = 72
  set {char}($text_comma_1_tab + 1) = 101
  set {char}($text_comma_1_tab + 2) = 108
  set {char}($text_comma_1_tab + 3) = 108
  set {char}($text_comma_1_tab + 4) = 111
  set {char}($text_comma_1_tab + 5) = 44
  set {char}($text_comma_1_tab + 6) = 9
  set {char}($text_comma_1_tab + 7) = 119
  set {char}($text_comma_1_tab + 8) = 111
  set {char}($text_comma_1_tab + 9) = 114
  set {char}($text_comma_1_tab + 10) = 108
  set {char}($text_comma_1_tab + 11) = 100
  set {char}($text_comma_1_tab + 12) = 46
  set {char}($text_comma_1_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_1_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=1 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-1-tab.wav
  set $text_comma_1_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_1_no_comma + 0) = 72
  set {char}($text_comma_1_no_comma + 1) = 101
  set {char}($text_comma_1_no_comma + 2) = 108
  set {char}($text_comma_1_no_comma + 3) = 108
  set {char}($text_comma_1_no_comma + 4) = 111
  set {char}($text_comma_1_no_comma + 5) = 32
  set {char}($text_comma_1_no_comma + 6) = 119
  set {char}($text_comma_1_no_comma + 7) = 111
  set {char}($text_comma_1_no_comma + 8) = 114
  set {char}($text_comma_1_no_comma + 9) = 108
  set {char}($text_comma_1_no_comma + 10) = 100
  set {char}($text_comma_1_no_comma + 11) = 46
  set {char}($text_comma_1_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_1_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=1 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-1-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(199, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=199 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_199_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_199_space + 0) = 72
  set {char}($text_comma_199_space + 1) = 101
  set {char}($text_comma_199_space + 2) = 108
  set {char}($text_comma_199_space + 3) = 108
  set {char}($text_comma_199_space + 4) = 111
  set {char}($text_comma_199_space + 5) = 44
  set {char}($text_comma_199_space + 6) = 32
  set {char}($text_comma_199_space + 7) = 119
  set {char}($text_comma_199_space + 8) = 111
  set {char}($text_comma_199_space + 9) = 114
  set {char}($text_comma_199_space + 10) = 108
  set {char}($text_comma_199_space + 11) = 100
  set {char}($text_comma_199_space + 12) = 46
  set {char}($text_comma_199_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_199_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=199 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-199-space.wav
  set $text_comma_199_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_199_adjacent + 0) = 72
  set {char}($text_comma_199_adjacent + 1) = 101
  set {char}($text_comma_199_adjacent + 2) = 108
  set {char}($text_comma_199_adjacent + 3) = 108
  set {char}($text_comma_199_adjacent + 4) = 111
  set {char}($text_comma_199_adjacent + 5) = 44
  set {char}($text_comma_199_adjacent + 6) = 119
  set {char}($text_comma_199_adjacent + 7) = 111
  set {char}($text_comma_199_adjacent + 8) = 114
  set {char}($text_comma_199_adjacent + 9) = 108
  set {char}($text_comma_199_adjacent + 10) = 100
  set {char}($text_comma_199_adjacent + 11) = 46
  set {char}($text_comma_199_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_199_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=199 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-199-adjacent.wav
  set $text_comma_199_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_199_double + 0) = 72
  set {char}($text_comma_199_double + 1) = 101
  set {char}($text_comma_199_double + 2) = 108
  set {char}($text_comma_199_double + 3) = 108
  set {char}($text_comma_199_double + 4) = 111
  set {char}($text_comma_199_double + 5) = 44
  set {char}($text_comma_199_double + 6) = 32
  set {char}($text_comma_199_double + 7) = 32
  set {char}($text_comma_199_double + 8) = 119
  set {char}($text_comma_199_double + 9) = 111
  set {char}($text_comma_199_double + 10) = 114
  set {char}($text_comma_199_double + 11) = 108
  set {char}($text_comma_199_double + 12) = 100
  set {char}($text_comma_199_double + 13) = 46
  set {char}($text_comma_199_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_199_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=199 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-199-double.wav
  set $text_comma_199_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_199_tab + 0) = 72
  set {char}($text_comma_199_tab + 1) = 101
  set {char}($text_comma_199_tab + 2) = 108
  set {char}($text_comma_199_tab + 3) = 108
  set {char}($text_comma_199_tab + 4) = 111
  set {char}($text_comma_199_tab + 5) = 44
  set {char}($text_comma_199_tab + 6) = 9
  set {char}($text_comma_199_tab + 7) = 119
  set {char}($text_comma_199_tab + 8) = 111
  set {char}($text_comma_199_tab + 9) = 114
  set {char}($text_comma_199_tab + 10) = 108
  set {char}($text_comma_199_tab + 11) = 100
  set {char}($text_comma_199_tab + 12) = 46
  set {char}($text_comma_199_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_199_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=199 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-199-tab.wav
  set $text_comma_199_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_199_no_comma + 0) = 72
  set {char}($text_comma_199_no_comma + 1) = 101
  set {char}($text_comma_199_no_comma + 2) = 108
  set {char}($text_comma_199_no_comma + 3) = 108
  set {char}($text_comma_199_no_comma + 4) = 111
  set {char}($text_comma_199_no_comma + 5) = 32
  set {char}($text_comma_199_no_comma + 6) = 119
  set {char}($text_comma_199_no_comma + 7) = 111
  set {char}($text_comma_199_no_comma + 8) = 114
  set {char}($text_comma_199_no_comma + 9) = 108
  set {char}($text_comma_199_no_comma + 10) = 100
  set {char}($text_comma_199_no_comma + 11) = 46
  set {char}($text_comma_199_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_199_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=199 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-199-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=200 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_200_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_200_space + 0) = 72
  set {char}($text_comma_200_space + 1) = 101
  set {char}($text_comma_200_space + 2) = 108
  set {char}($text_comma_200_space + 3) = 108
  set {char}($text_comma_200_space + 4) = 111
  set {char}($text_comma_200_space + 5) = 44
  set {char}($text_comma_200_space + 6) = 32
  set {char}($text_comma_200_space + 7) = 119
  set {char}($text_comma_200_space + 8) = 111
  set {char}($text_comma_200_space + 9) = 114
  set {char}($text_comma_200_space + 10) = 108
  set {char}($text_comma_200_space + 11) = 100
  set {char}($text_comma_200_space + 12) = 46
  set {char}($text_comma_200_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_200_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=200 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-200-space.wav
  set $text_comma_200_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_200_adjacent + 0) = 72
  set {char}($text_comma_200_adjacent + 1) = 101
  set {char}($text_comma_200_adjacent + 2) = 108
  set {char}($text_comma_200_adjacent + 3) = 108
  set {char}($text_comma_200_adjacent + 4) = 111
  set {char}($text_comma_200_adjacent + 5) = 44
  set {char}($text_comma_200_adjacent + 6) = 119
  set {char}($text_comma_200_adjacent + 7) = 111
  set {char}($text_comma_200_adjacent + 8) = 114
  set {char}($text_comma_200_adjacent + 9) = 108
  set {char}($text_comma_200_adjacent + 10) = 100
  set {char}($text_comma_200_adjacent + 11) = 46
  set {char}($text_comma_200_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_200_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=200 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-200-adjacent.wav
  set $text_comma_200_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_200_double + 0) = 72
  set {char}($text_comma_200_double + 1) = 101
  set {char}($text_comma_200_double + 2) = 108
  set {char}($text_comma_200_double + 3) = 108
  set {char}($text_comma_200_double + 4) = 111
  set {char}($text_comma_200_double + 5) = 44
  set {char}($text_comma_200_double + 6) = 32
  set {char}($text_comma_200_double + 7) = 32
  set {char}($text_comma_200_double + 8) = 119
  set {char}($text_comma_200_double + 9) = 111
  set {char}($text_comma_200_double + 10) = 114
  set {char}($text_comma_200_double + 11) = 108
  set {char}($text_comma_200_double + 12) = 100
  set {char}($text_comma_200_double + 13) = 46
  set {char}($text_comma_200_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_200_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=200 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-200-double.wav
  set $text_comma_200_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_200_tab + 0) = 72
  set {char}($text_comma_200_tab + 1) = 101
  set {char}($text_comma_200_tab + 2) = 108
  set {char}($text_comma_200_tab + 3) = 108
  set {char}($text_comma_200_tab + 4) = 111
  set {char}($text_comma_200_tab + 5) = 44
  set {char}($text_comma_200_tab + 6) = 9
  set {char}($text_comma_200_tab + 7) = 119
  set {char}($text_comma_200_tab + 8) = 111
  set {char}($text_comma_200_tab + 9) = 114
  set {char}($text_comma_200_tab + 10) = 108
  set {char}($text_comma_200_tab + 11) = 100
  set {char}($text_comma_200_tab + 12) = 46
  set {char}($text_comma_200_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_200_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=200 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-200-tab.wav
  set $text_comma_200_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_200_no_comma + 0) = 72
  set {char}($text_comma_200_no_comma + 1) = 101
  set {char}($text_comma_200_no_comma + 2) = 108
  set {char}($text_comma_200_no_comma + 3) = 108
  set {char}($text_comma_200_no_comma + 4) = 111
  set {char}($text_comma_200_no_comma + 5) = 32
  set {char}($text_comma_200_no_comma + 6) = 119
  set {char}($text_comma_200_no_comma + 7) = 111
  set {char}($text_comma_200_no_comma + 8) = 114
  set {char}($text_comma_200_no_comma + 9) = 108
  set {char}($text_comma_200_no_comma + 10) = 100
  set {char}($text_comma_200_no_comma + 11) = 46
  set {char}($text_comma_200_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_200_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=200 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-200-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(201, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=201 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_201_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_201_space + 0) = 72
  set {char}($text_comma_201_space + 1) = 101
  set {char}($text_comma_201_space + 2) = 108
  set {char}($text_comma_201_space + 3) = 108
  set {char}($text_comma_201_space + 4) = 111
  set {char}($text_comma_201_space + 5) = 44
  set {char}($text_comma_201_space + 6) = 32
  set {char}($text_comma_201_space + 7) = 119
  set {char}($text_comma_201_space + 8) = 111
  set {char}($text_comma_201_space + 9) = 114
  set {char}($text_comma_201_space + 10) = 108
  set {char}($text_comma_201_space + 11) = 100
  set {char}($text_comma_201_space + 12) = 46
  set {char}($text_comma_201_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_201_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=201 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-201-space.wav
  set $text_comma_201_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_201_adjacent + 0) = 72
  set {char}($text_comma_201_adjacent + 1) = 101
  set {char}($text_comma_201_adjacent + 2) = 108
  set {char}($text_comma_201_adjacent + 3) = 108
  set {char}($text_comma_201_adjacent + 4) = 111
  set {char}($text_comma_201_adjacent + 5) = 44
  set {char}($text_comma_201_adjacent + 6) = 119
  set {char}($text_comma_201_adjacent + 7) = 111
  set {char}($text_comma_201_adjacent + 8) = 114
  set {char}($text_comma_201_adjacent + 9) = 108
  set {char}($text_comma_201_adjacent + 10) = 100
  set {char}($text_comma_201_adjacent + 11) = 46
  set {char}($text_comma_201_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_201_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=201 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-201-adjacent.wav
  set $text_comma_201_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_201_double + 0) = 72
  set {char}($text_comma_201_double + 1) = 101
  set {char}($text_comma_201_double + 2) = 108
  set {char}($text_comma_201_double + 3) = 108
  set {char}($text_comma_201_double + 4) = 111
  set {char}($text_comma_201_double + 5) = 44
  set {char}($text_comma_201_double + 6) = 32
  set {char}($text_comma_201_double + 7) = 32
  set {char}($text_comma_201_double + 8) = 119
  set {char}($text_comma_201_double + 9) = 111
  set {char}($text_comma_201_double + 10) = 114
  set {char}($text_comma_201_double + 11) = 108
  set {char}($text_comma_201_double + 12) = 100
  set {char}($text_comma_201_double + 13) = 46
  set {char}($text_comma_201_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_201_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=201 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-201-double.wav
  set $text_comma_201_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_201_tab + 0) = 72
  set {char}($text_comma_201_tab + 1) = 101
  set {char}($text_comma_201_tab + 2) = 108
  set {char}($text_comma_201_tab + 3) = 108
  set {char}($text_comma_201_tab + 4) = 111
  set {char}($text_comma_201_tab + 5) = 44
  set {char}($text_comma_201_tab + 6) = 9
  set {char}($text_comma_201_tab + 7) = 119
  set {char}($text_comma_201_tab + 8) = 111
  set {char}($text_comma_201_tab + 9) = 114
  set {char}($text_comma_201_tab + 10) = 108
  set {char}($text_comma_201_tab + 11) = 100
  set {char}($text_comma_201_tab + 12) = 46
  set {char}($text_comma_201_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_201_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=201 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-201-tab.wav
  set $text_comma_201_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_201_no_comma + 0) = 72
  set {char}($text_comma_201_no_comma + 1) = 101
  set {char}($text_comma_201_no_comma + 2) = 108
  set {char}($text_comma_201_no_comma + 3) = 108
  set {char}($text_comma_201_no_comma + 4) = 111
  set {char}($text_comma_201_no_comma + 5) = 32
  set {char}($text_comma_201_no_comma + 6) = 119
  set {char}($text_comma_201_no_comma + 7) = 111
  set {char}($text_comma_201_no_comma + 8) = 114
  set {char}($text_comma_201_no_comma + 9) = 108
  set {char}($text_comma_201_no_comma + 10) = 100
  set {char}($text_comma_201_no_comma + 11) = 46
  set {char}($text_comma_201_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_201_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=201 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-201-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(249, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=249 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_249_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_249_space + 0) = 72
  set {char}($text_comma_249_space + 1) = 101
  set {char}($text_comma_249_space + 2) = 108
  set {char}($text_comma_249_space + 3) = 108
  set {char}($text_comma_249_space + 4) = 111
  set {char}($text_comma_249_space + 5) = 44
  set {char}($text_comma_249_space + 6) = 32
  set {char}($text_comma_249_space + 7) = 119
  set {char}($text_comma_249_space + 8) = 111
  set {char}($text_comma_249_space + 9) = 114
  set {char}($text_comma_249_space + 10) = 108
  set {char}($text_comma_249_space + 11) = 100
  set {char}($text_comma_249_space + 12) = 46
  set {char}($text_comma_249_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_249_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=249 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-249-space.wav
  set $text_comma_249_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_249_adjacent + 0) = 72
  set {char}($text_comma_249_adjacent + 1) = 101
  set {char}($text_comma_249_adjacent + 2) = 108
  set {char}($text_comma_249_adjacent + 3) = 108
  set {char}($text_comma_249_adjacent + 4) = 111
  set {char}($text_comma_249_adjacent + 5) = 44
  set {char}($text_comma_249_adjacent + 6) = 119
  set {char}($text_comma_249_adjacent + 7) = 111
  set {char}($text_comma_249_adjacent + 8) = 114
  set {char}($text_comma_249_adjacent + 9) = 108
  set {char}($text_comma_249_adjacent + 10) = 100
  set {char}($text_comma_249_adjacent + 11) = 46
  set {char}($text_comma_249_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_249_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=249 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-249-adjacent.wav
  set $text_comma_249_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_249_double + 0) = 72
  set {char}($text_comma_249_double + 1) = 101
  set {char}($text_comma_249_double + 2) = 108
  set {char}($text_comma_249_double + 3) = 108
  set {char}($text_comma_249_double + 4) = 111
  set {char}($text_comma_249_double + 5) = 44
  set {char}($text_comma_249_double + 6) = 32
  set {char}($text_comma_249_double + 7) = 32
  set {char}($text_comma_249_double + 8) = 119
  set {char}($text_comma_249_double + 9) = 111
  set {char}($text_comma_249_double + 10) = 114
  set {char}($text_comma_249_double + 11) = 108
  set {char}($text_comma_249_double + 12) = 100
  set {char}($text_comma_249_double + 13) = 46
  set {char}($text_comma_249_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_249_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=249 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-249-double.wav
  set $text_comma_249_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_249_tab + 0) = 72
  set {char}($text_comma_249_tab + 1) = 101
  set {char}($text_comma_249_tab + 2) = 108
  set {char}($text_comma_249_tab + 3) = 108
  set {char}($text_comma_249_tab + 4) = 111
  set {char}($text_comma_249_tab + 5) = 44
  set {char}($text_comma_249_tab + 6) = 9
  set {char}($text_comma_249_tab + 7) = 119
  set {char}($text_comma_249_tab + 8) = 111
  set {char}($text_comma_249_tab + 9) = 114
  set {char}($text_comma_249_tab + 10) = 108
  set {char}($text_comma_249_tab + 11) = 100
  set {char}($text_comma_249_tab + 12) = 46
  set {char}($text_comma_249_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_249_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=249 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-249-tab.wav
  set $text_comma_249_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_249_no_comma + 0) = 72
  set {char}($text_comma_249_no_comma + 1) = 101
  set {char}($text_comma_249_no_comma + 2) = 108
  set {char}($text_comma_249_no_comma + 3) = 108
  set {char}($text_comma_249_no_comma + 4) = 111
  set {char}($text_comma_249_no_comma + 5) = 32
  set {char}($text_comma_249_no_comma + 6) = 119
  set {char}($text_comma_249_no_comma + 7) = 111
  set {char}($text_comma_249_no_comma + 8) = 114
  set {char}($text_comma_249_no_comma + 9) = 108
  set {char}($text_comma_249_no_comma + 10) = 100
  set {char}($text_comma_249_no_comma + 11) = 46
  set {char}($text_comma_249_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_249_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=249 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-249-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(250, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=250 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_250_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_250_space + 0) = 72
  set {char}($text_comma_250_space + 1) = 101
  set {char}($text_comma_250_space + 2) = 108
  set {char}($text_comma_250_space + 3) = 108
  set {char}($text_comma_250_space + 4) = 111
  set {char}($text_comma_250_space + 5) = 44
  set {char}($text_comma_250_space + 6) = 32
  set {char}($text_comma_250_space + 7) = 119
  set {char}($text_comma_250_space + 8) = 111
  set {char}($text_comma_250_space + 9) = 114
  set {char}($text_comma_250_space + 10) = 108
  set {char}($text_comma_250_space + 11) = 100
  set {char}($text_comma_250_space + 12) = 46
  set {char}($text_comma_250_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_250_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=250 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-250-space.wav
  set $text_comma_250_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_250_adjacent + 0) = 72
  set {char}($text_comma_250_adjacent + 1) = 101
  set {char}($text_comma_250_adjacent + 2) = 108
  set {char}($text_comma_250_adjacent + 3) = 108
  set {char}($text_comma_250_adjacent + 4) = 111
  set {char}($text_comma_250_adjacent + 5) = 44
  set {char}($text_comma_250_adjacent + 6) = 119
  set {char}($text_comma_250_adjacent + 7) = 111
  set {char}($text_comma_250_adjacent + 8) = 114
  set {char}($text_comma_250_adjacent + 9) = 108
  set {char}($text_comma_250_adjacent + 10) = 100
  set {char}($text_comma_250_adjacent + 11) = 46
  set {char}($text_comma_250_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_250_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=250 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-250-adjacent.wav
  set $text_comma_250_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_250_double + 0) = 72
  set {char}($text_comma_250_double + 1) = 101
  set {char}($text_comma_250_double + 2) = 108
  set {char}($text_comma_250_double + 3) = 108
  set {char}($text_comma_250_double + 4) = 111
  set {char}($text_comma_250_double + 5) = 44
  set {char}($text_comma_250_double + 6) = 32
  set {char}($text_comma_250_double + 7) = 32
  set {char}($text_comma_250_double + 8) = 119
  set {char}($text_comma_250_double + 9) = 111
  set {char}($text_comma_250_double + 10) = 114
  set {char}($text_comma_250_double + 11) = 108
  set {char}($text_comma_250_double + 12) = 100
  set {char}($text_comma_250_double + 13) = 46
  set {char}($text_comma_250_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_250_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=250 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-250-double.wav
  set $text_comma_250_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_250_tab + 0) = 72
  set {char}($text_comma_250_tab + 1) = 101
  set {char}($text_comma_250_tab + 2) = 108
  set {char}($text_comma_250_tab + 3) = 108
  set {char}($text_comma_250_tab + 4) = 111
  set {char}($text_comma_250_tab + 5) = 44
  set {char}($text_comma_250_tab + 6) = 9
  set {char}($text_comma_250_tab + 7) = 119
  set {char}($text_comma_250_tab + 8) = 111
  set {char}($text_comma_250_tab + 9) = 114
  set {char}($text_comma_250_tab + 10) = 108
  set {char}($text_comma_250_tab + 11) = 100
  set {char}($text_comma_250_tab + 12) = 46
  set {char}($text_comma_250_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_250_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=250 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-250-tab.wav
  set $text_comma_250_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_250_no_comma + 0) = 72
  set {char}($text_comma_250_no_comma + 1) = 101
  set {char}($text_comma_250_no_comma + 2) = 108
  set {char}($text_comma_250_no_comma + 3) = 108
  set {char}($text_comma_250_no_comma + 4) = 111
  set {char}($text_comma_250_no_comma + 5) = 32
  set {char}($text_comma_250_no_comma + 6) = 119
  set {char}($text_comma_250_no_comma + 7) = 111
  set {char}($text_comma_250_no_comma + 8) = 114
  set {char}($text_comma_250_no_comma + 9) = 108
  set {char}($text_comma_250_no_comma + 10) = 100
  set {char}($text_comma_250_no_comma + 11) = 46
  set {char}($text_comma_250_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_250_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=250 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-250-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(251, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=251 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_251_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_251_space + 0) = 72
  set {char}($text_comma_251_space + 1) = 101
  set {char}($text_comma_251_space + 2) = 108
  set {char}($text_comma_251_space + 3) = 108
  set {char}($text_comma_251_space + 4) = 111
  set {char}($text_comma_251_space + 5) = 44
  set {char}($text_comma_251_space + 6) = 32
  set {char}($text_comma_251_space + 7) = 119
  set {char}($text_comma_251_space + 8) = 111
  set {char}($text_comma_251_space + 9) = 114
  set {char}($text_comma_251_space + 10) = 108
  set {char}($text_comma_251_space + 11) = 100
  set {char}($text_comma_251_space + 12) = 46
  set {char}($text_comma_251_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_251_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=251 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-251-space.wav
  set $text_comma_251_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_251_adjacent + 0) = 72
  set {char}($text_comma_251_adjacent + 1) = 101
  set {char}($text_comma_251_adjacent + 2) = 108
  set {char}($text_comma_251_adjacent + 3) = 108
  set {char}($text_comma_251_adjacent + 4) = 111
  set {char}($text_comma_251_adjacent + 5) = 44
  set {char}($text_comma_251_adjacent + 6) = 119
  set {char}($text_comma_251_adjacent + 7) = 111
  set {char}($text_comma_251_adjacent + 8) = 114
  set {char}($text_comma_251_adjacent + 9) = 108
  set {char}($text_comma_251_adjacent + 10) = 100
  set {char}($text_comma_251_adjacent + 11) = 46
  set {char}($text_comma_251_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_251_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=251 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-251-adjacent.wav
  set $text_comma_251_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_251_double + 0) = 72
  set {char}($text_comma_251_double + 1) = 101
  set {char}($text_comma_251_double + 2) = 108
  set {char}($text_comma_251_double + 3) = 108
  set {char}($text_comma_251_double + 4) = 111
  set {char}($text_comma_251_double + 5) = 44
  set {char}($text_comma_251_double + 6) = 32
  set {char}($text_comma_251_double + 7) = 32
  set {char}($text_comma_251_double + 8) = 119
  set {char}($text_comma_251_double + 9) = 111
  set {char}($text_comma_251_double + 10) = 114
  set {char}($text_comma_251_double + 11) = 108
  set {char}($text_comma_251_double + 12) = 100
  set {char}($text_comma_251_double + 13) = 46
  set {char}($text_comma_251_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_251_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=251 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-251-double.wav
  set $text_comma_251_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_251_tab + 0) = 72
  set {char}($text_comma_251_tab + 1) = 101
  set {char}($text_comma_251_tab + 2) = 108
  set {char}($text_comma_251_tab + 3) = 108
  set {char}($text_comma_251_tab + 4) = 111
  set {char}($text_comma_251_tab + 5) = 44
  set {char}($text_comma_251_tab + 6) = 9
  set {char}($text_comma_251_tab + 7) = 119
  set {char}($text_comma_251_tab + 8) = 111
  set {char}($text_comma_251_tab + 9) = 114
  set {char}($text_comma_251_tab + 10) = 108
  set {char}($text_comma_251_tab + 11) = 100
  set {char}($text_comma_251_tab + 12) = 46
  set {char}($text_comma_251_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_251_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=251 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-251-tab.wav
  set $text_comma_251_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_251_no_comma + 0) = 72
  set {char}($text_comma_251_no_comma + 1) = 101
  set {char}($text_comma_251_no_comma + 2) = 108
  set {char}($text_comma_251_no_comma + 3) = 108
  set {char}($text_comma_251_no_comma + 4) = 111
  set {char}($text_comma_251_no_comma + 5) = 32
  set {char}($text_comma_251_no_comma + 6) = 119
  set {char}($text_comma_251_no_comma + 7) = 111
  set {char}($text_comma_251_no_comma + 8) = 114
  set {char}($text_comma_251_no_comma + 9) = 108
  set {char}($text_comma_251_no_comma + 10) = 100
  set {char}($text_comma_251_no_comma + 11) = 46
  set {char}($text_comma_251_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_251_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=251 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-251-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(499, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=499 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_499_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_499_space + 0) = 72
  set {char}($text_comma_499_space + 1) = 101
  set {char}($text_comma_499_space + 2) = 108
  set {char}($text_comma_499_space + 3) = 108
  set {char}($text_comma_499_space + 4) = 111
  set {char}($text_comma_499_space + 5) = 44
  set {char}($text_comma_499_space + 6) = 32
  set {char}($text_comma_499_space + 7) = 119
  set {char}($text_comma_499_space + 8) = 111
  set {char}($text_comma_499_space + 9) = 114
  set {char}($text_comma_499_space + 10) = 108
  set {char}($text_comma_499_space + 11) = 100
  set {char}($text_comma_499_space + 12) = 46
  set {char}($text_comma_499_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_499_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=499 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-499-space.wav
  set $text_comma_499_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_499_adjacent + 0) = 72
  set {char}($text_comma_499_adjacent + 1) = 101
  set {char}($text_comma_499_adjacent + 2) = 108
  set {char}($text_comma_499_adjacent + 3) = 108
  set {char}($text_comma_499_adjacent + 4) = 111
  set {char}($text_comma_499_adjacent + 5) = 44
  set {char}($text_comma_499_adjacent + 6) = 119
  set {char}($text_comma_499_adjacent + 7) = 111
  set {char}($text_comma_499_adjacent + 8) = 114
  set {char}($text_comma_499_adjacent + 9) = 108
  set {char}($text_comma_499_adjacent + 10) = 100
  set {char}($text_comma_499_adjacent + 11) = 46
  set {char}($text_comma_499_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_499_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=499 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-499-adjacent.wav
  set $text_comma_499_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_499_double + 0) = 72
  set {char}($text_comma_499_double + 1) = 101
  set {char}($text_comma_499_double + 2) = 108
  set {char}($text_comma_499_double + 3) = 108
  set {char}($text_comma_499_double + 4) = 111
  set {char}($text_comma_499_double + 5) = 44
  set {char}($text_comma_499_double + 6) = 32
  set {char}($text_comma_499_double + 7) = 32
  set {char}($text_comma_499_double + 8) = 119
  set {char}($text_comma_499_double + 9) = 111
  set {char}($text_comma_499_double + 10) = 114
  set {char}($text_comma_499_double + 11) = 108
  set {char}($text_comma_499_double + 12) = 100
  set {char}($text_comma_499_double + 13) = 46
  set {char}($text_comma_499_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_499_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=499 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-499-double.wav
  set $text_comma_499_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_499_tab + 0) = 72
  set {char}($text_comma_499_tab + 1) = 101
  set {char}($text_comma_499_tab + 2) = 108
  set {char}($text_comma_499_tab + 3) = 108
  set {char}($text_comma_499_tab + 4) = 111
  set {char}($text_comma_499_tab + 5) = 44
  set {char}($text_comma_499_tab + 6) = 9
  set {char}($text_comma_499_tab + 7) = 119
  set {char}($text_comma_499_tab + 8) = 111
  set {char}($text_comma_499_tab + 9) = 114
  set {char}($text_comma_499_tab + 10) = 108
  set {char}($text_comma_499_tab + 11) = 100
  set {char}($text_comma_499_tab + 12) = 46
  set {char}($text_comma_499_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_499_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=499 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-499-tab.wav
  set $text_comma_499_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_499_no_comma + 0) = 72
  set {char}($text_comma_499_no_comma + 1) = 101
  set {char}($text_comma_499_no_comma + 2) = 108
  set {char}($text_comma_499_no_comma + 3) = 108
  set {char}($text_comma_499_no_comma + 4) = 111
  set {char}($text_comma_499_no_comma + 5) = 32
  set {char}($text_comma_499_no_comma + 6) = 119
  set {char}($text_comma_499_no_comma + 7) = 111
  set {char}($text_comma_499_no_comma + 8) = 114
  set {char}($text_comma_499_no_comma + 9) = 108
  set {char}($text_comma_499_no_comma + 10) = 100
  set {char}($text_comma_499_no_comma + 11) = 46
  set {char}($text_comma_499_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_499_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=499 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-499-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(500, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=500 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_500_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_500_space + 0) = 72
  set {char}($text_comma_500_space + 1) = 101
  set {char}($text_comma_500_space + 2) = 108
  set {char}($text_comma_500_space + 3) = 108
  set {char}($text_comma_500_space + 4) = 111
  set {char}($text_comma_500_space + 5) = 44
  set {char}($text_comma_500_space + 6) = 32
  set {char}($text_comma_500_space + 7) = 119
  set {char}($text_comma_500_space + 8) = 111
  set {char}($text_comma_500_space + 9) = 114
  set {char}($text_comma_500_space + 10) = 108
  set {char}($text_comma_500_space + 11) = 100
  set {char}($text_comma_500_space + 12) = 46
  set {char}($text_comma_500_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_500_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=500 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-500-space.wav
  set $text_comma_500_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_500_adjacent + 0) = 72
  set {char}($text_comma_500_adjacent + 1) = 101
  set {char}($text_comma_500_adjacent + 2) = 108
  set {char}($text_comma_500_adjacent + 3) = 108
  set {char}($text_comma_500_adjacent + 4) = 111
  set {char}($text_comma_500_adjacent + 5) = 44
  set {char}($text_comma_500_adjacent + 6) = 119
  set {char}($text_comma_500_adjacent + 7) = 111
  set {char}($text_comma_500_adjacent + 8) = 114
  set {char}($text_comma_500_adjacent + 9) = 108
  set {char}($text_comma_500_adjacent + 10) = 100
  set {char}($text_comma_500_adjacent + 11) = 46
  set {char}($text_comma_500_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_500_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=500 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-500-adjacent.wav
  set $text_comma_500_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_500_double + 0) = 72
  set {char}($text_comma_500_double + 1) = 101
  set {char}($text_comma_500_double + 2) = 108
  set {char}($text_comma_500_double + 3) = 108
  set {char}($text_comma_500_double + 4) = 111
  set {char}($text_comma_500_double + 5) = 44
  set {char}($text_comma_500_double + 6) = 32
  set {char}($text_comma_500_double + 7) = 32
  set {char}($text_comma_500_double + 8) = 119
  set {char}($text_comma_500_double + 9) = 111
  set {char}($text_comma_500_double + 10) = 114
  set {char}($text_comma_500_double + 11) = 108
  set {char}($text_comma_500_double + 12) = 100
  set {char}($text_comma_500_double + 13) = 46
  set {char}($text_comma_500_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_500_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=500 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-500-double.wav
  set $text_comma_500_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_500_tab + 0) = 72
  set {char}($text_comma_500_tab + 1) = 101
  set {char}($text_comma_500_tab + 2) = 108
  set {char}($text_comma_500_tab + 3) = 108
  set {char}($text_comma_500_tab + 4) = 111
  set {char}($text_comma_500_tab + 5) = 44
  set {char}($text_comma_500_tab + 6) = 9
  set {char}($text_comma_500_tab + 7) = 119
  set {char}($text_comma_500_tab + 8) = 111
  set {char}($text_comma_500_tab + 9) = 114
  set {char}($text_comma_500_tab + 10) = 108
  set {char}($text_comma_500_tab + 11) = 100
  set {char}($text_comma_500_tab + 12) = 46
  set {char}($text_comma_500_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_500_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=500 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-500-tab.wav
  set $text_comma_500_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_500_no_comma + 0) = 72
  set {char}($text_comma_500_no_comma + 1) = 101
  set {char}($text_comma_500_no_comma + 2) = 108
  set {char}($text_comma_500_no_comma + 3) = 108
  set {char}($text_comma_500_no_comma + 4) = 111
  set {char}($text_comma_500_no_comma + 5) = 32
  set {char}($text_comma_500_no_comma + 6) = 119
  set {char}($text_comma_500_no_comma + 7) = 111
  set {char}($text_comma_500_no_comma + 8) = 114
  set {char}($text_comma_500_no_comma + 9) = 108
  set {char}($text_comma_500_no_comma + 10) = 100
  set {char}($text_comma_500_no_comma + 11) = 46
  set {char}($text_comma_500_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_500_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=500 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-500-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(501, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=501 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_501_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_501_space + 0) = 72
  set {char}($text_comma_501_space + 1) = 101
  set {char}($text_comma_501_space + 2) = 108
  set {char}($text_comma_501_space + 3) = 108
  set {char}($text_comma_501_space + 4) = 111
  set {char}($text_comma_501_space + 5) = 44
  set {char}($text_comma_501_space + 6) = 32
  set {char}($text_comma_501_space + 7) = 119
  set {char}($text_comma_501_space + 8) = 111
  set {char}($text_comma_501_space + 9) = 114
  set {char}($text_comma_501_space + 10) = 108
  set {char}($text_comma_501_space + 11) = 100
  set {char}($text_comma_501_space + 12) = 46
  set {char}($text_comma_501_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_501_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=501 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-501-space.wav
  set $text_comma_501_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_501_adjacent + 0) = 72
  set {char}($text_comma_501_adjacent + 1) = 101
  set {char}($text_comma_501_adjacent + 2) = 108
  set {char}($text_comma_501_adjacent + 3) = 108
  set {char}($text_comma_501_adjacent + 4) = 111
  set {char}($text_comma_501_adjacent + 5) = 44
  set {char}($text_comma_501_adjacent + 6) = 119
  set {char}($text_comma_501_adjacent + 7) = 111
  set {char}($text_comma_501_adjacent + 8) = 114
  set {char}($text_comma_501_adjacent + 9) = 108
  set {char}($text_comma_501_adjacent + 10) = 100
  set {char}($text_comma_501_adjacent + 11) = 46
  set {char}($text_comma_501_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_501_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=501 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-501-adjacent.wav
  set $text_comma_501_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_501_double + 0) = 72
  set {char}($text_comma_501_double + 1) = 101
  set {char}($text_comma_501_double + 2) = 108
  set {char}($text_comma_501_double + 3) = 108
  set {char}($text_comma_501_double + 4) = 111
  set {char}($text_comma_501_double + 5) = 44
  set {char}($text_comma_501_double + 6) = 32
  set {char}($text_comma_501_double + 7) = 32
  set {char}($text_comma_501_double + 8) = 119
  set {char}($text_comma_501_double + 9) = 111
  set {char}($text_comma_501_double + 10) = 114
  set {char}($text_comma_501_double + 11) = 108
  set {char}($text_comma_501_double + 12) = 100
  set {char}($text_comma_501_double + 13) = 46
  set {char}($text_comma_501_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_501_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=501 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-501-double.wav
  set $text_comma_501_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_501_tab + 0) = 72
  set {char}($text_comma_501_tab + 1) = 101
  set {char}($text_comma_501_tab + 2) = 108
  set {char}($text_comma_501_tab + 3) = 108
  set {char}($text_comma_501_tab + 4) = 111
  set {char}($text_comma_501_tab + 5) = 44
  set {char}($text_comma_501_tab + 6) = 9
  set {char}($text_comma_501_tab + 7) = 119
  set {char}($text_comma_501_tab + 8) = 111
  set {char}($text_comma_501_tab + 9) = 114
  set {char}($text_comma_501_tab + 10) = 108
  set {char}($text_comma_501_tab + 11) = 100
  set {char}($text_comma_501_tab + 12) = 46
  set {char}($text_comma_501_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_501_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=501 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-501-tab.wav
  set $text_comma_501_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_501_no_comma + 0) = 72
  set {char}($text_comma_501_no_comma + 1) = 101
  set {char}($text_comma_501_no_comma + 2) = 108
  set {char}($text_comma_501_no_comma + 3) = 108
  set {char}($text_comma_501_no_comma + 4) = 111
  set {char}($text_comma_501_no_comma + 5) = 32
  set {char}($text_comma_501_no_comma + 6) = 119
  set {char}($text_comma_501_no_comma + 7) = 111
  set {char}($text_comma_501_no_comma + 8) = 114
  set {char}($text_comma_501_no_comma + 9) = 108
  set {char}($text_comma_501_no_comma + 10) = 100
  set {char}($text_comma_501_no_comma + 11) = 46
  set {char}($text_comma_501_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_501_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=501 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-501-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(924, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=924 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_924_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_924_space + 0) = 72
  set {char}($text_comma_924_space + 1) = 101
  set {char}($text_comma_924_space + 2) = 108
  set {char}($text_comma_924_space + 3) = 108
  set {char}($text_comma_924_space + 4) = 111
  set {char}($text_comma_924_space + 5) = 44
  set {char}($text_comma_924_space + 6) = 32
  set {char}($text_comma_924_space + 7) = 119
  set {char}($text_comma_924_space + 8) = 111
  set {char}($text_comma_924_space + 9) = 114
  set {char}($text_comma_924_space + 10) = 108
  set {char}($text_comma_924_space + 11) = 100
  set {char}($text_comma_924_space + 12) = 46
  set {char}($text_comma_924_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_924_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=924 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-924-space.wav
  set $text_comma_924_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_924_adjacent + 0) = 72
  set {char}($text_comma_924_adjacent + 1) = 101
  set {char}($text_comma_924_adjacent + 2) = 108
  set {char}($text_comma_924_adjacent + 3) = 108
  set {char}($text_comma_924_adjacent + 4) = 111
  set {char}($text_comma_924_adjacent + 5) = 44
  set {char}($text_comma_924_adjacent + 6) = 119
  set {char}($text_comma_924_adjacent + 7) = 111
  set {char}($text_comma_924_adjacent + 8) = 114
  set {char}($text_comma_924_adjacent + 9) = 108
  set {char}($text_comma_924_adjacent + 10) = 100
  set {char}($text_comma_924_adjacent + 11) = 46
  set {char}($text_comma_924_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_924_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=924 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-924-adjacent.wav
  set $text_comma_924_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_924_double + 0) = 72
  set {char}($text_comma_924_double + 1) = 101
  set {char}($text_comma_924_double + 2) = 108
  set {char}($text_comma_924_double + 3) = 108
  set {char}($text_comma_924_double + 4) = 111
  set {char}($text_comma_924_double + 5) = 44
  set {char}($text_comma_924_double + 6) = 32
  set {char}($text_comma_924_double + 7) = 32
  set {char}($text_comma_924_double + 8) = 119
  set {char}($text_comma_924_double + 9) = 111
  set {char}($text_comma_924_double + 10) = 114
  set {char}($text_comma_924_double + 11) = 108
  set {char}($text_comma_924_double + 12) = 100
  set {char}($text_comma_924_double + 13) = 46
  set {char}($text_comma_924_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_924_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=924 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-924-double.wav
  set $text_comma_924_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_924_tab + 0) = 72
  set {char}($text_comma_924_tab + 1) = 101
  set {char}($text_comma_924_tab + 2) = 108
  set {char}($text_comma_924_tab + 3) = 108
  set {char}($text_comma_924_tab + 4) = 111
  set {char}($text_comma_924_tab + 5) = 44
  set {char}($text_comma_924_tab + 6) = 9
  set {char}($text_comma_924_tab + 7) = 119
  set {char}($text_comma_924_tab + 8) = 111
  set {char}($text_comma_924_tab + 9) = 114
  set {char}($text_comma_924_tab + 10) = 108
  set {char}($text_comma_924_tab + 11) = 100
  set {char}($text_comma_924_tab + 12) = 46
  set {char}($text_comma_924_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_924_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=924 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-924-tab.wav
  set $text_comma_924_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_924_no_comma + 0) = 72
  set {char}($text_comma_924_no_comma + 1) = 101
  set {char}($text_comma_924_no_comma + 2) = 108
  set {char}($text_comma_924_no_comma + 3) = 108
  set {char}($text_comma_924_no_comma + 4) = 111
  set {char}($text_comma_924_no_comma + 5) = 32
  set {char}($text_comma_924_no_comma + 6) = 119
  set {char}($text_comma_924_no_comma + 7) = 111
  set {char}($text_comma_924_no_comma + 8) = 114
  set {char}($text_comma_924_no_comma + 9) = 108
  set {char}($text_comma_924_no_comma + 10) = 100
  set {char}($text_comma_924_no_comma + 11) = 46
  set {char}($text_comma_924_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_924_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=924 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-924-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(925, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=925 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_925_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_925_space + 0) = 72
  set {char}($text_comma_925_space + 1) = 101
  set {char}($text_comma_925_space + 2) = 108
  set {char}($text_comma_925_space + 3) = 108
  set {char}($text_comma_925_space + 4) = 111
  set {char}($text_comma_925_space + 5) = 44
  set {char}($text_comma_925_space + 6) = 32
  set {char}($text_comma_925_space + 7) = 119
  set {char}($text_comma_925_space + 8) = 111
  set {char}($text_comma_925_space + 9) = 114
  set {char}($text_comma_925_space + 10) = 108
  set {char}($text_comma_925_space + 11) = 100
  set {char}($text_comma_925_space + 12) = 46
  set {char}($text_comma_925_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_925_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=925 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-925-space.wav
  set $text_comma_925_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_925_adjacent + 0) = 72
  set {char}($text_comma_925_adjacent + 1) = 101
  set {char}($text_comma_925_adjacent + 2) = 108
  set {char}($text_comma_925_adjacent + 3) = 108
  set {char}($text_comma_925_adjacent + 4) = 111
  set {char}($text_comma_925_adjacent + 5) = 44
  set {char}($text_comma_925_adjacent + 6) = 119
  set {char}($text_comma_925_adjacent + 7) = 111
  set {char}($text_comma_925_adjacent + 8) = 114
  set {char}($text_comma_925_adjacent + 9) = 108
  set {char}($text_comma_925_adjacent + 10) = 100
  set {char}($text_comma_925_adjacent + 11) = 46
  set {char}($text_comma_925_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_925_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=925 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-925-adjacent.wav
  set $text_comma_925_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_925_double + 0) = 72
  set {char}($text_comma_925_double + 1) = 101
  set {char}($text_comma_925_double + 2) = 108
  set {char}($text_comma_925_double + 3) = 108
  set {char}($text_comma_925_double + 4) = 111
  set {char}($text_comma_925_double + 5) = 44
  set {char}($text_comma_925_double + 6) = 32
  set {char}($text_comma_925_double + 7) = 32
  set {char}($text_comma_925_double + 8) = 119
  set {char}($text_comma_925_double + 9) = 111
  set {char}($text_comma_925_double + 10) = 114
  set {char}($text_comma_925_double + 11) = 108
  set {char}($text_comma_925_double + 12) = 100
  set {char}($text_comma_925_double + 13) = 46
  set {char}($text_comma_925_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_925_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=925 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-925-double.wav
  set $text_comma_925_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_925_tab + 0) = 72
  set {char}($text_comma_925_tab + 1) = 101
  set {char}($text_comma_925_tab + 2) = 108
  set {char}($text_comma_925_tab + 3) = 108
  set {char}($text_comma_925_tab + 4) = 111
  set {char}($text_comma_925_tab + 5) = 44
  set {char}($text_comma_925_tab + 6) = 9
  set {char}($text_comma_925_tab + 7) = 119
  set {char}($text_comma_925_tab + 8) = 111
  set {char}($text_comma_925_tab + 9) = 114
  set {char}($text_comma_925_tab + 10) = 108
  set {char}($text_comma_925_tab + 11) = 100
  set {char}($text_comma_925_tab + 12) = 46
  set {char}($text_comma_925_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_925_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=925 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-925-tab.wav
  set $text_comma_925_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_925_no_comma + 0) = 72
  set {char}($text_comma_925_no_comma + 1) = 101
  set {char}($text_comma_925_no_comma + 2) = 108
  set {char}($text_comma_925_no_comma + 3) = 108
  set {char}($text_comma_925_no_comma + 4) = 111
  set {char}($text_comma_925_no_comma + 5) = 32
  set {char}($text_comma_925_no_comma + 6) = 119
  set {char}($text_comma_925_no_comma + 7) = 111
  set {char}($text_comma_925_no_comma + 8) = 114
  set {char}($text_comma_925_no_comma + 9) = 108
  set {char}($text_comma_925_no_comma + 10) = 100
  set {char}($text_comma_925_no_comma + 11) = 46
  set {char}($text_comma_925_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_925_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=925 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-925-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(926, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=926 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_926_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_926_space + 0) = 72
  set {char}($text_comma_926_space + 1) = 101
  set {char}($text_comma_926_space + 2) = 108
  set {char}($text_comma_926_space + 3) = 108
  set {char}($text_comma_926_space + 4) = 111
  set {char}($text_comma_926_space + 5) = 44
  set {char}($text_comma_926_space + 6) = 32
  set {char}($text_comma_926_space + 7) = 119
  set {char}($text_comma_926_space + 8) = 111
  set {char}($text_comma_926_space + 9) = 114
  set {char}($text_comma_926_space + 10) = 108
  set {char}($text_comma_926_space + 11) = 100
  set {char}($text_comma_926_space + 12) = 46
  set {char}($text_comma_926_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_926_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=926 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-926-space.wav
  set $text_comma_926_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_926_adjacent + 0) = 72
  set {char}($text_comma_926_adjacent + 1) = 101
  set {char}($text_comma_926_adjacent + 2) = 108
  set {char}($text_comma_926_adjacent + 3) = 108
  set {char}($text_comma_926_adjacent + 4) = 111
  set {char}($text_comma_926_adjacent + 5) = 44
  set {char}($text_comma_926_adjacent + 6) = 119
  set {char}($text_comma_926_adjacent + 7) = 111
  set {char}($text_comma_926_adjacent + 8) = 114
  set {char}($text_comma_926_adjacent + 9) = 108
  set {char}($text_comma_926_adjacent + 10) = 100
  set {char}($text_comma_926_adjacent + 11) = 46
  set {char}($text_comma_926_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_926_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=926 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-926-adjacent.wav
  set $text_comma_926_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_926_double + 0) = 72
  set {char}($text_comma_926_double + 1) = 101
  set {char}($text_comma_926_double + 2) = 108
  set {char}($text_comma_926_double + 3) = 108
  set {char}($text_comma_926_double + 4) = 111
  set {char}($text_comma_926_double + 5) = 44
  set {char}($text_comma_926_double + 6) = 32
  set {char}($text_comma_926_double + 7) = 32
  set {char}($text_comma_926_double + 8) = 119
  set {char}($text_comma_926_double + 9) = 111
  set {char}($text_comma_926_double + 10) = 114
  set {char}($text_comma_926_double + 11) = 108
  set {char}($text_comma_926_double + 12) = 100
  set {char}($text_comma_926_double + 13) = 46
  set {char}($text_comma_926_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_926_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=926 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-926-double.wav
  set $text_comma_926_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_926_tab + 0) = 72
  set {char}($text_comma_926_tab + 1) = 101
  set {char}($text_comma_926_tab + 2) = 108
  set {char}($text_comma_926_tab + 3) = 108
  set {char}($text_comma_926_tab + 4) = 111
  set {char}($text_comma_926_tab + 5) = 44
  set {char}($text_comma_926_tab + 6) = 9
  set {char}($text_comma_926_tab + 7) = 119
  set {char}($text_comma_926_tab + 8) = 111
  set {char}($text_comma_926_tab + 9) = 114
  set {char}($text_comma_926_tab + 10) = 108
  set {char}($text_comma_926_tab + 11) = 100
  set {char}($text_comma_926_tab + 12) = 46
  set {char}($text_comma_926_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_926_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=926 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-926-tab.wav
  set $text_comma_926_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_926_no_comma + 0) = 72
  set {char}($text_comma_926_no_comma + 1) = 101
  set {char}($text_comma_926_no_comma + 2) = 108
  set {char}($text_comma_926_no_comma + 3) = 108
  set {char}($text_comma_926_no_comma + 4) = 111
  set {char}($text_comma_926_no_comma + 5) = 32
  set {char}($text_comma_926_no_comma + 6) = 119
  set {char}($text_comma_926_no_comma + 7) = 111
  set {char}($text_comma_926_no_comma + 8) = 114
  set {char}($text_comma_926_no_comma + 9) = 108
  set {char}($text_comma_926_no_comma + 10) = 100
  set {char}($text_comma_926_no_comma + 11) = 46
  set {char}($text_comma_926_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_926_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=926 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-926-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(65534, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=65534 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_65534_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65534_space + 0) = 72
  set {char}($text_comma_65534_space + 1) = 101
  set {char}($text_comma_65534_space + 2) = 108
  set {char}($text_comma_65534_space + 3) = 108
  set {char}($text_comma_65534_space + 4) = 111
  set {char}($text_comma_65534_space + 5) = 44
  set {char}($text_comma_65534_space + 6) = 32
  set {char}($text_comma_65534_space + 7) = 119
  set {char}($text_comma_65534_space + 8) = 111
  set {char}($text_comma_65534_space + 9) = 114
  set {char}($text_comma_65534_space + 10) = 108
  set {char}($text_comma_65534_space + 11) = 100
  set {char}($text_comma_65534_space + 12) = 46
  set {char}($text_comma_65534_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65534_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65534 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65534-space.wav
  set $text_comma_65534_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65534_adjacent + 0) = 72
  set {char}($text_comma_65534_adjacent + 1) = 101
  set {char}($text_comma_65534_adjacent + 2) = 108
  set {char}($text_comma_65534_adjacent + 3) = 108
  set {char}($text_comma_65534_adjacent + 4) = 111
  set {char}($text_comma_65534_adjacent + 5) = 44
  set {char}($text_comma_65534_adjacent + 6) = 119
  set {char}($text_comma_65534_adjacent + 7) = 111
  set {char}($text_comma_65534_adjacent + 8) = 114
  set {char}($text_comma_65534_adjacent + 9) = 108
  set {char}($text_comma_65534_adjacent + 10) = 100
  set {char}($text_comma_65534_adjacent + 11) = 46
  set {char}($text_comma_65534_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65534_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65534 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65534-adjacent.wav
  set $text_comma_65534_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65534_double + 0) = 72
  set {char}($text_comma_65534_double + 1) = 101
  set {char}($text_comma_65534_double + 2) = 108
  set {char}($text_comma_65534_double + 3) = 108
  set {char}($text_comma_65534_double + 4) = 111
  set {char}($text_comma_65534_double + 5) = 44
  set {char}($text_comma_65534_double + 6) = 32
  set {char}($text_comma_65534_double + 7) = 32
  set {char}($text_comma_65534_double + 8) = 119
  set {char}($text_comma_65534_double + 9) = 111
  set {char}($text_comma_65534_double + 10) = 114
  set {char}($text_comma_65534_double + 11) = 108
  set {char}($text_comma_65534_double + 12) = 100
  set {char}($text_comma_65534_double + 13) = 46
  set {char}($text_comma_65534_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65534_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65534 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65534-double.wav
  set $text_comma_65534_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65534_tab + 0) = 72
  set {char}($text_comma_65534_tab + 1) = 101
  set {char}($text_comma_65534_tab + 2) = 108
  set {char}($text_comma_65534_tab + 3) = 108
  set {char}($text_comma_65534_tab + 4) = 111
  set {char}($text_comma_65534_tab + 5) = 44
  set {char}($text_comma_65534_tab + 6) = 9
  set {char}($text_comma_65534_tab + 7) = 119
  set {char}($text_comma_65534_tab + 8) = 111
  set {char}($text_comma_65534_tab + 9) = 114
  set {char}($text_comma_65534_tab + 10) = 108
  set {char}($text_comma_65534_tab + 11) = 100
  set {char}($text_comma_65534_tab + 12) = 46
  set {char}($text_comma_65534_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65534_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65534 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65534-tab.wav
  set $text_comma_65534_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65534_no_comma + 0) = 72
  set {char}($text_comma_65534_no_comma + 1) = 101
  set {char}($text_comma_65534_no_comma + 2) = 108
  set {char}($text_comma_65534_no_comma + 3) = 108
  set {char}($text_comma_65534_no_comma + 4) = 111
  set {char}($text_comma_65534_no_comma + 5) = 32
  set {char}($text_comma_65534_no_comma + 6) = 119
  set {char}($text_comma_65534_no_comma + 7) = 111
  set {char}($text_comma_65534_no_comma + 8) = 114
  set {char}($text_comma_65534_no_comma + 9) = 108
  set {char}($text_comma_65534_no_comma + 10) = 100
  set {char}($text_comma_65534_no_comma + 11) = 46
  set {char}($text_comma_65534_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65534_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65534 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65534-no_comma.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(65535, $speaker)
  set $get_result = ((int (*)(int *, int))0x100281f0)($out, $speaker)
  printf "COMMA_PAUSE value=65535 getter_ret=%d getter_value=%d\n", $get_result, *$out
  set $text_comma_65535_space = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65535_space + 0) = 72
  set {char}($text_comma_65535_space + 1) = 101
  set {char}($text_comma_65535_space + 2) = 108
  set {char}($text_comma_65535_space + 3) = 108
  set {char}($text_comma_65535_space + 4) = 111
  set {char}($text_comma_65535_space + 5) = 44
  set {char}($text_comma_65535_space + 6) = 32
  set {char}($text_comma_65535_space + 7) = 119
  set {char}($text_comma_65535_space + 8) = 111
  set {char}($text_comma_65535_space + 9) = 114
  set {char}($text_comma_65535_space + 10) = 108
  set {char}($text_comma_65535_space + 11) = 100
  set {char}($text_comma_65535_space + 12) = 46
  set {char}($text_comma_65535_space + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65535_space, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65535 context=space synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65535-space.wav
  set $text_comma_65535_adjacent = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65535_adjacent + 0) = 72
  set {char}($text_comma_65535_adjacent + 1) = 101
  set {char}($text_comma_65535_adjacent + 2) = 108
  set {char}($text_comma_65535_adjacent + 3) = 108
  set {char}($text_comma_65535_adjacent + 4) = 111
  set {char}($text_comma_65535_adjacent + 5) = 44
  set {char}($text_comma_65535_adjacent + 6) = 119
  set {char}($text_comma_65535_adjacent + 7) = 111
  set {char}($text_comma_65535_adjacent + 8) = 114
  set {char}($text_comma_65535_adjacent + 9) = 108
  set {char}($text_comma_65535_adjacent + 10) = 100
  set {char}($text_comma_65535_adjacent + 11) = 46
  set {char}($text_comma_65535_adjacent + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65535_adjacent, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65535 context=adjacent synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65535-adjacent.wav
  set $text_comma_65535_double = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65535_double + 0) = 72
  set {char}($text_comma_65535_double + 1) = 101
  set {char}($text_comma_65535_double + 2) = 108
  set {char}($text_comma_65535_double + 3) = 108
  set {char}($text_comma_65535_double + 4) = 111
  set {char}($text_comma_65535_double + 5) = 44
  set {char}($text_comma_65535_double + 6) = 32
  set {char}($text_comma_65535_double + 7) = 32
  set {char}($text_comma_65535_double + 8) = 119
  set {char}($text_comma_65535_double + 9) = 111
  set {char}($text_comma_65535_double + 10) = 114
  set {char}($text_comma_65535_double + 11) = 108
  set {char}($text_comma_65535_double + 12) = 100
  set {char}($text_comma_65535_double + 13) = 46
  set {char}($text_comma_65535_double + 14) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65535_double, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65535 context=double synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65535-double.wav
  set $text_comma_65535_tab = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65535_tab + 0) = 72
  set {char}($text_comma_65535_tab + 1) = 101
  set {char}($text_comma_65535_tab + 2) = 108
  set {char}($text_comma_65535_tab + 3) = 108
  set {char}($text_comma_65535_tab + 4) = 111
  set {char}($text_comma_65535_tab + 5) = 44
  set {char}($text_comma_65535_tab + 6) = 9
  set {char}($text_comma_65535_tab + 7) = 119
  set {char}($text_comma_65535_tab + 8) = 111
  set {char}($text_comma_65535_tab + 9) = 114
  set {char}($text_comma_65535_tab + 10) = 108
  set {char}($text_comma_65535_tab + 11) = 100
  set {char}($text_comma_65535_tab + 12) = 46
  set {char}($text_comma_65535_tab + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65535_tab, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65535 context=tab synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65535-tab.wav
  set $text_comma_65535_no_comma = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_comma_65535_no_comma + 0) = 72
  set {char}($text_comma_65535_no_comma + 1) = 101
  set {char}($text_comma_65535_no_comma + 2) = 108
  set {char}($text_comma_65535_no_comma + 3) = 108
  set {char}($text_comma_65535_no_comma + 4) = 111
  set {char}($text_comma_65535_no_comma + 5) = 32
  set {char}($text_comma_65535_no_comma + 6) = 119
  set {char}($text_comma_65535_no_comma + 7) = 111
  set {char}($text_comma_65535_no_comma + 8) = 114
  set {char}($text_comma_65535_no_comma + 9) = 108
  set {char}($text_comma_65535_no_comma + 10) = 100
  set {char}($text_comma_65535_no_comma + 11) = 46
  set {char}($text_comma_65535_no_comma + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_comma_65535_no_comma, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "COMMA_PAUSE value=65535 context=no_comma synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/comma-pause-65535-no_comma.wav
  kill
  quit
end

continue
