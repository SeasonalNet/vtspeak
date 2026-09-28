set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/pah0-legacy-key-substitution-metrics/gdb.log
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
  printf "LEGACY_KEY_PATCH_ENTER n=%u slot=%u bank=%u leaf=%u before:", $split_calls, $slot, $bank, $leaf
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
  printf "LEGACY_KEY_PATCH_ACTIVE n=%u slot=%u signature:", $split_calls, $slot
  x/7ub $signature
  tbreak *$return
  commands
    silent
    set {unsigned char}($signature + 1) = $old1
    set {unsigned char}($signature + 2) = $old2
    set {unsigned char}($signature + 3) = $old3
    set {unsigned char}($signature + 5) = $old5
    printf "LEGACY_KEY_PATCH_RESTORED n=%u slot=%u accepted=%u signature:", $split_calls, $slot, $al
    x/7ub $signature
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
    printf "LEGACY_KEY_METRIC_ENTRY caller=%#x count=%u ids_weights:", $caller, $count
    set $i = 0
    while $i < $count && $i < 64
      set $id = *(unsigned int *)($ids + $i * 4)
      printf " %u:%u", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      printf "LEGACY_KEY_METRIC_RETURN sum=%u\n", $eax
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
  printf "PAH0_LEGACY_KEY_SUBSTITUTION_METRICS_READY\n"
  continue
end

continue
