set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-2013-producer/adapted-producer-gdb.log
set logging overwrite on
set logging enabled on

set $producer_calls = 0
hbreak *0x10018770
commands
  silent
  set $producer_calls = $producer_calls + 1
  set $signature = *(unsigned int *)($esp + 4)
  set $slot = *(unsigned int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $mode = *(unsigned short *)($esp + 20)
  set $row = $state + ($slot * 3 + 0x76314) * 2
  printf "2013_HELLO_QUERY_PRODUCER n=%u slot=%u mode=%u row:", $producer_calls, $slot, $mode
  x/6ub $row
  printf "2013_HELLO_QUERY_SOURCE_SIGNATURE:"
  x/7ub $signature
  continue
end

hbreak *0x408187
commands
  silent
  printf "2013_HELLO_QUERY_PRODUCER_TRACE_READY\n"
  continue
end

continue
