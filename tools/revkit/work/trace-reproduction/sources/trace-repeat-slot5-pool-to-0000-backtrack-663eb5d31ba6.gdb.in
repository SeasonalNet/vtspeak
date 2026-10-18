set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot5-pool-to-0000-backtrack/adapted-backtrack-gdb.log
set logging overwrite on
set logging enabled on

break *0x10024900
commands
  silent
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    set $context = $first
    while $context <= $last
      set $record = $state + 0xae894 + $context * 0xfc
      set $count = *(short *)($state + 0xae988 + $context * 0xfc)
      set $index = *(short *)($state + 0xae0c4 + $context * 2)
      set $unit = *(unsigned int *)($record + 0x7c + $index * 4)
      printf "SLOT5_ALIAS_ROW context=%d count=%d selected_index=%d selected_unit=%u\n", $context, $count, $index, $unit
      set $context = $context + 1
    end
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "SLOT5_ALIAS_BACKTRACK_TRACE_READY\n"
  continue
end

continue
