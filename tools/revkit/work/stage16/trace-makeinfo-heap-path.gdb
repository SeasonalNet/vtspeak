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
  set {char}($path + 21) = 109
  set {char}($path + 22) = 97
  set {char}($path + 23) = 107
  set {char}($path + 24) = 101
  set {char}($path + 25) = 105
  set {char}($path + 26) = 110
  set {char}($path + 27) = 102
  set {char}($path + 28) = 111
  set {char}($path + 29) = 45
  set {char}($path + 30) = 112
  set {char}($path + 31) = 97
  set {char}($path + 32) = 116
  set {char}($path + 33) = 104
  set {char}($path + 34) = 45
  set {char}($path + 35) = 50
  set {char}($path + 36) = 48
  set {char}($path + 37) = 50
  set {char}($path + 38) = 54
  set {char}($path + 39) = 48
  set {char}($path + 40) = 57
  set {char}($path + 41) = 50
  set {char}($path + 42) = 54
  set {char}($path + 43) = 0
  printf "MAKEINFO_HEAP_PATH_BEFORE ptr=%#x value=", $path
  x/s $path
  set $raw = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de90)((char *)$text, $path, 1, -1, -1, -1, -1, -1, -1)
  printf "MAKEINFO_HEAP_PATH_AFTER raw_eax=%#x signed=%d value=", $raw, $raw
  x/s $path
  continue
end

continue
