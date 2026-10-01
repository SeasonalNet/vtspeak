set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/kate-key5-continuity-cutoff/continuity-cutoff-gdb.log
set logging overwrite on
set logging enabled on
set $last_sum = -1
set $selected_calls = 0

# Restore the legacy attr_48 marker to the current continuity-byte position
# for the four native 2006 Hello-selected rows. This is a runtime-only edit.
break *0x10023350
commands
  silent
  set $position = *(int *)($esp + 4)
  if $position == 0
    set $model = *(unsigned int *)($esp + 12)
    set $signatures = *(unsigned int *)($model + 0x64)
    set $unit = 272822
    while $unit <= 272824
      set $byte = $signatures + $unit * 7 + 6
      printf "KATE_CONTINUITY_RESTORE unit=%u before=%#x", $unit, *(unsigned char *)$byte
      set {unsigned char}$byte = *(unsigned char *)$byte | 0x80
      printf " after=%#x\n", *(unsigned char *)$byte
      set $unit = $unit + 1
    end
  end
  continue
end

# Measure the 2013 candidate metric sum and emulate the legacy nonempty >2
# acceptance rule for controlled whole-position comparisons.
break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  set $last_sum = -1
  tbreak *$return
  commands
    silent
    set $native = $al
    if $native == 0 && $last_sum > 2
      printf "KATE_CUTOFF_OVERRIDE slot=%u sum=%u native=%u overridden=1\n", $slot, $last_sum, $native
      set $eax = 1
    end
    continue
  end
  continue
end

break *0x10023060
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $ids = *(unsigned int *)($esp + 4)
    set $count = *(unsigned short *)($esp + 8)
    set $context = *(unsigned int *)($esp + 12)
    set $weights = *(unsigned int *)($context + 0x8c)
    set $return = $caller
    set $i = 0
    set $last_sum = 0
    while $i < $count && $i < 64
      set $id = *(unsigned int *)($ids + $i * 4)
      set $last_sum = $last_sum + *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      continue
    end
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "KATE_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "KATE_KEY5_CONTINUITY_CUTOFF_TRACE_READY\n"
  continue
end

continue
