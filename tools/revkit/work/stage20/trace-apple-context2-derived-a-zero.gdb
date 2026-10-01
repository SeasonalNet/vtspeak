set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $transition_calls = 0
set $selected_calls = 0

break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  if $context == 2
    set $transition_calls = $transition_calls + 1
    set $weight_address = 0x1007c2a8
    set $saved_weight = *(float *)$weight_address
    set $return = *(unsigned int *)$esp
    printf "APPLE_DERIVED_A_BEGIN context=2 call=%u weight=%g\n", $transition_calls, $saved_weight
    set {float}$weight_address = 0.0
    printf "APPLE_DERIVED_A_ACTIVE context=2 call=%u weight=%g\n", $transition_calls, *(float *)$weight_address
    tbreak *$return
    commands
      silent
      set {float}$weight_address = $saved_weight
      printf "APPLE_DERIVED_A_END context=2 call=%u restored=%g\n", $transition_calls, *(float *)$weight_address
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
  printf "APPLE_DERIVED_A_TRACE_READY\n"
  continue
end

continue
