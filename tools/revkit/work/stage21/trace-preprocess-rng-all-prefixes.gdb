set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
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
  set {unsigned char}($path + 16) = 112
  set {unsigned char}($path + 17) = 48
  set {unsigned char}($path + 18) = 0
  set {unsigned char}($path + 18) = 0

  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 8
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=0 seed=8 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  set {unsigned char}($path + 17) = 49
  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 7
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=1 seed=7 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  set {unsigned char}($path + 17) = 50
  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 6
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=2 seed=6 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  set {unsigned char}($path + 17) = 51
  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 5
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=3 seed=5 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  set {unsigned char}($path + 17) = 52
  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 4
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=4 seed=4 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  set {unsigned char}($path + 17) = 53
  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 3
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=5 seed=3 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  set {unsigned char}($path + 17) = 54
  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 2
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=6 seed=2 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  set {unsigned char}($path + 17) = 55
  set {unsigned int}0x1007cfb0 = 0
  set {unsigned int}0x1009fc4c = 1
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_RNG_PREFIX index=7 seed=1 state=%u raw=%#x path=%s\n", *(unsigned int *)0x1009fc4c, $raw, $path

  continue
end

continue
