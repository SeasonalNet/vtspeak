set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001cd84
commands 1
  silent
  printf "PREPROCESS_RNG gettickcount=%u\n", $eax
  disable 1
  continue
end

break *0x1001da50
commands 2
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 2
  set $path = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {unsigned char}($path + 0) = 90
  set {unsigned char}($path + 1) = 58
  set {unsigned char}($path + 2) = 47
  set {unsigned char}($path + 3) = 119
  set {unsigned char}($path + 4) = 111
  set {unsigned char}($path + 5) = 114
  set {unsigned char}($path + 6) = 107
  set {unsigned char}($path + 7) = 47
  set {unsigned char}($path + 8) = 115
  set {unsigned char}($path + 9) = 116
  set {unsigned char}($path + 10) = 97
  set {unsigned char}($path + 11) = 103
  set {unsigned char}($path + 12) = 101
  set {unsigned char}($path + 13) = 50
  set {unsigned char}($path + 14) = 49
  set {unsigned char}($path + 15) = 47
  set {unsigned char}($path + 16) = 114
  set {unsigned char}($path + 17) = 49
  set {unsigned char}($path + 18) = 0
  printf "PREPROCESS_RNG_INIT flag=%u state=%u\n", *(unsigned int *)0x1007cfb0, *(unsigned int *)0x1009fc4c
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_CALL n=1 raw=%#x state=%u\n", $raw, *(unsigned int *)0x1009fc4c
  set {unsigned char}($path + 17) = 50
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_CALL n=2 raw=%#x state=%u\n", $raw, *(unsigned int *)0x1009fc4c
  set {unsigned char}($path + 17) = 51
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_CALL n=3 raw=%#x state=%u\n", $raw, *(unsigned int *)0x1009fc4c
  continue
end

continue
