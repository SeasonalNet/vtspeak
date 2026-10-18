set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $path = *(unsigned int *)($esp + 12)
  disable 1
  printf "HEAP_EXPORT before_unload=0x%08x\n", *(unsigned int *)0x100ff11c
  call ((void (*)(int))0x10027ea0)(1)
  printf "DB_UNLOAD called speaker=1\n"
  printf "HEAP_EXPORT after_unload=0x%08x\n", *(unsigned int *)0x100ff11c
  set $file = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, (char *)$text, (char *)$path, -1, -1, -1, -1, -1, -1, -1)
  printf "DB_UNLOADED file_return=%d\n", $file
  set $buffer = (char *)malloc(1048576)
  set $length = (int *)malloc(4)
  set *$length = 1048576
  set $buffer_result = ((int (*)(int, char *, char *, int *, int, int, int, int, int, int, int, int, int))0x1001dc40)(0, (char *)$text, $buffer, $length, 0, 0, 1, -1, -1, -1, -1, -1, -1)
  printf "DB_UNLOADED buffer_return=%d output_len=%d\n", $buffer_result, *$length
  set $play = ((short (*)(void *, unsigned int, char *, int, int, int, int, int, int, int))0x10027050)(0, 0, (char *)$text, -1, -1, -1, -1, -1, -1, -1)
  printf "DB_UNLOADED play_return=%d\n", $play
  continue
end

continue
