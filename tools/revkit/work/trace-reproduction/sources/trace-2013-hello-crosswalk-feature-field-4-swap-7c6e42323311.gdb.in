set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-plain-crosswalk-feature-field-4/adapted-feature-swap-gdb.log
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
      printf "FEATURE_FIELD_4_SWAP_CUTOFF_OVERRIDE slot=%u sum=%d\n", $slot, $last_sum
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

# Swap only scalar feature array slot 4 between each old/new competing row pair
# while FUN_100182e0 scores those rows. Keys, signatures, attr_b, and transition
# codes retain their original values.
break *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if ($unit == 272823 || $unit == 273370 || $unit == 272824 || $unit == 273371 || $unit == 272825 || $unit == 255282)
    set $model = *(unsigned int *)($esp + 24)
    set $p0 = *(unsigned int *)($model + 0x48)
    set $p1 = *(unsigned int *)($model + 0x4c)
    set $p2 = *(unsigned int *)($model + 0x50)
    set $p3 = *(unsigned int *)($model + 0x54)
    set $p4 = *(unsigned int *)($model + 0x58)
    set $p5 = *(unsigned int *)($model + 0x5c)
    if $unit == 272823 || $unit == 273370
      set $a = 272823
      set $b = 273370
    else
      if $unit == 272824 || $unit == 273371
        set $a = 272824
        set $b = 273371
      else
        set $a = 272825
        set $b = 255282
      end
    end
    set $v40 = *(unsigned char *)($p4 + $a)
    set $v41 = *(unsigned char *)($p4 + $b)
    set {unsigned char}($p4 + $a) = $v41
    set {unsigned char}($p4 + $b) = $v40
    set $return = *(unsigned int *)$esp
    printf "FEATURE_FIELD_4_SWAP_ENTRY unit=%u pair=%u,%u\n", $unit, $a, $b
    tbreak *$return
    commands
      silent
      printf "FEATURE_FIELD_4_SWAP_RETURN unit=%u pair=%u,%u score=%g\n", $unit, $a, $b, $st0
      set {unsigned char}($p4 + $a) = $v40
      set {unsigned char}($p4 + $b) = $v41
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
    printf "FEATURE_FIELD_4_SWAP_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "FEATURE_FIELD_4_SWAP_TRACE_READY\n"
  continue
end

continue
