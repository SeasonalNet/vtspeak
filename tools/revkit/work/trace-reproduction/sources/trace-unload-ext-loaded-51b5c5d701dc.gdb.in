set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $text = *(unsigned int *)($esp + 8)
  set $path = *(unsigned int *)($esp + 12)
  set $dbsize = (int *)malloc(4)
  set *$dbsize = 0
  set $before_result = ((int (*)(int *, int))0x10028130)($dbsize, 1)
  printf "UNLOAD_EXT_BEFORE dbsize_result=%d dbsize=%d speaker=%s\n", $before_result, *$dbsize, ((char *(*)(int))0x10028230)(1)
  call ((void (*)(int))0x10027ea0)(1)
  set *$dbsize = 0
  set $after_result = ((int (*)(int *, int))0x10028130)($dbsize, 1)
  printf "UNLOAD_EXT_AFTER dbsize_result=%d dbsize=%d speaker=%s heap_start=%#x\n", $after_result, *$dbsize, ((char *(*)(int))0x10028230)(1), *(unsigned int *)0x100ff11c
  set $file = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, (char *)$text, (char *)$path, -1, -1, -1, -1, -1, -1, -1)
  printf "UNLOAD_EXT_FILE result=%d\n", $file
  set $buffer = (char *)malloc(1048576)
  set $length = (int *)malloc(4)
  set *$length = 1048576
  set $buffer_result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, $buffer, $length, 0, 0, 1, -1, -1, -1, -1, -1, -1)
  printf "UNLOAD_EXT_BUFFER result=%d length=%d\n", $buffer_result, *$length
  continue
end

continue
