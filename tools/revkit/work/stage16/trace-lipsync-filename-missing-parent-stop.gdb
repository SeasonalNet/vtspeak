set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV stop print nopass

break *0x10025e4d
commands
  silent
  printf "LIPSYNC_MISSING_PARENT_WRITER eip=%#x edx=%#x ecx=%#x eax=%#x esp=%#x\n", $eip, $edx, $ecx, $eax, $esp
  x/i $eip
  if $edx != 0
    x/wx $edx
  end
  bt
  disable 1
  handle SIGSEGV nostop noprint pass
  continue
end

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 2
  set $name = (char *)malloc(128)
  set {char}($name + 0) = 110
  set {char}($name + 1) = 111
  set {char}($name + 2) = 45
  set {char}($name + 3) = 115
  set {char}($name + 4) = 117
  set {char}($name + 5) = 99
  set {char}($name + 6) = 104
  set {char}($name + 7) = 45
  set {char}($name + 8) = 108
  set {char}($name + 9) = 105
  set {char}($name + 10) = 112
  set {char}($name + 11) = 115
  set {char}($name + 12) = 121
  set {char}($name + 13) = 110
  set {char}($name + 14) = 99
  set {char}($name + 15) = 45
  set {char}($name + 16) = 100
  set {char}($name + 17) = 105
  set {char}($name + 18) = 114
  set {char}($name + 19) = 47
  set {char}($name + 20) = 114
  set {char}($name + 21) = 101
  set {char}($name + 22) = 112
  set {char}($name + 23) = 111
  set {char}($name + 24) = 114
  set {char}($name + 25) = 116
  set {char}($name + 26) = 46
  set {char}($name + 27) = 116
  set {char}($name + 28) = 120
  set {char}($name + 29) = 116
  set {char}($name + 30) = 0
  x/s $name
  printf "LIPSYNC_FILENAME_CALL case=missing-parent\n"
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $name, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_FILENAME case=missing-parent name=no-such-lipsync-dir/report.txt raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
