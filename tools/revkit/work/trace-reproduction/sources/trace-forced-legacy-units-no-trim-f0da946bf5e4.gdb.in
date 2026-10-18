set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/forced-legacy-units-no-trim/adapted-forced-pah0-gdb.log
set logging overwrite on
set logging enabled on
set $tree_calls = 0
set $selected_calls = 0
set $timeline_calls = 0

break *0x10001670
commands
  silent
  set $tree_calls = $tree_calls + 1
  set $tree = *(unsigned int *)($esp + 4)
  set $features = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "TREE3_SCALAR_ENTRY n=%u tree=%#x features:", $tree_calls, $tree
  x/16hd $features
  tbreak *$ret
  commands
    silent
    printf "TREE3_SCALAR_RETURN n=%u value=%d\n", $tree_calls, (short)$eax
    continue
  end
  continue
end

break *0x100016a0
commands
  silent
  set $tree_calls = $tree_calls + 1
  set $tree = *(unsigned int *)($esp + 4)
  set $features = *(unsigned int *)($esp + 8)
  set $out = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "TREE3_VECTOR_ENTRY n=%u tree=%#x features:", $tree_calls, $tree
  x/16hd $features
  tbreak *$ret
  commands
    silent
    printf "TREE3_VECTOR_RETURN n=%u values:\n", $tree_calls
    x/12hd $out
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  set $original_unit = *(unsigned int *)($esp + 8)
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
  printf "FORCED_UNIT n=%u original=%u forced=%u\n", $selected_calls, $original_unit, $forced_unit
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
