set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/legacy-query-expansion-gdb.log
set logging overwrite on
set logging enabled on

set $generated_calls = 0
hbreak *0x1001da00
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x1001da80 && $caller < 0x1001dfb0
    set $generated_calls = $generated_calls + 1
    set $key = *(unsigned int *)($esp + 4)
    set $mode = *(unsigned short *)($esp + 8)
    printf "LEGACY_QUERY_VARIANT n=%u caller=%#x mode=%u key:", $generated_calls, $caller, $mode
    x/5ub $key
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_QUERY_EXPANSION_TRACE_READY\n"
  continue
end

continue
