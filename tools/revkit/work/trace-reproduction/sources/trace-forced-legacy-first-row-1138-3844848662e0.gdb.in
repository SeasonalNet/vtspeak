set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/forced-legacy-first-row-1138-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0
set $timeline_calls = 0

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls == 1
    set $forced_unit = 272822
  else
    if $selected_calls == 2
      set $forced_unit = 272823
    else
      if $selected_calls == 3
        set $forced_unit = 272824
      else
        set $forced_unit = 272825
      end
    end
  end
  set *(unsigned int *)($esp + 8) = $forced_unit
  printf "FORCED_UNIT n=%u forced=%u\n", $selected_calls, $forced_unit
  continue
end

break *0x1002c220
commands
  silent
  set $timeline_calls = $timeline_calls + 1
  set $rows = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $start = *(unsigned short *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "TIMELINE_ENTRY n=%u start=%u input_count=%u\n", $timeline_calls, $start, *(unsigned int *)($state + 0xec624)
  tbreak *$ret
  commands
    silent
    set $count = *(unsigned short *)($rows + 0xdf70)
    if $count == 7
      set *(unsigned short *)($rows + 0xdf70) = 4
      printf "TRIM_TIMELINE n=%u before=7 after=4\n", $timeline_calls
      set $count = 4
    end
    if $count == 4
      set $old_first_count = *(unsigned int *)($rows + 0x0c)
      set *(unsigned int *)($rows + 0x0c) = 1138
      printf "EXTEND_FIRST_ROW old=%u new=1138\n", $old_first_count
    end
    printf "TIMELINE_RETURN n=%u records=%u\n", $timeline_calls, $count
    set $i = 0
    while $i < $count && $i < 32
      printf "TIMELINE_ROW i=%u words:", $i
      x/13wx $rows + ($i * 0x34)
      set $i = $i + 1
    end
    continue
  end
  continue
end

continue
