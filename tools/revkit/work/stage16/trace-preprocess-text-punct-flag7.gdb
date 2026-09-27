set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $path = ((char *(*)(unsigned int))0x1001d9c0)(128)
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
  set {char}($path + 16) = 104
  set {char}($path + 17) = 101
  set {char}($path + 18) = 97
  set {char}($path + 19) = 112
  set {char}($path + 20) = 45
  set {char}($path + 21) = 102
  set {char}($path + 22) = 108
  set {char}($path + 23) = 97
  set {char}($path + 24) = 103
  set {char}($path + 25) = 55
  set {char}($path + 26) = 45
  set {char}($path + 27) = 99
  set {char}($path + 28) = 97
  set {char}($path + 29) = 115
  set {char}($path + 30) = 101
  set {char}($path + 31) = 45
  set {char}($path + 32) = 112
  set {char}($path + 33) = 117
  set {char}($path + 34) = 110
  set {char}($path + 35) = 99
  set {char}($path + 36) = 116
  set {char}($path + 37) = 45
  set {char}($path + 38) = 50
  set {char}($path + 39) = 48
  set {char}($path + 40) = 50
  set {char}($path + 41) = 54
  set {char}($path + 42) = 48
  set {char}($path + 43) = 57
  set {char}($path + 44) = 50
  set {char}($path + 45) = 54
  set {char}($path + 46) = 0
  printf "PREPROCESS_CASE_BEFORE case=punct flag=7 path="
  x/s $path
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_CASE_RESULT case=punct flag=7 raw_eax=%#x low_ax=%d path=%s\n", $raw, ((short)$raw), $path
  continue
end

continue
