set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-slot0-old-winner-order-override/adapted-order-and-chain-gdb.log
set logging overwrite on
set logging enabled on

# Give the old engine's first selected row priority among exact rank ties.
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
        printf "OLD_WINNER_ORDER before=%g id=273369\n", *(float *)($node + 4)
        set {float}($node + 4) = -2.0
        printf "OLD_WINNER_ORDER after=%g id=273369\n", *(float *)($node + 4)
      end
      set $index = $index + 1
    end
  end
  continue
end

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
      printf "ORDER_OVERRIDE_ROW context=%d count=%d selected_index=%d selected_unit=%u\n", $context, $count, $index, $unit
      set $context = $context + 1
    end
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "ORDER_OVERRIDE_TRACE_READY\n"
  continue
end

continue
