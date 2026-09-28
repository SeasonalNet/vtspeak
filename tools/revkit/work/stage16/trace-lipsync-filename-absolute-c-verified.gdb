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
  set {char}($name + 0) = 67
  set {char}($name + 1) = 58
  set {char}($name + 2) = 92
  set {char}($name + 3) = 119
  set {char}($name + 4) = 105
  set {char}($name + 5) = 110
  set {char}($name + 6) = 100
  set {char}($name + 7) = 111
  set {char}($name + 8) = 119
  set {char}($name + 9) = 115
  set {char}($name + 10) = 92
  set {char}($name + 11) = 116
  set {char}($name + 12) = 101
  set {char}($name + 13) = 109
  set {char}($name + 14) = 112
  set {char}($name + 15) = 92
  set {char}($name + 16) = 118
  set {char}($name + 17) = 116
  set {char}($name + 18) = 115
  set {char}($name + 19) = 112
  set {char}($name + 20) = 101
  set {char}($name + 21) = 97
  set {char}($name + 22) = 107
  set {char}($name + 23) = 45
  set {char}($name + 24) = 108
  set {char}($name + 25) = 105
  set {char}($name + 26) = 112
  set {char}($name + 27) = 115
  set {char}($name + 28) = 121
  set {char}($name + 29) = 110
  set {char}($name + 30) = 99
  set {char}($name + 31) = 45
  set {char}($name + 32) = 99
  set {char}($name + 33) = 45
  set {char}($name + 34) = 100
  set {char}($name + 35) = 114
  set {char}($name + 36) = 105
  set {char}($name + 37) = 118
  set {char}($name + 38) = 101
  set {char}($name + 39) = 45
  set {char}($name + 40) = 111
  set {char}($name + 41) = 117
  set {char}($name + 42) = 116
  set {char}($name + 43) = 112
  set {char}($name + 44) = 117
  set {char}($name + 45) = 116
  set {char}($name + 46) = 45
  set {char}($name + 47) = 118
  set {char}($name + 48) = 101
  set {char}($name + 49) = 114
  set {char}($name + 50) = 105
  set {char}($name + 51) = 102
  set {char}($name + 52) = 105
  set {char}($name + 53) = 101
  set {char}($name + 54) = 100
  set {char}($name + 55) = 46
  set {char}($name + 56) = 116
  set {char}($name + 57) = 120
  set {char}($name + 58) = 116
  set {char}($name + 59) = 0
  x/s $name
  printf "LIPSYNC_FILENAME_CALL case=absolute-c-verified\n"
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $name, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_FILENAME case=absolute-c-verified raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
