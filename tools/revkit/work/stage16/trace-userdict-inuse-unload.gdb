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
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 0
  set $out = ((char *(*)(unsigned int))0x1001d9c0)(42)
  set {char}($out + 0) = 90
  set {char}($out + 1) = 58
  set {char}($out + 2) = 47
  set {char}($out + 3) = 119
  set {char}($out + 4) = 111
  set {char}($out + 5) = 114
  set {char}($out + 6) = 107
  set {char}($out + 7) = 47
  set {char}($out + 8) = 115
  set {char}($out + 9) = 116
  set {char}($out + 10) = 97
  set {char}($out + 11) = 103
  set {char}($out + 12) = 101
  set {char}($out + 13) = 49
  set {char}($out + 14) = 54
  set {char}($out + 15) = 47
  set {char}($out + 16) = 117
  set {char}($out + 17) = 115
  set {char}($out + 18) = 101
  set {char}($out + 19) = 114
  set {char}($out + 20) = 100
  set {char}($out + 21) = 105
  set {char}($out + 22) = 99
  set {char}($out + 23) = 116
  set {char}($out + 24) = 45
  set {char}($out + 25) = 105
  set {char}($out + 26) = 110
  set {char}($out + 27) = 117
  set {char}($out + 28) = 115
  set {char}($out + 29) = 101
  set {char}($out + 30) = 45
  set {char}($out + 31) = 117
  set {char}($out + 32) = 110
  set {char}($out + 33) = 108
  set {char}($out + 34) = 111
  set {char}($out + 35) = 97
  set {char}($out + 36) = 100
  set {char}($out + 37) = 46
  set {char}($out + 38) = 119
  set {char}($out + 39) = 97
  set {char}($out + 40) = 118
  set {char}($out + 41) = 0
  set $load = ((short (*)(int, char *))0x10027960)(200, $path)
  printf "USERDICT_INUSE load index=200 ax=%d\n", $load
  set $gate_before = *(unsigned char *)0x100a7489
  set {unsigned char}0x100a7489 = 1
  printf "USERDICT_INUSE gate_before=%u gate_forced=1\n", $gate_before
  break *0x100260ec
  commands 2
    silent
    disable 2
    printf "USERDICT_INUSE active_state=%#x selected_dict_ptr=%#x slot200_ptr=%#x\n", $esi, *(unsigned int *)($esi + 0x1312c0), *(unsigned int *)(0x100a647c + 200 * 4)
    set $during = ((short (*)(int))0x10027980)(200)
    printf "USERDICT_INUSE unload_during_selection index=200 ax=%d\n", $during
    continue
  end
  set $synth = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $out, 1, -1, -1, -1, -1, 200, 0)
  printf "USERDICT_INUSE synth index=200 ax=%d\n", $synth
  set $after = ((short (*)(int))0x10027980)(200)
  printf "USERDICT_INUSE unload_after_synthesis index=200 ax=%d\n", $after
  set {unsigned char}0x100a7489 = $gate_before
  printf "USERDICT_INUSE gate_restored=%u\n", *(unsigned char *)0x100a7489
  continue
end

continue
