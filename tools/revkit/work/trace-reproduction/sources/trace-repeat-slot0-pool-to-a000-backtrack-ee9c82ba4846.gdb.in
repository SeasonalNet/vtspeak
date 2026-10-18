set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot0-pool-to-a000-v2/adapted-backtrack-gdb.log
set logging overwrite on
set logging enabled on
set $backtrack_calls = 0

break *0x10024900
commands
  silent
  set $backtrack_calls = $backtrack_calls + 1
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 12)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    printf "SLOT0_ALIAS_BACKTRACK call=%u first=%d last=%d\n", $backtrack_calls, $first, $last
    set $context = $first
    while $context <= $last
      set $record = $state + 0xae894 + $context * 0xfc
      set $count = *(short *)($state + 0xae988 + $context * 0xfc)
      set $index = *(short *)($state + 0xae0c4 + $context * 2)
      set $unit = *(unsigned int *)($record + 0x7c + $index * 4)
      printf "SLOT0_ALIAS_ROW context=%d count=%d selected_index=%d selected_unit=%u\n", $context, $count, $index, $unit
      set $context = $context + 1
    end
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "SLOT0_ALIAS_BACKTRACK_TRACE_READY\n"
  continue
end

continue
