set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-span-boundary/adapted-span-boundary-gdb.log
set logging overwrite on
set logging enabled on

# FUN_100230a0 annotates the candidate rows assembled by FUN_10023350. Capture
# the position's input IDs and the two target unit signatures before the call,
# then capture their calculated coverage fields at the return address.
break *0x100230a0
commands
  silent
  set $position = *(int *)($esp + 4)
  if $position == 6
    set $state = *(unsigned int *)($esp + 8)
    set $model = *(unsigned int *)($esp + 12)
    set $return = *(unsigned int *)$esp
    set $record = $state + 0xae894 + $position * 0xfc
    set $count = *(unsigned short *)($record + 0xf4)
    set $unit_signatures = *(unsigned int *)($model + 0x64)
    set $position_row = $state + (0x76314 + $position * 3) * 2
    printf "SPAN2013_ENTRY position=%u record_count=%u packed_row=", $position, $count
    x/6ub $position_row
    printf "SPAN2013_POSITION_IDS="
    set $i = 0
    set $position_count = *(unsigned short *)$record
    while $i < $position_count && $i < 12
      printf "%u ", *(unsigned int *)($record + 4 + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    set $i = 0
    while $i < $count && $i < 100
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      set $unit = *(unsigned int *)($node + 8)
      if $unit == 264071 || $unit == 49874
        printf "SPAN2013_UNIT_SIGNATURE position=%u index=%u unit=%u bytes=", $position, $i, $unit
        x/7ub ($unit_signatures + $unit * 7)
      end
      set $i = $i + 1
    end
    tbreak *$return
    commands
      silent
      set $i = 0
      while $i < $count && $i < 100
        set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
        set $unit = *(unsigned int *)($node + 8)
        if $unit == 264071 || $unit == 49874
          printf "SPAN2013_RESULT position=%u index=%u unit=%u right=%d total=%d weighted=%d\n", $position, $i, $unit, *(short *)($node + 0xc), *(short *)($node + 0xe), *(short *)($node + 0x10)
        end
        set $i = $i + 1
      end
      continue
    end
  end
  continue
end

break *0x408187
commands
  silent
  printf "SPAN2013_TRACE_READY\n"
  continue
end

continue
