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
  set {char}($name + 13) = 92
  set {char}($name + 14) = 98
  set {char}($name + 15) = 97
  set {char}($name + 16) = 99
  set {char}($name + 17) = 107
  set {char}($name + 18) = 115
  set {char}($name + 19) = 108
  set {char}($name + 20) = 97
  set {char}($name + 21) = 115
  set {char}($name + 22) = 104
  set {char}($name + 23) = 46
  set {char}($name + 24) = 116
  set {char}($name + 25) = 120
  set {char}($name + 26) = 116
  set {char}($name + 27) = 0
  x/s $name
  printf "LIPSYNC_FILENAME_CALL case=backslash-relative\n"
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $name, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_FILENAME case=backslash-relative name=lipsync-probe\backslash.txt raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
