set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-volume100-rows-gdb.log
set logging overwrite on
set logging enabled on
set $row_calls = 0

hbreak *0x10027fe0
commands
  silent
  printf "HI_VOLUME100_SETTER original=%d\n", *(int *)($esp + 12)
  set *(int *)($esp + 12) = 100
  continue
end

hbreak *0x1002c8b0
commands
  silent
  set $row = *(unsigned int *)($esp + 8)
  set $row_calls = $row_calls + 1
  printf "HI_VOLUME100_ROW call=%u samples=%d scale=%d\n", $row_calls, *(int *)($row + 0xc), *(int *)($row + 8)
  continue
end

continue
