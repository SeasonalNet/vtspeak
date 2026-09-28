set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $path = ((char *(*)(unsigned int))0x1001d9c0)(48)
  set {char}($path + 0) = 90
  set {char}($path + 1) = 58
  set {char}($path + 2) = 47
  set {char}($path + 3) = 119
  set {char}($path + 4) = 111
  set {char}($path + 5) = 114
  set {char}($path + 6) = 107
  set {char}($path + 7) = 47
  set {char}($path + 8) = 115
  set {char}($path + 9) = 116
  set {char}($path + 10) = 97
  set {char}($path + 11) = 103
  set {char}($path + 12) = 101
  set {char}($path + 13) = 49
  set {char}($path + 14) = 54
  set {char}($path + 15) = 47
  set {char}($path + 16) = 117
  set {char}($path + 17) = 115
  set {char}($path + 18) = 101
  set {char}($path + 19) = 114
  set {char}($path + 20) = 100
  set {char}($path + 21) = 105
  set {char}($path + 22) = 99
  set {char}($path + 23) = 116
  set {char}($path + 24) = 45
  set {char}($path + 25) = 118
  set {char}($path + 26) = 97
  set {char}($path + 27) = 108
  set {char}($path + 28) = 105
  set {char}($path + 29) = 100
  set {char}($path + 30) = 97
  set {char}($path + 31) = 116
  set {char}($path + 32) = 105
  set {char}($path + 33) = 111
  set {char}($path + 34) = 110
  set {char}($path + 35) = 45
  set {char}($path + 36) = 112
  set {char}($path + 37) = 108
  set {char}($path + 38) = 97
  set {char}($path + 39) = 105
  set {char}($path + 40) = 110
  set {char}($path + 41) = 45
  set {char}($path + 42) = 112
  set {char}($path + 43) = 46
  set {char}($path + 44) = 99
  set {char}($path + 45) = 115
  set {char}($path + 46) = 118
  set {char}($path + 47) = 0
  set $load = ((short (*)(int, char *))0x10027960)(200, $path)
  printf "USERDICT_INUSE load index=200 ax=%d\n", $load
  set $dict_pointer = *(unsigned int *)(0x100a647c + 200 * 4)
  set $speaker_state = *(unsigned int *)0x100a0468
  set $capacity = *(int *)($speaker_state + 0x4d14)
  set $reference_slot = (unsigned int *)0x100a147c
  set $reference_before = *$reference_slot
  set $active_context = ((char *(*)(unsigned int))0x1001d9c0)(0x1312e0)
  set {unsigned int}($active_context + 0x1312c0) = $dict_pointer
  printf "USERDICT_INUSE state=0x%x capacity=%d dictionary=0x%x reference_slot=0x%x old_reference=0x%x active_context=0x%x context_dict=0x%x\n", $speaker_state, $capacity, $dict_pointer, $reference_slot, $reference_before, $active_context, *(unsigned int *)($active_context + 0x1312c0)
  set *$reference_slot = $active_context
  set $busy = ((short (*)(int))0x10027980)(200)
  printf "USERDICT_INUSE forced_reference_unload index=200 ax=%d\n", $busy
  set *$reference_slot = $reference_before
  printf "USERDICT_INUSE reference_restored=0x%x\n", *$reference_slot
  call ((void (*)(void *))0x1001da30)($active_context)
  printf "USERDICT_INUSE scratch_context_freed=1\n"
  set $after = ((short (*)(int))0x10027980)(200)
  printf "USERDICT_INUSE idle_unload_after_restore index=200 ax=%d\n", $after
  continue
end

continue
