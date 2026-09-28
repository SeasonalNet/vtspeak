set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $name = (char *)malloc(128)
  set {char}($name + 0) = 90
  set {char}($name + 1) = 58
  set {char}($name + 2) = 47
  set {char}($name + 3) = 119
  set {char}($name + 4) = 111
  set {char}($name + 5) = 114
  set {char}($name + 6) = 107
  set {char}($name + 7) = 47
  set {char}($name + 8) = 115
  set {char}($name + 9) = 116
  set {char}($name + 10) = 97
  set {char}($name + 11) = 103
  set {char}($name + 12) = 101
  set {char}($name + 13) = 49
  set {char}($name + 14) = 54
  set {char}($name + 15) = 47
  set {char}($name + 16) = 108
  set {char}($name + 17) = 105
  set {char}($name + 18) = 112
  set {char}($name + 19) = 115
  set {char}($name + 20) = 121
  set {char}($name + 21) = 110
  set {char}($name + 22) = 99
  set {char}($name + 23) = 45
  set {char}($name + 24) = 102
  set {char}($name + 25) = 105
  set {char}($name + 26) = 108
  set {char}($name + 27) = 101
  set {char}($name + 28) = 110
  set {char}($name + 29) = 97
  set {char}($name + 30) = 109
  set {char}($name + 31) = 101
  set {char}($name + 32) = 45
  set {char}($name + 33) = 98
  set {char}($name + 34) = 121
  set {char}($name + 35) = 116
  set {char}($name + 36) = 101
  set {char}($name + 37) = 119
  set {char}($name + 38) = 105
  set {char}($name + 39) = 115
  set {char}($name + 40) = 101
  set {char}($name + 41) = 45
  set {char}($name + 42) = 97
  set {char}($name + 43) = 98
  set {char}($name + 44) = 115
  set {char}($name + 45) = 111
  set {char}($name + 46) = 108
  set {char}($name + 47) = 117
  set {char}($name + 48) = 116
  set {char}($name + 49) = 101
  set {char}($name + 50) = 45
  set {char}($name + 51) = 122
  set {char}($name + 52) = 45
  set {char}($name + 53) = 111
  set {char}($name + 54) = 117
  set {char}($name + 55) = 116
  set {char}($name + 56) = 112
  set {char}($name + 57) = 117
  set {char}($name + 58) = 116
  set {char}($name + 59) = 46
  set {char}($name + 60) = 116
  set {char}($name + 61) = 120
  set {char}($name + 62) = 116
  set {char}($name + 63) = 0
  x/s $name
  printf "LIPSYNC_FILENAME_CALL case=absolute-z\n"
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $name, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_FILENAME case=absolute-z name=Z:/work/stage16/lipsync-filename-bytewise-absolute-z-output.txt raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
