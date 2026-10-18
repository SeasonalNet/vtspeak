set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot0-pool-to-a000-v2/adapted-path-costs-gdb.log
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
    set $count = *(short *)($state + 0xae988)
    set $selected_index = *(short *)($state + 0xae0c4)
    printf "SLOT0_ALIAS_PATH_COSTS first=%d last=%d count=%d selected_index=%d\n", $first, $last, $count, $selected_index
    set $index = 0
    while $index < $count && $index < 100
      set $record = $state + 0xae894
      set $unit = *(unsigned int *)($record + 0x7c + $index * 4)
      set $cost = *(float *)($state + 0x8212c + $index * 4)
      printf "SLOT0_ALIAS_PATH index=%d unit=%u cumulative=%g selected=%d\n", $index, $unit, $cost, $index == $selected_index
      set $index = $index + 1
    end
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "SLOT0_ALIAS_PATH_TRACE_READY\n"
  continue
end

continue
