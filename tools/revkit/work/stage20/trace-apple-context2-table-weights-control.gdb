set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $zero_table_weights = 0
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
    printf "APPLE_TABLE_WEIGHTS_BEGIN context=2 call=%u zero=%d raw_a_b=%g,%g,%g\n", $transition_calls, $zero_table_weights, $raw_weight, $a_weight, $b_weight
    if $zero_table_weights
      set {float}$raw_address = 0.0
      set {float}$b_address = 0.0
      set {float}$a_address = 0.0
    end
    printf "APPLE_TABLE_WEIGHTS_ACTIVE context=2 call=%u raw_a_b=%g,%g,%g\n", $transition_calls, *(float *)$raw_address, *(float *)$a_address, *(float *)$b_address
    tbreak *$return
    commands
      silent
      if $zero_table_weights
        set {float}$raw_address = $raw_weight
        set {float}$b_address = $b_weight
        set {float}$a_address = $a_weight
      end
      printf "APPLE_TABLE_WEIGHTS_END context=2 call=%u restored_raw_a_b=%g,%g,%g\n", $transition_calls, *(float *)$raw_address, *(float *)$a_address, *(float *)$b_address
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
  printf "APPLE_TABLE_WEIGHTS_TRACE_READY\n"
  continue
end

continue
