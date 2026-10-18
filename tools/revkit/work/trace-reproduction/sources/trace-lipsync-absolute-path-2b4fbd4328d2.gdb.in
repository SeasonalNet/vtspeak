set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $path = (char *)malloc(64)
  set {char[64]}$path = "Z:/work/stage16/lipsync-absolute-path-report.txt"
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $path, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_PATH absolute=Z:/work/stage16/lipsync-absolute-path-report.txt raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
