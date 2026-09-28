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
  set {char}($name + 0) = 108
  set {char}($name + 1) = 105
  set {char}($name + 2) = 112
  set {char}($name + 3) = 115
  set {char}($name + 4) = 121
  set {char}($name + 5) = 110
  set {char}($name + 6) = 99
  set {char}($name + 7) = 45
  set {char}($name + 8) = 112
  set {char}($name + 9) = 114
  set {char}($name + 10) = 111
  set {char}($name + 11) = 98
  set {char}($name + 12) = 101
  set {char}($name + 13) = 47
  set {char}($name + 14) = 118
  set {char}($name + 15) = 116
  set {char}($name + 16) = 115
  set {char}($name + 17) = 112
  set {char}($name + 18) = 101
  set {char}($name + 19) = 97
  set {char}($name + 20) = 107
  set {char}($name + 21) = 45
  set {char}($name + 22) = 108
  set {char}($name + 23) = 105
  set {char}($name + 24) = 112
  set {char}($name + 25) = 114
  set {char}($name + 26) = 101
  set {char}($name + 27) = 108
  set {char}($name + 28) = 97
  set {char}($name + 29) = 116
  set {char}($name + 30) = 105
  set {char}($name + 31) = 118
  set {char}($name + 32) = 101
  set {char}($name + 33) = 46
  set {char}($name + 34) = 116
  set {char}($name + 35) = 120
  set {char}($name + 36) = 116
  set {char}($name + 37) = 0
  x/s $name
  printf "LIPSYNC_FILENAME_CALL case=subdir-relative-verified\n"
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $name, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_FILENAME case=subdir-relative-verified name=lipsync-probe/vtspeak-liprelative.txt raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
