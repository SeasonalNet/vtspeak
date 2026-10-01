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
  set $text = (char *)malloc(64)
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 0, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=-1 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0--1-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=-1 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0--1-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=-1 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0--1-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=0 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0-0-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=0 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0-0-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=0 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0-0-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=250 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0-250-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=250 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0-250-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=0 call=250 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-0-250-control.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 925, $speaker)
  call ((void (*)(int, int))0x100281b0)(200, $speaker)
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=-1 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925--1-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=-1 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925--1-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=-1 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925--1-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=0 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925-0-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=0 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925-0-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=0 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925-0-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=250 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925-250-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=250 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925-250-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=sentence stored=925 call=250 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-sentence-925-250-control.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 200, $speaker)
  call ((void (*)(int, int))0x100281b0)(0, $speaker)
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=-1 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0--1-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=-1 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0--1-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=-1 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0--1-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=0 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0-0-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=0 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0-0-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=0 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0-0-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=250 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0-250-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=250 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0-250-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=0 call=250 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-0-250-control.wav
  call ((void (*)(int, int, int, int, int))0x10027fe0)(100, 100, 200, 200, $speaker)
  call ((void (*)(int, int))0x100281b0)(925, $speaker)
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=-1 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925--1-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=-1 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925--1-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=-1 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925--1-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=0 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925-0-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=0 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925-0-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 0, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=0 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925-0-control.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 46
  set {char}($text + 6) = 32
  set {char}($text + 7) = 87
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=250 context=period ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925-250-period.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 44
  set {char}($text + 6) = 32
  set {char}($text + 7) = 119
  set {char}($text + 8) = 111
  set {char}($text + 9) = 114
  set {char}($text + 10) = 108
  set {char}($text + 11) = 100
  set {char}($text + 12) = 46
  set {char}($text + 13) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=250 context=comma ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925-250-comma.wav
  set {char}($text + 0) = 72
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 119
  set {char}($text + 7) = 111
  set {char}($text + 8) = 114
  set {char}($text + 9) = 108
  set {char}($text + 10) = 100
  set {char}($text + 11) = 46
  set {char}($text + 12) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($fmt, $text, $path, $speaker, -1, -1, -1, 250, -1, -1)
  printf "PAUSE_PRECEDENCE axis=comma stored=925 call=250 context=control ret=%d\n", $result
  shell cp /work/stage21/sandbox/stage5/output.wav /work/stage21/pause-precedence-comma-925-250-control.wav
  kill
  quit
end

continue
