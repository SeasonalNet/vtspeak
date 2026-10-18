set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $zero_raw = 1
set $zero_a = 1
set $zero_b = 0
set $transition_calls = 0
set $selected_calls = 0

break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  if $context == 2
    set $transition_calls = $transition_calls + 1
    set $raw_address = 0x1007c2a0
    set $b_address = 0x1007c2a4
    set $a_address = 0x1007c2a8
    set $raw_weight = *(float *)$raw_address
    set $b_weight = *(float *)$b_address
    set $a_weight = *(float *)$a_address
    set $return = *(unsigned int *)$esp
    printf "APPLE_WEIGHT_SUBSET_BEGIN context=2 raw_a_b=%g,%g,%g\n", $raw_weight, $a_weight, $b_weight
    if $zero_raw
      set {float}$raw_address = 0.0
    end
    if $zero_a
      set {float}$a_address = 0.0
    end
    if $zero_b
      set {float}$b_address = 0.0
    end
    printf "APPLE_WEIGHT_SUBSET_ACTIVE context=2 raw_a_b=%g,%g,%g\n", *(float *)$raw_address, *(float *)$a_address, *(float *)$b_address
    tbreak *$return
    commands
      silent
      set {float}$raw_address = $raw_weight
      set {float}$b_address = $b_weight
      set {float}$a_address = $a_weight
      printf "APPLE_WEIGHT_SUBSET_END context=2 restored_raw_a_b=%g,%g,%g\n", *(float *)$raw_address, *(float *)$a_address, *(float *)$b_address
      continue
    end
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "APPLE_SELECTED_UNIT_HARDWARE n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_WEIGHT_SUBSET_TRACE_READY\n"
  continue
end

continue
