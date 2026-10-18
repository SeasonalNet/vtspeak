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
  @INITIALIZE_NAME@
  x/s $name
  printf "LIPSYNC_FILENAME_CALL case=absolute-z\n"
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $name, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_FILENAME case=absolute-z name=Z:/work/stage16/lipsync-filename-bytewise-absolute-z-output.txt raw_eax=%#x signed=%d\n", $result, $result
  continue
end

continue
