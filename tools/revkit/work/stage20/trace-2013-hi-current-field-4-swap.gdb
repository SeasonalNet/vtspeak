set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-2013-current-field-4/adapted-feature-swap-gdb.log
set logging overwrite on
set logging enabled on

set $last_sum = -1
break *0x10024060
commands
  silent
  set $last_sum = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    if $al == 0 && $last_sum > 2
      printf "HI_CURRENT_FIELD_4_CUTOFF_OVERRIDE slot=%u sum=%d\n", $slot, $last_sum
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
    set $return = $caller
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      continue
    end
  end
  continue
end

# Exchange only scalar feature field 4 between the two shared position-0
# records while FUN_100182e0 scores them; restore both values on function return.
break *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if $unit == 272821 || $unit == 280449
    set $model = *(unsigned int *)($esp + 24)
    set $values = *(unsigned int *)($model + 0x58)
    set $old_a = *(unsigned char *)($values + 272821)
    set $old_b = *(unsigned char *)($values + 280449)
    set {unsigned char}($values + 272821) = $old_b
    set {unsigned char}($values + 280449) = $old_a
    set $return = *(unsigned int *)$esp
    printf "HI_CURRENT_FIELD_4_SWAP_ENTRY unit=%u pair=272821,280449 values=%u,%u\n", $unit, $old_a, $old_b
    tbreak *$return
    commands
      silent
      printf "HI_CURRENT_FIELD_4_SWAP_RETURN unit=%u score=%g\n", $unit, $st0
      set {unsigned char}($values + 272821) = $old_a
      set {unsigned char}($values + 280449) = $old_b
      continue
    end
  end
  continue
end

set $selected_calls = 0
break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 16
    printf "HI_CURRENT_FIELD_4_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_CURRENT_FIELD_4_TRACE_READY\n"
  continue
end

continue
