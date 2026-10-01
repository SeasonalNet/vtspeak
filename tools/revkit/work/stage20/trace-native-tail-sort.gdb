set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $sort_captures = 0

# FUN_10023350 passes its unscored tail to FUN_1001b5f0 at 0x10023948.
# At this callsite the cdecl stack contains the count followed by the pointer
# array. Capture only the first bounded tail large enough to use partitioning.
break *0x10023948
commands
  silent
  set $count = *(int *)$esp
  if $count > 17 && $count <= 200 && $sort_captures == 0
    set $sort_captures = $sort_captures + 1
    set $nodes = *(unsigned int *)($esp + 4)
    printf "NATIVE_TAIL_SORT_PRE count=%d", $count
    set $index = 0
    while $index < $count
      set $node = *(unsigned int *)($nodes + $index * 4)
      printf " [%d]=%u:%08x", $index, *(unsigned int *)($node + 8), *(unsigned int *)($node + 4)
      set $index = $index + 1
    end
    printf "\n"
    tbreak *0x1002394d
    commands
      silent
      printf "NATIVE_TAIL_SORT_POST count=%d", $count
      set $index = 0
      while $index < $count
        set $node = *(unsigned int *)($nodes + $index * 4)
        printf " [%d]=%u:%08x", $index, *(unsigned int *)($node + 8), *(unsigned int *)($node + 4)
        set $index = $index + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "NATIVE_TAIL_SORT_TRACE_READY\n"
  continue
end

continue
