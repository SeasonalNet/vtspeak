set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $empty = (char *)malloc(1)
  set {char}($empty) = 0
  set $result = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $empty, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_VALID speaker=1 selector_string=empty options=-1 raw_eax=%#x signed=%d\n", $result, $result
  set $fallback = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)((char *)$text, $empty, -1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_VALID speaker=-1 selector_string=empty options=-1 raw_eax=%#x signed=%d\n", $fallback, $fallback
  set $null_text = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)(0, $empty, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_NULL_TEXT raw_eax=%#x signed=%d\n", $null_text, $null_text
  set $empty_text = ((unsigned int (*)(char *, char *, int, int, int, int, int, int, int))0x1001de50)($empty, $empty, 1, -1, -1, -1, -1, -1, -1)
  printf "LIPSYNC_EMPTY_TEXT raw_eax=%#x signed=%d\n", $empty_text, $empty_text
  continue
end

continue
