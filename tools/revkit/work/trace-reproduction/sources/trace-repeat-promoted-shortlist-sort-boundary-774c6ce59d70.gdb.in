set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-promoted-shortlist-sort-boundary/adapted-promoted-shortlist-sort-gdb.log
set logging overwrite on
set logging enabled on

# Preserve the 2006 first choice, then promote 264071 before the initial
# candidate sort at this position.
break *0x10023331
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  if $position == 0 || $position == 6
    set $count = *(short *)($state + 0xae988 + $position * 0xfc)
    set $index = 0
    while $index < $count
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if $position == 0 && *(unsigned int *)($node + 8) == 273369
        set {float}($node + 4) = -2.0
      end
      if $position == 6 && *(unsigned int *)($node + 8) == 264071
        printf "ORDER_OVERRIDE position=6 unit=264071 previous_key=%g", *(float *)($node + 4)
        set {float}($node + 4) = -2.0
        printf " forced_key=%g pre_sort_index=%d\n", *(float *)($node + 4), $index
        tbreak *0x10023336
        commands
          silent
          set $count = *(short *)($state + 0xae988 + 6 * 0xfc)
          printf "INITIAL_SORTED_POSITION6 count=%d", $count
          set $i = 0
          while $i < $count && $i < 30
            set $sorted_node = *(unsigned int *)($state + 0x477ac + $i * 4)
            printf " [%d]=%u(key=%g)", $i, *(unsigned int *)($sorted_node + 8), *(float *)($sorted_node + 4)
            set $i = $i + 1
          end
          printf "\n"
          continue
        end
      end
      set $index = $index + 1
    end
  end
  continue
end

# Record final local scorer values for the target old/new units.
break *0x1002389a
commands
  silent
  set $position = *(int *)($ebp + 8)
  if $position == 6
    set $node = *(unsigned int *)$esi
    set $unit = *(unsigned int *)($node + 8)
    if $unit == 264071 || $unit == 49874
      printf "SHORTLIST_UNIT_SCORE position=6 unit=%u score=%g pre_sort_key=%g span=%d weighted=%d\n", $unit, $st0, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
    end
  end
  continue
end

# Capture the actual sort input and its retained 30-row prefix.
break *0x10023948
commands
  silent
  set $position = *(int *)($ebp + 8)
  if $position == 6
    set $state = *(unsigned int *)($ebp + 0xc)
    set $total = $ebx
    set $prefix = *(int *)($ebp - 0x20)
    set $sort_count = *(int *)$esp
    set $sort_ptr = *(unsigned int *)($esp + 4)
    printf "SHORTLIST_SORT_INPUT position=6 total=%d presorted_prefix=%d suffix_count=%d suffix_ptr=%#x\n", $total, $prefix, $sort_count, $sort_ptr
    set $i = 0
    while $i < $total && $i < 100
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      set $unit = *(unsigned int *)($node + 8)
      if $unit == 264071 || $unit == 49874
        printf "SHORTLIST_SORT_TARGET position=6 index=%d unit=%u score=%g key=%g span=%d weighted=%d flags=%#x\n", $i, $unit, *(float *)$node, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), *(unsigned char *)($node + 0x12)
      end
      set $i = $i + 1
    end
    tbreak *0x1002394d
    commands
      silent
      set $i = 0
      printf "SHORTLIST_SORTED_TOP30 position=6"
      while $i < 30
        set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
        printf " [%d]=%u(score=%g,key=%g,span=%d,weighted=%d,flags=%#x)", $i, *(unsigned int *)($node + 8), *(float *)$node, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), *(unsigned char *)($node + 0x12)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x408187
commands
  silent
  printf "PROMOTED_SHORTLIST_SORT_TRACE_READY\n"
  continue
end

continue
