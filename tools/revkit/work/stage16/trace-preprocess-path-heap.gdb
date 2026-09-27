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
  set {char}($path + 16) = 112
  set {char}($path + 17) = 114
  set {char}($path + 18) = 101
  set {char}($path + 19) = 112
  set {char}($path + 20) = 114
  set {char}($path + 21) = 111
  set {char}($path + 22) = 99
  set {char}($path + 23) = 101
  set {char}($path + 24) = 115
  set {char}($path + 25) = 115
  set {char}($path + 26) = 45
  set {char}($path + 27) = 104
  set {char}($path + 28) = 101
  set {char}($path + 29) = 97
  set {char}($path + 30) = 112
  set {char}($path + 31) = 45
  set {char}($path + 32) = 98
  set {char}($path + 33) = 121
  set {char}($path + 34) = 116
  set {char}($path + 35) = 101
  set {char}($path + 36) = 115
  set {char}($path + 37) = 45
  set {char}($path + 38) = 112
  set {char}($path + 39) = 97
  set {char}($path + 40) = 116
  set {char}($path + 41) = 104
  set {char}($path + 42) = 45
  set {char}($path + 43) = 50
  set {char}($path + 44) = 48
  set {char}($path + 45) = 50
  set {char}($path + 46) = 54
  set {char}($path + 47) = 48
  set {char}($path + 48) = 57
  set {char}($path + 49) = 50
  set {char}($path + 50) = 54
  set {char}($path + 51) = 0
  printf "PREPROCESS_HEAP_PATH_BEFORE ptr=%#x value=", $path
  x/s $path
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 5, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_HEAP_PATH_AFTER raw_eax=%#x low_ax=%d value=", $raw, ((short)$raw)
  x/s $path
  continue
end

continue
