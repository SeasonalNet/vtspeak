set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {char}$p = 60
  set {char}($p + 1) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=lt ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 62
  set {char}($p + 1) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=gt ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 60
  set {char}($p + 1) = 62
  set {char}($p + 2) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=empty_pair ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 60
  set {char}($p + 1) = 47
  set {char}($p + 2) = 62
  set {char}($p + 3) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=empty_close ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 60
  set {char}($p + 1) = 65
  set {char}($p + 2) = 62
  set {char}($p + 3) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=single_tag ax=%u eax=%u\n", $r & 65535, $r
  set {char}($p + 2) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=unterminated_open ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 65
  set {char}($p + 1) = 62
  set {char}($p + 2) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=bare_gt ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 60
  set {char}($p + 1) = 60
  set {char}($p + 2) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=double_lt ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 62
  set {char}($p + 1) = 62
  set {char}($p + 2) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=double_gt ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 60
  set {char}($p + 1) = 65
  set {char}($p + 2) = 47
  set {char}($p + 3) = 62
  set {char}($p + 4) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=open_then_close ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 60
  set {char}($p + 1) = 47
  set {char}($p + 2) = 65
  set {char}($p + 3) = 62
  set {char}($p + 4) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=close_only ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 60
  set {char}($p + 1) = 65
  set {char}($p + 2) = 62
  set {char}($p + 3) = 60
  set {char}($p + 4) = 47
  set {char}($p + 5) = 65
  set {char}($p + 6) = 62
  set {char}($p + 7) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=balanced_pair ax=%u eax=%u\n", $r & 65535, $r
  set {char}($p + 5) = 66
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=mismatched_pair ax=%u eax=%u\n", $r & 65535, $r
  set {char}$p = 65
  set {char}($p + 1) = 60
  set {char}($p + 2) = 66
  set {char}($p + 3) = 62
  set {char}($p + 4) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=plain_then_open ax=%u eax=%u\n", $r & 65535, $r
  set {char}($p + 1) = 60
  set {char}($p + 2) = 47
  set {char}($p + 3) = 66
  set {char}($p + 4) = 62
  set {char}($p + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
  printf "VERIFY_TTS_MARKUP id=plain_then_close ax=%u eax=%u\n", $r & 65535, $r
  continue
end

continue
