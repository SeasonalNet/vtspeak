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
  set $out = (int *)malloc(16)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 0, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=0 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_0_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_0_upper + 0) = 72
  set {char}($text_sentence_0_upper + 1) = 101
  set {char}($text_sentence_0_upper + 2) = 108
  set {char}($text_sentence_0_upper + 3) = 108
  set {char}($text_sentence_0_upper + 4) = 111
  set {char}($text_sentence_0_upper + 5) = 46
  set {char}($text_sentence_0_upper + 6) = 32
  set {char}($text_sentence_0_upper + 7) = 87
  set {char}($text_sentence_0_upper + 8) = 111
  set {char}($text_sentence_0_upper + 9) = 114
  set {char}($text_sentence_0_upper + 10) = 108
  set {char}($text_sentence_0_upper + 11) = 100
  set {char}($text_sentence_0_upper + 12) = 46
  set {char}($text_sentence_0_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_0_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=0 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-0-upper.wav
  set $text_sentence_0_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_0_upper_lower + 0) = 72
  set {char}($text_sentence_0_upper_lower + 1) = 101
  set {char}($text_sentence_0_upper_lower + 2) = 108
  set {char}($text_sentence_0_upper_lower + 3) = 108
  set {char}($text_sentence_0_upper_lower + 4) = 111
  set {char}($text_sentence_0_upper_lower + 5) = 46
  set {char}($text_sentence_0_upper_lower + 6) = 32
  set {char}($text_sentence_0_upper_lower + 7) = 119
  set {char}($text_sentence_0_upper_lower + 8) = 111
  set {char}($text_sentence_0_upper_lower + 9) = 114
  set {char}($text_sentence_0_upper_lower + 10) = 108
  set {char}($text_sentence_0_upper_lower + 11) = 100
  set {char}($text_sentence_0_upper_lower + 12) = 46
  set {char}($text_sentence_0_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_0_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=0 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-0-upper_lower.wav
  set $text_sentence_0_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_0_lower + 0) = 104
  set {char}($text_sentence_0_lower + 1) = 101
  set {char}($text_sentence_0_lower + 2) = 108
  set {char}($text_sentence_0_lower + 3) = 108
  set {char}($text_sentence_0_lower + 4) = 111
  set {char}($text_sentence_0_lower + 5) = 46
  set {char}($text_sentence_0_lower + 6) = 32
  set {char}($text_sentence_0_lower + 7) = 119
  set {char}($text_sentence_0_lower + 8) = 111
  set {char}($text_sentence_0_lower + 9) = 114
  set {char}($text_sentence_0_lower + 10) = 108
  set {char}($text_sentence_0_lower + 11) = 100
  set {char}($text_sentence_0_lower + 12) = 46
  set {char}($text_sentence_0_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_0_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=0 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-0-lower.wav
  set $text_sentence_0_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_0_no_period + 0) = 72
  set {char}($text_sentence_0_no_period + 1) = 101
  set {char}($text_sentence_0_no_period + 2) = 108
  set {char}($text_sentence_0_no_period + 3) = 108
  set {char}($text_sentence_0_no_period + 4) = 111
  set {char}($text_sentence_0_no_period + 5) = 32
  set {char}($text_sentence_0_no_period + 6) = 119
  set {char}($text_sentence_0_no_period + 7) = 111
  set {char}($text_sentence_0_no_period + 8) = 114
  set {char}($text_sentence_0_no_period + 9) = 108
  set {char}($text_sentence_0_no_period + 10) = 100
  set {char}($text_sentence_0_no_period + 11) = 46
  set {char}($text_sentence_0_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_0_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=0 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-0-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 1, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=1 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_1_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_1_upper + 0) = 72
  set {char}($text_sentence_1_upper + 1) = 101
  set {char}($text_sentence_1_upper + 2) = 108
  set {char}($text_sentence_1_upper + 3) = 108
  set {char}($text_sentence_1_upper + 4) = 111
  set {char}($text_sentence_1_upper + 5) = 46
  set {char}($text_sentence_1_upper + 6) = 32
  set {char}($text_sentence_1_upper + 7) = 87
  set {char}($text_sentence_1_upper + 8) = 111
  set {char}($text_sentence_1_upper + 9) = 114
  set {char}($text_sentence_1_upper + 10) = 108
  set {char}($text_sentence_1_upper + 11) = 100
  set {char}($text_sentence_1_upper + 12) = 46
  set {char}($text_sentence_1_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_1_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=1 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-1-upper.wav
  set $text_sentence_1_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_1_upper_lower + 0) = 72
  set {char}($text_sentence_1_upper_lower + 1) = 101
  set {char}($text_sentence_1_upper_lower + 2) = 108
  set {char}($text_sentence_1_upper_lower + 3) = 108
  set {char}($text_sentence_1_upper_lower + 4) = 111
  set {char}($text_sentence_1_upper_lower + 5) = 46
  set {char}($text_sentence_1_upper_lower + 6) = 32
  set {char}($text_sentence_1_upper_lower + 7) = 119
  set {char}($text_sentence_1_upper_lower + 8) = 111
  set {char}($text_sentence_1_upper_lower + 9) = 114
  set {char}($text_sentence_1_upper_lower + 10) = 108
  set {char}($text_sentence_1_upper_lower + 11) = 100
  set {char}($text_sentence_1_upper_lower + 12) = 46
  set {char}($text_sentence_1_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_1_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=1 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-1-upper_lower.wav
  set $text_sentence_1_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_1_lower + 0) = 104
  set {char}($text_sentence_1_lower + 1) = 101
  set {char}($text_sentence_1_lower + 2) = 108
  set {char}($text_sentence_1_lower + 3) = 108
  set {char}($text_sentence_1_lower + 4) = 111
  set {char}($text_sentence_1_lower + 5) = 46
  set {char}($text_sentence_1_lower + 6) = 32
  set {char}($text_sentence_1_lower + 7) = 119
  set {char}($text_sentence_1_lower + 8) = 111
  set {char}($text_sentence_1_lower + 9) = 114
  set {char}($text_sentence_1_lower + 10) = 108
  set {char}($text_sentence_1_lower + 11) = 100
  set {char}($text_sentence_1_lower + 12) = 46
  set {char}($text_sentence_1_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_1_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=1 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-1-lower.wav
  set $text_sentence_1_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_1_no_period + 0) = 72
  set {char}($text_sentence_1_no_period + 1) = 101
  set {char}($text_sentence_1_no_period + 2) = 108
  set {char}($text_sentence_1_no_period + 3) = 108
  set {char}($text_sentence_1_no_period + 4) = 111
  set {char}($text_sentence_1_no_period + 5) = 32
  set {char}($text_sentence_1_no_period + 6) = 119
  set {char}($text_sentence_1_no_period + 7) = 111
  set {char}($text_sentence_1_no_period + 8) = 114
  set {char}($text_sentence_1_no_period + 9) = 108
  set {char}($text_sentence_1_no_period + 10) = 100
  set {char}($text_sentence_1_no_period + 11) = 46
  set {char}($text_sentence_1_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_1_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=1 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-1-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 199, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=199 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_199_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_199_upper + 0) = 72
  set {char}($text_sentence_199_upper + 1) = 101
  set {char}($text_sentence_199_upper + 2) = 108
  set {char}($text_sentence_199_upper + 3) = 108
  set {char}($text_sentence_199_upper + 4) = 111
  set {char}($text_sentence_199_upper + 5) = 46
  set {char}($text_sentence_199_upper + 6) = 32
  set {char}($text_sentence_199_upper + 7) = 87
  set {char}($text_sentence_199_upper + 8) = 111
  set {char}($text_sentence_199_upper + 9) = 114
  set {char}($text_sentence_199_upper + 10) = 108
  set {char}($text_sentence_199_upper + 11) = 100
  set {char}($text_sentence_199_upper + 12) = 46
  set {char}($text_sentence_199_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_199_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=199 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-199-upper.wav
  set $text_sentence_199_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_199_upper_lower + 0) = 72
  set {char}($text_sentence_199_upper_lower + 1) = 101
  set {char}($text_sentence_199_upper_lower + 2) = 108
  set {char}($text_sentence_199_upper_lower + 3) = 108
  set {char}($text_sentence_199_upper_lower + 4) = 111
  set {char}($text_sentence_199_upper_lower + 5) = 46
  set {char}($text_sentence_199_upper_lower + 6) = 32
  set {char}($text_sentence_199_upper_lower + 7) = 119
  set {char}($text_sentence_199_upper_lower + 8) = 111
  set {char}($text_sentence_199_upper_lower + 9) = 114
  set {char}($text_sentence_199_upper_lower + 10) = 108
  set {char}($text_sentence_199_upper_lower + 11) = 100
  set {char}($text_sentence_199_upper_lower + 12) = 46
  set {char}($text_sentence_199_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_199_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=199 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-199-upper_lower.wav
  set $text_sentence_199_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_199_lower + 0) = 104
  set {char}($text_sentence_199_lower + 1) = 101
  set {char}($text_sentence_199_lower + 2) = 108
  set {char}($text_sentence_199_lower + 3) = 108
  set {char}($text_sentence_199_lower + 4) = 111
  set {char}($text_sentence_199_lower + 5) = 46
  set {char}($text_sentence_199_lower + 6) = 32
  set {char}($text_sentence_199_lower + 7) = 119
  set {char}($text_sentence_199_lower + 8) = 111
  set {char}($text_sentence_199_lower + 9) = 114
  set {char}($text_sentence_199_lower + 10) = 108
  set {char}($text_sentence_199_lower + 11) = 100
  set {char}($text_sentence_199_lower + 12) = 46
  set {char}($text_sentence_199_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_199_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=199 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-199-lower.wav
  set $text_sentence_199_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_199_no_period + 0) = 72
  set {char}($text_sentence_199_no_period + 1) = 101
  set {char}($text_sentence_199_no_period + 2) = 108
  set {char}($text_sentence_199_no_period + 3) = 108
  set {char}($text_sentence_199_no_period + 4) = 111
  set {char}($text_sentence_199_no_period + 5) = 32
  set {char}($text_sentence_199_no_period + 6) = 119
  set {char}($text_sentence_199_no_period + 7) = 111
  set {char}($text_sentence_199_no_period + 8) = 114
  set {char}($text_sentence_199_no_period + 9) = 108
  set {char}($text_sentence_199_no_period + 10) = 100
  set {char}($text_sentence_199_no_period + 11) = 46
  set {char}($text_sentence_199_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_199_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=199 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-199-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 200, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=200 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_200_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_200_upper + 0) = 72
  set {char}($text_sentence_200_upper + 1) = 101
  set {char}($text_sentence_200_upper + 2) = 108
  set {char}($text_sentence_200_upper + 3) = 108
  set {char}($text_sentence_200_upper + 4) = 111
  set {char}($text_sentence_200_upper + 5) = 46
  set {char}($text_sentence_200_upper + 6) = 32
  set {char}($text_sentence_200_upper + 7) = 87
  set {char}($text_sentence_200_upper + 8) = 111
  set {char}($text_sentence_200_upper + 9) = 114
  set {char}($text_sentence_200_upper + 10) = 108
  set {char}($text_sentence_200_upper + 11) = 100
  set {char}($text_sentence_200_upper + 12) = 46
  set {char}($text_sentence_200_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_200_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=200 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-200-upper.wav
  set $text_sentence_200_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_200_upper_lower + 0) = 72
  set {char}($text_sentence_200_upper_lower + 1) = 101
  set {char}($text_sentence_200_upper_lower + 2) = 108
  set {char}($text_sentence_200_upper_lower + 3) = 108
  set {char}($text_sentence_200_upper_lower + 4) = 111
  set {char}($text_sentence_200_upper_lower + 5) = 46
  set {char}($text_sentence_200_upper_lower + 6) = 32
  set {char}($text_sentence_200_upper_lower + 7) = 119
  set {char}($text_sentence_200_upper_lower + 8) = 111
  set {char}($text_sentence_200_upper_lower + 9) = 114
  set {char}($text_sentence_200_upper_lower + 10) = 108
  set {char}($text_sentence_200_upper_lower + 11) = 100
  set {char}($text_sentence_200_upper_lower + 12) = 46
  set {char}($text_sentence_200_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_200_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=200 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-200-upper_lower.wav
  set $text_sentence_200_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_200_lower + 0) = 104
  set {char}($text_sentence_200_lower + 1) = 101
  set {char}($text_sentence_200_lower + 2) = 108
  set {char}($text_sentence_200_lower + 3) = 108
  set {char}($text_sentence_200_lower + 4) = 111
  set {char}($text_sentence_200_lower + 5) = 46
  set {char}($text_sentence_200_lower + 6) = 32
  set {char}($text_sentence_200_lower + 7) = 119
  set {char}($text_sentence_200_lower + 8) = 111
  set {char}($text_sentence_200_lower + 9) = 114
  set {char}($text_sentence_200_lower + 10) = 108
  set {char}($text_sentence_200_lower + 11) = 100
  set {char}($text_sentence_200_lower + 12) = 46
  set {char}($text_sentence_200_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_200_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=200 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-200-lower.wav
  set $text_sentence_200_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_200_no_period + 0) = 72
  set {char}($text_sentence_200_no_period + 1) = 101
  set {char}($text_sentence_200_no_period + 2) = 108
  set {char}($text_sentence_200_no_period + 3) = 108
  set {char}($text_sentence_200_no_period + 4) = 111
  set {char}($text_sentence_200_no_period + 5) = 32
  set {char}($text_sentence_200_no_period + 6) = 119
  set {char}($text_sentence_200_no_period + 7) = 111
  set {char}($text_sentence_200_no_period + 8) = 114
  set {char}($text_sentence_200_no_period + 9) = 108
  set {char}($text_sentence_200_no_period + 10) = 100
  set {char}($text_sentence_200_no_period + 11) = 46
  set {char}($text_sentence_200_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_200_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=200 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-200-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 201, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=201 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_201_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_201_upper + 0) = 72
  set {char}($text_sentence_201_upper + 1) = 101
  set {char}($text_sentence_201_upper + 2) = 108
  set {char}($text_sentence_201_upper + 3) = 108
  set {char}($text_sentence_201_upper + 4) = 111
  set {char}($text_sentence_201_upper + 5) = 46
  set {char}($text_sentence_201_upper + 6) = 32
  set {char}($text_sentence_201_upper + 7) = 87
  set {char}($text_sentence_201_upper + 8) = 111
  set {char}($text_sentence_201_upper + 9) = 114
  set {char}($text_sentence_201_upper + 10) = 108
  set {char}($text_sentence_201_upper + 11) = 100
  set {char}($text_sentence_201_upper + 12) = 46
  set {char}($text_sentence_201_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_201_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=201 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-201-upper.wav
  set $text_sentence_201_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_201_upper_lower + 0) = 72
  set {char}($text_sentence_201_upper_lower + 1) = 101
  set {char}($text_sentence_201_upper_lower + 2) = 108
  set {char}($text_sentence_201_upper_lower + 3) = 108
  set {char}($text_sentence_201_upper_lower + 4) = 111
  set {char}($text_sentence_201_upper_lower + 5) = 46
  set {char}($text_sentence_201_upper_lower + 6) = 32
  set {char}($text_sentence_201_upper_lower + 7) = 119
  set {char}($text_sentence_201_upper_lower + 8) = 111
  set {char}($text_sentence_201_upper_lower + 9) = 114
  set {char}($text_sentence_201_upper_lower + 10) = 108
  set {char}($text_sentence_201_upper_lower + 11) = 100
  set {char}($text_sentence_201_upper_lower + 12) = 46
  set {char}($text_sentence_201_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_201_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=201 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-201-upper_lower.wav
  set $text_sentence_201_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_201_lower + 0) = 104
  set {char}($text_sentence_201_lower + 1) = 101
  set {char}($text_sentence_201_lower + 2) = 108
  set {char}($text_sentence_201_lower + 3) = 108
  set {char}($text_sentence_201_lower + 4) = 111
  set {char}($text_sentence_201_lower + 5) = 46
  set {char}($text_sentence_201_lower + 6) = 32
  set {char}($text_sentence_201_lower + 7) = 119
  set {char}($text_sentence_201_lower + 8) = 111
  set {char}($text_sentence_201_lower + 9) = 114
  set {char}($text_sentence_201_lower + 10) = 108
  set {char}($text_sentence_201_lower + 11) = 100
  set {char}($text_sentence_201_lower + 12) = 46
  set {char}($text_sentence_201_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_201_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=201 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-201-lower.wav
  set $text_sentence_201_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_201_no_period + 0) = 72
  set {char}($text_sentence_201_no_period + 1) = 101
  set {char}($text_sentence_201_no_period + 2) = 108
  set {char}($text_sentence_201_no_period + 3) = 108
  set {char}($text_sentence_201_no_period + 4) = 111
  set {char}($text_sentence_201_no_period + 5) = 32
  set {char}($text_sentence_201_no_period + 6) = 119
  set {char}($text_sentence_201_no_period + 7) = 111
  set {char}($text_sentence_201_no_period + 8) = 114
  set {char}($text_sentence_201_no_period + 9) = 108
  set {char}($text_sentence_201_no_period + 10) = 100
  set {char}($text_sentence_201_no_period + 11) = 46
  set {char}($text_sentence_201_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_201_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=201 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-201-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 250, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=250 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_250_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_250_upper + 0) = 72
  set {char}($text_sentence_250_upper + 1) = 101
  set {char}($text_sentence_250_upper + 2) = 108
  set {char}($text_sentence_250_upper + 3) = 108
  set {char}($text_sentence_250_upper + 4) = 111
  set {char}($text_sentence_250_upper + 5) = 46
  set {char}($text_sentence_250_upper + 6) = 32
  set {char}($text_sentence_250_upper + 7) = 87
  set {char}($text_sentence_250_upper + 8) = 111
  set {char}($text_sentence_250_upper + 9) = 114
  set {char}($text_sentence_250_upper + 10) = 108
  set {char}($text_sentence_250_upper + 11) = 100
  set {char}($text_sentence_250_upper + 12) = 46
  set {char}($text_sentence_250_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_250_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=250 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-250-upper.wav
  set $text_sentence_250_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_250_upper_lower + 0) = 72
  set {char}($text_sentence_250_upper_lower + 1) = 101
  set {char}($text_sentence_250_upper_lower + 2) = 108
  set {char}($text_sentence_250_upper_lower + 3) = 108
  set {char}($text_sentence_250_upper_lower + 4) = 111
  set {char}($text_sentence_250_upper_lower + 5) = 46
  set {char}($text_sentence_250_upper_lower + 6) = 32
  set {char}($text_sentence_250_upper_lower + 7) = 119
  set {char}($text_sentence_250_upper_lower + 8) = 111
  set {char}($text_sentence_250_upper_lower + 9) = 114
  set {char}($text_sentence_250_upper_lower + 10) = 108
  set {char}($text_sentence_250_upper_lower + 11) = 100
  set {char}($text_sentence_250_upper_lower + 12) = 46
  set {char}($text_sentence_250_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_250_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=250 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-250-upper_lower.wav
  set $text_sentence_250_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_250_lower + 0) = 104
  set {char}($text_sentence_250_lower + 1) = 101
  set {char}($text_sentence_250_lower + 2) = 108
  set {char}($text_sentence_250_lower + 3) = 108
  set {char}($text_sentence_250_lower + 4) = 111
  set {char}($text_sentence_250_lower + 5) = 46
  set {char}($text_sentence_250_lower + 6) = 32
  set {char}($text_sentence_250_lower + 7) = 119
  set {char}($text_sentence_250_lower + 8) = 111
  set {char}($text_sentence_250_lower + 9) = 114
  set {char}($text_sentence_250_lower + 10) = 108
  set {char}($text_sentence_250_lower + 11) = 100
  set {char}($text_sentence_250_lower + 12) = 46
  set {char}($text_sentence_250_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_250_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=250 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-250-lower.wav
  set $text_sentence_250_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_250_no_period + 0) = 72
  set {char}($text_sentence_250_no_period + 1) = 101
  set {char}($text_sentence_250_no_period + 2) = 108
  set {char}($text_sentence_250_no_period + 3) = 108
  set {char}($text_sentence_250_no_period + 4) = 111
  set {char}($text_sentence_250_no_period + 5) = 32
  set {char}($text_sentence_250_no_period + 6) = 119
  set {char}($text_sentence_250_no_period + 7) = 111
  set {char}($text_sentence_250_no_period + 8) = 114
  set {char}($text_sentence_250_no_period + 9) = 108
  set {char}($text_sentence_250_no_period + 10) = 100
  set {char}($text_sentence_250_no_period + 11) = 46
  set {char}($text_sentence_250_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_250_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=250 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-250-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 500, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=500 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_500_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_500_upper + 0) = 72
  set {char}($text_sentence_500_upper + 1) = 101
  set {char}($text_sentence_500_upper + 2) = 108
  set {char}($text_sentence_500_upper + 3) = 108
  set {char}($text_sentence_500_upper + 4) = 111
  set {char}($text_sentence_500_upper + 5) = 46
  set {char}($text_sentence_500_upper + 6) = 32
  set {char}($text_sentence_500_upper + 7) = 87
  set {char}($text_sentence_500_upper + 8) = 111
  set {char}($text_sentence_500_upper + 9) = 114
  set {char}($text_sentence_500_upper + 10) = 108
  set {char}($text_sentence_500_upper + 11) = 100
  set {char}($text_sentence_500_upper + 12) = 46
  set {char}($text_sentence_500_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_500_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=500 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-500-upper.wav
  set $text_sentence_500_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_500_upper_lower + 0) = 72
  set {char}($text_sentence_500_upper_lower + 1) = 101
  set {char}($text_sentence_500_upper_lower + 2) = 108
  set {char}($text_sentence_500_upper_lower + 3) = 108
  set {char}($text_sentence_500_upper_lower + 4) = 111
  set {char}($text_sentence_500_upper_lower + 5) = 46
  set {char}($text_sentence_500_upper_lower + 6) = 32
  set {char}($text_sentence_500_upper_lower + 7) = 119
  set {char}($text_sentence_500_upper_lower + 8) = 111
  set {char}($text_sentence_500_upper_lower + 9) = 114
  set {char}($text_sentence_500_upper_lower + 10) = 108
  set {char}($text_sentence_500_upper_lower + 11) = 100
  set {char}($text_sentence_500_upper_lower + 12) = 46
  set {char}($text_sentence_500_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_500_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=500 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-500-upper_lower.wav
  set $text_sentence_500_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_500_lower + 0) = 104
  set {char}($text_sentence_500_lower + 1) = 101
  set {char}($text_sentence_500_lower + 2) = 108
  set {char}($text_sentence_500_lower + 3) = 108
  set {char}($text_sentence_500_lower + 4) = 111
  set {char}($text_sentence_500_lower + 5) = 46
  set {char}($text_sentence_500_lower + 6) = 32
  set {char}($text_sentence_500_lower + 7) = 119
  set {char}($text_sentence_500_lower + 8) = 111
  set {char}($text_sentence_500_lower + 9) = 114
  set {char}($text_sentence_500_lower + 10) = 108
  set {char}($text_sentence_500_lower + 11) = 100
  set {char}($text_sentence_500_lower + 12) = 46
  set {char}($text_sentence_500_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_500_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=500 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-500-lower.wav
  set $text_sentence_500_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_500_no_period + 0) = 72
  set {char}($text_sentence_500_no_period + 1) = 101
  set {char}($text_sentence_500_no_period + 2) = 108
  set {char}($text_sentence_500_no_period + 3) = 108
  set {char}($text_sentence_500_no_period + 4) = 111
  set {char}($text_sentence_500_no_period + 5) = 32
  set {char}($text_sentence_500_no_period + 6) = 119
  set {char}($text_sentence_500_no_period + 7) = 111
  set {char}($text_sentence_500_no_period + 8) = 114
  set {char}($text_sentence_500_no_period + 9) = 108
  set {char}($text_sentence_500_no_period + 10) = 100
  set {char}($text_sentence_500_no_period + 11) = 46
  set {char}($text_sentence_500_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_500_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=500 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-500-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 924, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=924 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_924_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_924_upper + 0) = 72
  set {char}($text_sentence_924_upper + 1) = 101
  set {char}($text_sentence_924_upper + 2) = 108
  set {char}($text_sentence_924_upper + 3) = 108
  set {char}($text_sentence_924_upper + 4) = 111
  set {char}($text_sentence_924_upper + 5) = 46
  set {char}($text_sentence_924_upper + 6) = 32
  set {char}($text_sentence_924_upper + 7) = 87
  set {char}($text_sentence_924_upper + 8) = 111
  set {char}($text_sentence_924_upper + 9) = 114
  set {char}($text_sentence_924_upper + 10) = 108
  set {char}($text_sentence_924_upper + 11) = 100
  set {char}($text_sentence_924_upper + 12) = 46
  set {char}($text_sentence_924_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_924_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=924 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-924-upper.wav
  set $text_sentence_924_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_924_upper_lower + 0) = 72
  set {char}($text_sentence_924_upper_lower + 1) = 101
  set {char}($text_sentence_924_upper_lower + 2) = 108
  set {char}($text_sentence_924_upper_lower + 3) = 108
  set {char}($text_sentence_924_upper_lower + 4) = 111
  set {char}($text_sentence_924_upper_lower + 5) = 46
  set {char}($text_sentence_924_upper_lower + 6) = 32
  set {char}($text_sentence_924_upper_lower + 7) = 119
  set {char}($text_sentence_924_upper_lower + 8) = 111
  set {char}($text_sentence_924_upper_lower + 9) = 114
  set {char}($text_sentence_924_upper_lower + 10) = 108
  set {char}($text_sentence_924_upper_lower + 11) = 100
  set {char}($text_sentence_924_upper_lower + 12) = 46
  set {char}($text_sentence_924_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_924_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=924 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-924-upper_lower.wav
  set $text_sentence_924_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_924_lower + 0) = 104
  set {char}($text_sentence_924_lower + 1) = 101
  set {char}($text_sentence_924_lower + 2) = 108
  set {char}($text_sentence_924_lower + 3) = 108
  set {char}($text_sentence_924_lower + 4) = 111
  set {char}($text_sentence_924_lower + 5) = 46
  set {char}($text_sentence_924_lower + 6) = 32
  set {char}($text_sentence_924_lower + 7) = 119
  set {char}($text_sentence_924_lower + 8) = 111
  set {char}($text_sentence_924_lower + 9) = 114
  set {char}($text_sentence_924_lower + 10) = 108
  set {char}($text_sentence_924_lower + 11) = 100
  set {char}($text_sentence_924_lower + 12) = 46
  set {char}($text_sentence_924_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_924_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=924 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-924-lower.wav
  set $text_sentence_924_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_924_no_period + 0) = 72
  set {char}($text_sentence_924_no_period + 1) = 101
  set {char}($text_sentence_924_no_period + 2) = 108
  set {char}($text_sentence_924_no_period + 3) = 108
  set {char}($text_sentence_924_no_period + 4) = 111
  set {char}($text_sentence_924_no_period + 5) = 32
  set {char}($text_sentence_924_no_period + 6) = 119
  set {char}($text_sentence_924_no_period + 7) = 111
  set {char}($text_sentence_924_no_period + 8) = 114
  set {char}($text_sentence_924_no_period + 9) = 108
  set {char}($text_sentence_924_no_period + 10) = 100
  set {char}($text_sentence_924_no_period + 11) = 46
  set {char}($text_sentence_924_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_924_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=924 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-924-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=925 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_925_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_925_upper + 0) = 72
  set {char}($text_sentence_925_upper + 1) = 101
  set {char}($text_sentence_925_upper + 2) = 108
  set {char}($text_sentence_925_upper + 3) = 108
  set {char}($text_sentence_925_upper + 4) = 111
  set {char}($text_sentence_925_upper + 5) = 46
  set {char}($text_sentence_925_upper + 6) = 32
  set {char}($text_sentence_925_upper + 7) = 87
  set {char}($text_sentence_925_upper + 8) = 111
  set {char}($text_sentence_925_upper + 9) = 114
  set {char}($text_sentence_925_upper + 10) = 108
  set {char}($text_sentence_925_upper + 11) = 100
  set {char}($text_sentence_925_upper + 12) = 46
  set {char}($text_sentence_925_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_925_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=925 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-925-upper.wav
  set $text_sentence_925_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_925_upper_lower + 0) = 72
  set {char}($text_sentence_925_upper_lower + 1) = 101
  set {char}($text_sentence_925_upper_lower + 2) = 108
  set {char}($text_sentence_925_upper_lower + 3) = 108
  set {char}($text_sentence_925_upper_lower + 4) = 111
  set {char}($text_sentence_925_upper_lower + 5) = 46
  set {char}($text_sentence_925_upper_lower + 6) = 32
  set {char}($text_sentence_925_upper_lower + 7) = 119
  set {char}($text_sentence_925_upper_lower + 8) = 111
  set {char}($text_sentence_925_upper_lower + 9) = 114
  set {char}($text_sentence_925_upper_lower + 10) = 108
  set {char}($text_sentence_925_upper_lower + 11) = 100
  set {char}($text_sentence_925_upper_lower + 12) = 46
  set {char}($text_sentence_925_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_925_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=925 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-925-upper_lower.wav
  set $text_sentence_925_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_925_lower + 0) = 104
  set {char}($text_sentence_925_lower + 1) = 101
  set {char}($text_sentence_925_lower + 2) = 108
  set {char}($text_sentence_925_lower + 3) = 108
  set {char}($text_sentence_925_lower + 4) = 111
  set {char}($text_sentence_925_lower + 5) = 46
  set {char}($text_sentence_925_lower + 6) = 32
  set {char}($text_sentence_925_lower + 7) = 119
  set {char}($text_sentence_925_lower + 8) = 111
  set {char}($text_sentence_925_lower + 9) = 114
  set {char}($text_sentence_925_lower + 10) = 108
  set {char}($text_sentence_925_lower + 11) = 100
  set {char}($text_sentence_925_lower + 12) = 46
  set {char}($text_sentence_925_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_925_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=925 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-925-lower.wav
  set $text_sentence_925_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_925_no_period + 0) = 72
  set {char}($text_sentence_925_no_period + 1) = 101
  set {char}($text_sentence_925_no_period + 2) = 108
  set {char}($text_sentence_925_no_period + 3) = 108
  set {char}($text_sentence_925_no_period + 4) = 111
  set {char}($text_sentence_925_no_period + 5) = 32
  set {char}($text_sentence_925_no_period + 6) = 119
  set {char}($text_sentence_925_no_period + 7) = 111
  set {char}($text_sentence_925_no_period + 8) = 114
  set {char}($text_sentence_925_no_period + 9) = 108
  set {char}($text_sentence_925_no_period + 10) = 100
  set {char}($text_sentence_925_no_period + 11) = 46
  set {char}($text_sentence_925_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_925_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=925 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-925-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 926, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=926 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_926_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_926_upper + 0) = 72
  set {char}($text_sentence_926_upper + 1) = 101
  set {char}($text_sentence_926_upper + 2) = 108
  set {char}($text_sentence_926_upper + 3) = 108
  set {char}($text_sentence_926_upper + 4) = 111
  set {char}($text_sentence_926_upper + 5) = 46
  set {char}($text_sentence_926_upper + 6) = 32
  set {char}($text_sentence_926_upper + 7) = 87
  set {char}($text_sentence_926_upper + 8) = 111
  set {char}($text_sentence_926_upper + 9) = 114
  set {char}($text_sentence_926_upper + 10) = 108
  set {char}($text_sentence_926_upper + 11) = 100
  set {char}($text_sentence_926_upper + 12) = 46
  set {char}($text_sentence_926_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_926_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=926 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-926-upper.wav
  set $text_sentence_926_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_926_upper_lower + 0) = 72
  set {char}($text_sentence_926_upper_lower + 1) = 101
  set {char}($text_sentence_926_upper_lower + 2) = 108
  set {char}($text_sentence_926_upper_lower + 3) = 108
  set {char}($text_sentence_926_upper_lower + 4) = 111
  set {char}($text_sentence_926_upper_lower + 5) = 46
  set {char}($text_sentence_926_upper_lower + 6) = 32
  set {char}($text_sentence_926_upper_lower + 7) = 119
  set {char}($text_sentence_926_upper_lower + 8) = 111
  set {char}($text_sentence_926_upper_lower + 9) = 114
  set {char}($text_sentence_926_upper_lower + 10) = 108
  set {char}($text_sentence_926_upper_lower + 11) = 100
  set {char}($text_sentence_926_upper_lower + 12) = 46
  set {char}($text_sentence_926_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_926_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=926 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-926-upper_lower.wav
  set $text_sentence_926_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_926_lower + 0) = 104
  set {char}($text_sentence_926_lower + 1) = 101
  set {char}($text_sentence_926_lower + 2) = 108
  set {char}($text_sentence_926_lower + 3) = 108
  set {char}($text_sentence_926_lower + 4) = 111
  set {char}($text_sentence_926_lower + 5) = 46
  set {char}($text_sentence_926_lower + 6) = 32
  set {char}($text_sentence_926_lower + 7) = 119
  set {char}($text_sentence_926_lower + 8) = 111
  set {char}($text_sentence_926_lower + 9) = 114
  set {char}($text_sentence_926_lower + 10) = 108
  set {char}($text_sentence_926_lower + 11) = 100
  set {char}($text_sentence_926_lower + 12) = 46
  set {char}($text_sentence_926_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_926_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=926 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-926-lower.wav
  set $text_sentence_926_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_926_no_period + 0) = 72
  set {char}($text_sentence_926_no_period + 1) = 101
  set {char}($text_sentence_926_no_period + 2) = 108
  set {char}($text_sentence_926_no_period + 3) = 108
  set {char}($text_sentence_926_no_period + 4) = 111
  set {char}($text_sentence_926_no_period + 5) = 32
  set {char}($text_sentence_926_no_period + 6) = 119
  set {char}($text_sentence_926_no_period + 7) = 111
  set {char}($text_sentence_926_no_period + 8) = 114
  set {char}($text_sentence_926_no_period + 9) = 108
  set {char}($text_sentence_926_no_period + 10) = 100
  set {char}($text_sentence_926_no_period + 11) = 46
  set {char}($text_sentence_926_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_926_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=926 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-926-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 65534, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=65534 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_65534_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65534_upper + 0) = 72
  set {char}($text_sentence_65534_upper + 1) = 101
  set {char}($text_sentence_65534_upper + 2) = 108
  set {char}($text_sentence_65534_upper + 3) = 108
  set {char}($text_sentence_65534_upper + 4) = 111
  set {char}($text_sentence_65534_upper + 5) = 46
  set {char}($text_sentence_65534_upper + 6) = 32
  set {char}($text_sentence_65534_upper + 7) = 87
  set {char}($text_sentence_65534_upper + 8) = 111
  set {char}($text_sentence_65534_upper + 9) = 114
  set {char}($text_sentence_65534_upper + 10) = 108
  set {char}($text_sentence_65534_upper + 11) = 100
  set {char}($text_sentence_65534_upper + 12) = 46
  set {char}($text_sentence_65534_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65534_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65534 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65534-upper.wav
  set $text_sentence_65534_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65534_upper_lower + 0) = 72
  set {char}($text_sentence_65534_upper_lower + 1) = 101
  set {char}($text_sentence_65534_upper_lower + 2) = 108
  set {char}($text_sentence_65534_upper_lower + 3) = 108
  set {char}($text_sentence_65534_upper_lower + 4) = 111
  set {char}($text_sentence_65534_upper_lower + 5) = 46
  set {char}($text_sentence_65534_upper_lower + 6) = 32
  set {char}($text_sentence_65534_upper_lower + 7) = 119
  set {char}($text_sentence_65534_upper_lower + 8) = 111
  set {char}($text_sentence_65534_upper_lower + 9) = 114
  set {char}($text_sentence_65534_upper_lower + 10) = 108
  set {char}($text_sentence_65534_upper_lower + 11) = 100
  set {char}($text_sentence_65534_upper_lower + 12) = 46
  set {char}($text_sentence_65534_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65534_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65534 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65534-upper_lower.wav
  set $text_sentence_65534_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65534_lower + 0) = 104
  set {char}($text_sentence_65534_lower + 1) = 101
  set {char}($text_sentence_65534_lower + 2) = 108
  set {char}($text_sentence_65534_lower + 3) = 108
  set {char}($text_sentence_65534_lower + 4) = 111
  set {char}($text_sentence_65534_lower + 5) = 46
  set {char}($text_sentence_65534_lower + 6) = 32
  set {char}($text_sentence_65534_lower + 7) = 119
  set {char}($text_sentence_65534_lower + 8) = 111
  set {char}($text_sentence_65534_lower + 9) = 114
  set {char}($text_sentence_65534_lower + 10) = 108
  set {char}($text_sentence_65534_lower + 11) = 100
  set {char}($text_sentence_65534_lower + 12) = 46
  set {char}($text_sentence_65534_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65534_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65534 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65534-lower.wav
  set $text_sentence_65534_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65534_no_period + 0) = 72
  set {char}($text_sentence_65534_no_period + 1) = 101
  set {char}($text_sentence_65534_no_period + 2) = 108
  set {char}($text_sentence_65534_no_period + 3) = 108
  set {char}($text_sentence_65534_no_period + 4) = 111
  set {char}($text_sentence_65534_no_period + 5) = 32
  set {char}($text_sentence_65534_no_period + 6) = 119
  set {char}($text_sentence_65534_no_period + 7) = 111
  set {char}($text_sentence_65534_no_period + 8) = 114
  set {char}($text_sentence_65534_no_period + 9) = 108
  set {char}($text_sentence_65534_no_period + 10) = 100
  set {char}($text_sentence_65534_no_period + 11) = 46
  set {char}($text_sentence_65534_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65534_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65534 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65534-no_period.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 65535, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set $get_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out, $out + 1, $out + 2, $out + 3, $speaker)
  printf "SENT_PAUSE value=65535 getter_ret=%d pitch=%d speed=%d volume=%d getter_value=%d\n", $get_result, *$out, *($out + 1), *($out + 2), *($out + 3)
  set $text_sentence_65535_upper = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65535_upper + 0) = 72
  set {char}($text_sentence_65535_upper + 1) = 101
  set {char}($text_sentence_65535_upper + 2) = 108
  set {char}($text_sentence_65535_upper + 3) = 108
  set {char}($text_sentence_65535_upper + 4) = 111
  set {char}($text_sentence_65535_upper + 5) = 46
  set {char}($text_sentence_65535_upper + 6) = 32
  set {char}($text_sentence_65535_upper + 7) = 87
  set {char}($text_sentence_65535_upper + 8) = 111
  set {char}($text_sentence_65535_upper + 9) = 114
  set {char}($text_sentence_65535_upper + 10) = 108
  set {char}($text_sentence_65535_upper + 11) = 100
  set {char}($text_sentence_65535_upper + 12) = 46
  set {char}($text_sentence_65535_upper + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65535_upper, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65535 context=upper synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65535-upper.wav
  set $text_sentence_65535_upper_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65535_upper_lower + 0) = 72
  set {char}($text_sentence_65535_upper_lower + 1) = 101
  set {char}($text_sentence_65535_upper_lower + 2) = 108
  set {char}($text_sentence_65535_upper_lower + 3) = 108
  set {char}($text_sentence_65535_upper_lower + 4) = 111
  set {char}($text_sentence_65535_upper_lower + 5) = 46
  set {char}($text_sentence_65535_upper_lower + 6) = 32
  set {char}($text_sentence_65535_upper_lower + 7) = 119
  set {char}($text_sentence_65535_upper_lower + 8) = 111
  set {char}($text_sentence_65535_upper_lower + 9) = 114
  set {char}($text_sentence_65535_upper_lower + 10) = 108
  set {char}($text_sentence_65535_upper_lower + 11) = 100
  set {char}($text_sentence_65535_upper_lower + 12) = 46
  set {char}($text_sentence_65535_upper_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65535_upper_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65535 context=upper_lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65535-upper_lower.wav
  set $text_sentence_65535_lower = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65535_lower + 0) = 104
  set {char}($text_sentence_65535_lower + 1) = 101
  set {char}($text_sentence_65535_lower + 2) = 108
  set {char}($text_sentence_65535_lower + 3) = 108
  set {char}($text_sentence_65535_lower + 4) = 111
  set {char}($text_sentence_65535_lower + 5) = 46
  set {char}($text_sentence_65535_lower + 6) = 32
  set {char}($text_sentence_65535_lower + 7) = 119
  set {char}($text_sentence_65535_lower + 8) = 111
  set {char}($text_sentence_65535_lower + 9) = 114
  set {char}($text_sentence_65535_lower + 10) = 108
  set {char}($text_sentence_65535_lower + 11) = 100
  set {char}($text_sentence_65535_lower + 12) = 46
  set {char}($text_sentence_65535_lower + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65535_lower, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65535 context=lower synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65535-lower.wav
  set $text_sentence_65535_no_period = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}($text_sentence_65535_no_period + 0) = 72
  set {char}($text_sentence_65535_no_period + 1) = 101
  set {char}($text_sentence_65535_no_period + 2) = 108
  set {char}($text_sentence_65535_no_period + 3) = 108
  set {char}($text_sentence_65535_no_period + 4) = 111
  set {char}($text_sentence_65535_no_period + 5) = 32
  set {char}($text_sentence_65535_no_period + 6) = 119
  set {char}($text_sentence_65535_no_period + 7) = 111
  set {char}($text_sentence_65535_no_period + 8) = 114
  set {char}($text_sentence_65535_no_period + 9) = 108
  set {char}($text_sentence_65535_no_period + 10) = 100
  set {char}($text_sentence_65535_no_period + 11) = 46
  set {char}($text_sentence_65535_no_period + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text_sentence_65535_no_period, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "SENT_PAUSE value=65535 context=no_period synth_ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/sentence-pause-65535-no_period.wav
  kill
  quit
end

continue
