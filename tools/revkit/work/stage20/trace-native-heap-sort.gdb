set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $heap_sort_captures = 0

# FUN_10023350 passes its unscored tail to FUN_1001b5f0 at 0x10023948.
# Replace only the process-local node scores with a sentinel pattern that
# forces the 32:1 imbalance path for a 75-entry list. No fixture or vendor file
# is changed by this intervention.
break *0x10023948
commands
  silent
  set $count = *(int *)$esp
  if $count == 75 && $heap_sort_captures == 0
    set $heap_sort_captures = $heap_sort_captures + 1
    set $nodes = *(unsigned int *)($esp + 4)
    printf "NATIVE_HEAP_SORT_PRE count=%d", $count
    set $index = 0
    while $index < $count
      set $node = *(unsigned int *)($nodes + $index * 4)
      if $index == 0
        set {float}($node + 4) = 0.0
      else
        if $index == 37
          set {float}($node + 4) = 1.0
        else
          if $index == 74
            set {float}($node + 4) = 2.0
          else
            set {float}($node + 4) = 3.0
          end
        end
      end
      printf " [%d]=%u:%08x", $index, *(unsigned int *)($node + 8), *(unsigned int *)($node + 4)
      set $index = $index + 1
    end
    printf "\n"
    tbreak *0x1002394d
    commands
      silent
      printf "NATIVE_HEAP_SORT_POST count=%d", $count
      set $index = 0
      while $index < $count
        set $node = *(unsigned int *)($nodes + $index * 4)
        printf " [%d]=%u", $index, *(unsigned int *)($node + 8)
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
  printf "NATIVE_HEAP_SORT_TRACE_READY\n"
  continue
end

continue
