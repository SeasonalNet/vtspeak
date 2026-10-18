set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $ablate_raw = 0
set $transition_calls = 0
set $selected_calls = 0

break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  if $context == 2
    set $transition_calls = $transition_calls + 1
    set $weight_address = 0x1007c2a0
    set $raw_weight = *(float *)$weight_address
    set $return = *(unsigned int *)$esp
    printf "APPLE_RAW_WEIGHT_BEGIN context=2 call=%u ablate=%d raw_weight=%g\n", $transition_calls, $ablate_raw, $raw_weight
    if $ablate_raw
      set {float}$weight_address = 0.0
    end
    printf "APPLE_RAW_WEIGHT_ACTIVE context=2 call=%u raw_weight=%g\n", $transition_calls, *(float *)$weight_address
    tbreak *$return
    commands
      silent
      if $ablate_raw
        set {float}$weight_address = $raw_weight
      end
      printf "APPLE_RAW_WEIGHT_END context=2 call=%u restored=%g\n", $transition_calls, *(float *)$weight_address
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
  printf "APPLE_RAW_WEIGHT_TRACE_READY\n"
  continue
end

continue
