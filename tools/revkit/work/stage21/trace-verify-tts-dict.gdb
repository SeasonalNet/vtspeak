set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $path = ((char *(*)(unsigned int))0x1001d9c0)(41)
  set {char}($path + 0) = 90
  set {char}($path + 1) = 58
  set {char}($path + 2) = 47
  set {char}($path + 3) = 119
  set {char}($path + 4) = 111
  set {char}($path + 5) = 114
  set {char}($path + 6) = 107
  set {char}($path + 7) = 47
  set {char}($path + 8) = 115
  set {char}($path + 9) = 116
  set {char}($path + 10) = 97
  set {char}($path + 11) = 103
  set {char}($path + 12) = 101
  set {char}($path + 13) = 49
  set {char}($path + 14) = 54
  set {char}($path + 15) = 47
  set {char}($path + 16) = 117
  set {char}($path + 17) = 115
  set {char}($path + 18) = 101
  set {char}($path + 19) = 114
  set {char}($path + 20) = 100
  set {char}($path + 21) = 105
  set {char}($path + 22) = 99
  set {char}($path + 23) = 116
  set {char}($path + 24) = 45
  set {char}($path + 25) = 104
  set {char}($path + 26) = 101
  set {char}($path + 27) = 108
  set {char}($path + 28) = 108
  set {char}($path + 29) = 111
  set {char}($path + 30) = 45
  set {char}($path + 31) = 112
  set {char}($path + 32) = 108
  set {char}($path + 33) = 97
  set {char}($path + 34) = 105
  set {char}($path + 35) = 110
  set {char}($path + 36) = 46
  set {char}($path + 37) = 99
  set {char}($path + 38) = 115
  set {char}($path + 39) = 118
  set {char}($path + 40) = 0
  set $gate_original = *(unsigned char *)0x100a7489
  set $load = ((short (*)(int, char *))0x10027960)(27, $path)
  printf "VERIFY_TTS_DICT_LOAD index=27 ax=%d ptr=%#x gate=%u\n", $load, *(unsigned int *)(0x100a647c + 27 * 4), $gate_original
  set {unsigned char}0x100a7489 = 0
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, -1, 0)
  printf "VERIFY_TTS_DICT state=default text=match_lower dict=-1 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 119
  set {char}($text + 1) = 111
  set {char}($text + 2) = 114
  set {char}($text + 3) = 108
  set {char}($text + 4) = 100
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, -1, 0)
  printf "VERIFY_TTS_DICT state=default text=miss dict=-1 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 72
  set {char}($text + 1) = 69
  set {char}($text + 2) = 76
  set {char}($text + 3) = 76
  set {char}($text + 4) = 79
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, -1, 0)
  printf "VERIFY_TTS_DICT state=default text=case_variant dict=-1 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 104
  set {char}($text + 7) = 101
  set {char}($text + 8) = 108
  set {char}($text + 9) = 108
  set {char}($text + 10) = 111
  set {char}($text + 11) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, -1, 0)
  printf "VERIFY_TTS_DICT state=default text=repeat dict=-1 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 60
  set {char}($text + 6) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, -1, 0)
  printf "VERIFY_TTS_DICT state=default text=markup_tail dict=-1 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($text + 0) = 60
  set {char}($text + 1) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, -1, 0)
  printf "VERIFY_TTS_DICT state=default text=malformed_lt dict=-1 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_off text=match_lower dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 119
  set {char}($text + 1) = 111
  set {char}($text + 2) = 114
  set {char}($text + 3) = 108
  set {char}($text + 4) = 100
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_off text=miss dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 72
  set {char}($text + 1) = 69
  set {char}($text + 2) = 76
  set {char}($text + 3) = 76
  set {char}($text + 4) = 79
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_off text=case_variant dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 104
  set {char}($text + 7) = 101
  set {char}($text + 8) = 108
  set {char}($text + 9) = 108
  set {char}($text + 10) = 111
  set {char}($text + 11) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_off text=repeat dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 60
  set {char}($text + 6) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_off text=markup_tail dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($text + 0) = 60
  set {char}($text + 1) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_off text=malformed_lt dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set {unsigned char}0x100a7489 = 1
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_on text=match_lower dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 119
  set {char}($text + 1) = 111
  set {char}($text + 2) = 114
  set {char}($text + 3) = 108
  set {char}($text + 4) = 100
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_on text=miss dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 72
  set {char}($text + 1) = 69
  set {char}($text + 2) = 76
  set {char}($text + 3) = 76
  set {char}($text + 4) = 79
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_on text=case_variant dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 32
  set {char}($text + 6) = 104
  set {char}($text + 7) = 101
  set {char}($text + 8) = 108
  set {char}($text + 9) = 108
  set {char}($text + 10) = 111
  set {char}($text + 11) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_on text=repeat dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(7)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 60
  set {char}($text + 6) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_on text=markup_tail dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($text + 0) = 60
  set {char}($text + 1) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 27, 0)
  printf "VERIFY_TTS_DICT state=dict_gate_on text=malformed_lt dict=27 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 28, 0)
  printf "VERIFY_TTS_DICT state=empty_slot_gate_on text=match_lower dict=28 ax=%u eax=%d\n", $r & 65535, $r
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 119
  set {char}($text + 1) = 111
  set {char}($text + 2) = 114
  set {char}($text + 3) = 108
  set {char}($text + 4) = 100
  set {char}($text + 5) = 0
  set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($text, 1, 28, 0)
  printf "VERIFY_TTS_DICT state=empty_slot_gate_on text=miss dict=28 ax=%u eax=%d\n", $r & 65535, $r
  set {unsigned char}0x100a7489 = $gate_original
  set $unload = ((short (*)(int))0x10027a80)(27)
  printf "VERIFY_TTS_DICT_UNLOAD index=27 ax=%d restored_gate=%u\n", $unload, *(unsigned char *)0x100a7489
  continue
end

continue
