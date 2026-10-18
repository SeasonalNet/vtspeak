set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-bit7-chain-span/adapted-bit7-chain-span-gdb.log
set logging overwrite on
set logging enabled on

# Preserve the 2006 first choice for comparison with the completed-row trace.
break *0x10023331
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  if $position == 0
    set $count = *(short *)($state + 0xae988)
    set $index = 0
    while $index < $count
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if *(unsigned int *)($node + 8) == 273369
        set {float}($node + 4) = -2.0
      end
      set $index = $index + 1
    end
  end
  continue
end

# Save the third argument at the actual FUN_10023350 entry. The preceding
# breakpoint is in its caller and does not have this callee's argument frame.
break *0x10023350
commands
  silent
  set $position = *(int *)($esp + 4)
  if $position == 6
    set $model = *(unsigned int *)($esp + 12)
    set $signatures = *(unsigned int *)($model + 0x64)
    set $unit = 264071
    while $unit <= 264073
      set $unit_signature = $signatures + $unit * 7
      printf "SPAN_BIT7_INTERVENTION model=%#x unit=%u old_byte6=%#x\n", $model, $unit, *(unsigned char *)($unit_signature + 6)
      set {unsigned char}($unit_signature + 6) = *(unsigned char *)($unit_signature + 6) | 0x80
      printf "SPAN_BIT7_INTERVENTION unit=%u new_byte6=%#x\n", $unit, *(unsigned char *)($unit_signature + 6)
      set $unit = $unit + 1
    end
  end
  continue
end

# FUN_100182e0 returns in ST(0); its caller stores that value at node +0.
# Record the legacy row and its surviving 2013 alternative after unit scoring.
break *0x1002389a
commands
  silent
  set $position = *(int *)($ebp + 8)
  if $position == 6
    set $node = *(unsigned int *)$esi
    set $unit = *(unsigned int *)($node + 8)
    if $unit == 264071 || $unit == 49874
      printf "SHORTLIST_UNIT_SCORE position=6 unit=%u score=%g pre_sort_key=%g span=%d weighted=%d\n", $unit, $st0, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
      set $signatures = *(unsigned int *)($model + 0x64)
      printf "SHORTLIST_UNIT_SIGNATURE model=%#x address=%#x position=6 unit=%u signature=", $model, $signatures + $unit * 7, $unit
      x/7ub ($signatures + $unit * 7)
    end
  end
  continue
end

# Capture the actual sort input: scored candidates, the already ordered prefix,
# the suffix passed to FUN_1001b5f0, and the resulting 30-row DP shortlist.
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
  printf "SHORTLIST_SORT_TRACE_READY\n"
  continue
end

continue
