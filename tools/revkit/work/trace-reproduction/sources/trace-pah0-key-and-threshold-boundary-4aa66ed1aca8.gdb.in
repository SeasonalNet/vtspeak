set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/pah0-key-and-threshold-boundary/gdb.log
set logging overwrite on
set logging enabled on
set $split_calls = 0
set $selected_calls = 0

break *0x10024060
commands
  silent
  set $split_calls = $split_calls + 1
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $row = $state + ($slot * 3 + 0x76314) * 2
  set $bank = *(unsigned char *)$row
  set $leaf = *(unsigned char *)($row + 2)
  set $voice = *(unsigned int *)($state + 0x4c)
  set $signature = $voice + $bank * 0x3c0 + 0x6e2 + $leaf * 7
  set $return = *(unsigned int *)$esp
  set $old1 = *(unsigned char *)($signature + 1)
  set $old2 = *(unsigned char *)($signature + 2)
  set $old3 = *(unsigned char *)($signature + 3)
  set $old5 = *(unsigned char *)($signature + 5)
  printf "BOUNDARY_KEY_PATCH_ENTER n=%u slot=%u bank=%u leaf=%u before:", $split_calls, $slot, $bank, $leaf
  x/7ub $signature
  if $slot == 0
    set {unsigned char}($signature + 2) = 0x22
    set {unsigned char}($signature + 3) = 0x17
  else
    if $slot == 1
      set {unsigned char}($signature + 1) = 0x22
      set {unsigned char}($signature + 2) = 0x17
      set {unsigned char}($signature + 3) = 0x2b
      set {unsigned char}($signature + 5) = 0x20
    end
  end
  tbreak *$return
  commands
    silent
    set {unsigned char}($signature + 1) = $old1
    set {unsigned char}($signature + 2) = $old2
    set {unsigned char}($signature + 3) = $old3
    set {unsigned char}($signature + 5) = $old5
    printf "BOUNDARY_KEY_PATCH_RESTORED n=%u slot=%u accepted=%u\n", $split_calls, $slot, $al
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
    set $first_id = 0xffffffff
    if $count > 0
      set $first_id = *(unsigned int *)$ids
    end
    set $return = $caller
    printf "BOUNDARY_METRIC_ENTRY caller=%#x count=%u first_id=%u\n", $caller, $count, $first_id
    tbreak *$return
    commands
      silent
      if $count == 1 && $first_id == 13478 && $eax == 9
        printf "BOUNDARY_METRIC_OVERRIDE id=%u measured=9 forced=10\n", $first_id
        set $eax = 10
      else
        printf "BOUNDARY_METRIC_RESULT count=%u sum=%u\n", $count, $eax
      end
      continue
    end
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

break *0x408187
commands
  silent
  printf "PAH0_KEY_AND_THRESHOLD_BOUNDARY_READY\n"
  continue
end

continue
