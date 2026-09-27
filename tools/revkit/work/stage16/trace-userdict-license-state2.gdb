set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10029b80
commands
  silent
  set $license_return = *(unsigned int *)$esp
  tbreak *$license_return
  commands
    silent
    printf "LICENSE_CHECK_RESULT=%d\\n", (int)$eax
    continue
  end
  continue
end

break *0x1001da50
commands
  silent
  disable 2
  printf "LICENSE_GATE gate=%u capacity=%u\\n", *(unsigned char *)0x100a7489, *(unsigned int *)(*(unsigned int *)0x100a0468 + 0x4d14)
  set $path_p = ((char *(*)(unsigned int))0x1001d9c0)(41)
  set {char}($path_p + 0) = 90
  set {char}($path_p + 1) = 58
  set {char}($path_p + 2) = 47
  set {char}($path_p + 3) = 119
  set {char}($path_p + 4) = 111
  set {char}($path_p + 5) = 114
  set {char}($path_p + 6) = 107
  set {char}($path_p + 7) = 47
  set {char}($path_p + 8) = 115
  set {char}($path_p + 9) = 116
  set {char}($path_p + 10) = 97
  set {char}($path_p + 11) = 103
  set {char}($path_p + 12) = 101
  set {char}($path_p + 13) = 49
  set {char}($path_p + 14) = 54
  set {char}($path_p + 15) = 47
  set {char}($path_p + 16) = 117
  set {char}($path_p + 17) = 115
  set {char}($path_p + 18) = 101
  set {char}($path_p + 19) = 114
  set {char}($path_p + 20) = 100
  set {char}($path_p + 21) = 105
  set {char}($path_p + 22) = 99
  set {char}($path_p + 23) = 116
  set {char}($path_p + 24) = 45
  set {char}($path_p + 25) = 104
  set {char}($path_p + 26) = 101
  set {char}($path_p + 27) = 108
  set {char}($path_p + 28) = 108
  set {char}($path_p + 29) = 111
  set {char}($path_p + 30) = 45
  set {char}($path_p + 31) = 112
  set {char}($path_p + 32) = 108
  set {char}($path_p + 33) = 97
  set {char}($path_p + 34) = 105
  set {char}($path_p + 35) = 110
  set {char}($path_p + 36) = 46
  set {char}($path_p + 37) = 99
  set {char}($path_p + 38) = 115
  set {char}($path_p + 39) = 118
  set {char}($path_p + 40) = 0
  set $path_a = ((char *(*)(unsigned int))0x1001d9c0)(41)
  set {char}($path_a + 0) = 90
  set {char}($path_a + 1) = 58
  set {char}($path_a + 2) = 47
  set {char}($path_a + 3) = 119
  set {char}($path_a + 4) = 111
  set {char}($path_a + 5) = 114
  set {char}($path_a + 6) = 107
  set {char}($path_a + 7) = 47
  set {char}($path_a + 8) = 115
  set {char}($path_a + 9) = 116
  set {char}($path_a + 10) = 97
  set {char}($path_a + 11) = 103
  set {char}($path_a + 12) = 101
  set {char}($path_a + 13) = 49
  set {char}($path_a + 14) = 54
  set {char}($path_a + 15) = 47
  set {char}($path_a + 16) = 117
  set {char}($path_a + 17) = 115
  set {char}($path_a + 18) = 101
  set {char}($path_a + 19) = 114
  set {char}($path_a + 20) = 100
  set {char}($path_a + 21) = 105
  set {char}($path_a + 22) = 99
  set {char}($path_a + 23) = 116
  set {char}($path_a + 24) = 45
  set {char}($path_a + 25) = 97
  set {char}($path_a + 26) = 108
  set {char}($path_a + 27) = 112
  set {char}($path_a + 28) = 104
  set {char}($path_a + 29) = 97
  set {char}($path_a + 30) = 45
  set {char}($path_a + 31) = 119
  set {char}($path_a + 32) = 111
  set {char}($path_a + 33) = 114
  set {char}($path_a + 34) = 108
  set {char}($path_a + 35) = 100
  set {char}($path_a + 36) = 46
  set {char}($path_a + 37) = 99
  set {char}($path_a + 38) = 115
  set {char}($path_a + 39) = 118
  set {char}($path_a + 40) = 0
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 0
  set $path_control = ((char *(*)(unsigned int))0x1001d9c0)(52)
  set {char}($path_control + 0) = 90
  set {char}($path_control + 1) = 58
  set {char}($path_control + 2) = 47
  set {char}($path_control + 3) = 119
  set {char}($path_control + 4) = 111
  set {char}($path_control + 5) = 114
  set {char}($path_control + 6) = 107
  set {char}($path_control + 7) = 47
  set {char}($path_control + 8) = 115
  set {char}($path_control + 9) = 116
  set {char}($path_control + 10) = 97
  set {char}($path_control + 11) = 103
  set {char}($path_control + 12) = 101
  set {char}($path_control + 13) = 49
  set {char}($path_control + 14) = 54
  set {char}($path_control + 15) = 47
  set {char}($path_control + 16) = 117
  set {char}($path_control + 17) = 115
  set {char}($path_control + 18) = 101
  set {char}($path_control + 19) = 114
  set {char}($path_control + 20) = 100
  set {char}($path_control + 21) = 105
  set {char}($path_control + 22) = 99
  set {char}($path_control + 23) = 116
  set {char}($path_control + 24) = 45
  set {char}($path_control + 25) = 108
  set {char}($path_control + 26) = 105
  set {char}($path_control + 27) = 99
  set {char}($path_control + 28) = 101
  set {char}($path_control + 29) = 110
  set {char}($path_control + 30) = 115
  set {char}($path_control + 31) = 101
  set {char}($path_control + 32) = 45
  set {char}($path_control + 33) = 115
  set {char}($path_control + 34) = 116
  set {char}($path_control + 35) = 97
  set {char}($path_control + 36) = 116
  set {char}($path_control + 37) = 101
  set {char}($path_control + 38) = 50
  set {char}($path_control + 39) = 45
  set {char}($path_control + 40) = 99
  set {char}($path_control + 41) = 111
  set {char}($path_control + 42) = 110
  set {char}($path_control + 43) = 116
  set {char}($path_control + 44) = 114
  set {char}($path_control + 45) = 111
  set {char}($path_control + 46) = 108
  set {char}($path_control + 47) = 46
  set {char}($path_control + 48) = 119
  set {char}($path_control + 49) = 97
  set {char}($path_control + 50) = 118
  set {char}($path_control + 51) = 0
  set $path_p_output = ((char *(*)(unsigned int))0x1001d9c0)(46)
  set {char}($path_p_output + 0) = 90
  set {char}($path_p_output + 1) = 58
  set {char}($path_p_output + 2) = 47
  set {char}($path_p_output + 3) = 119
  set {char}($path_p_output + 4) = 111
  set {char}($path_p_output + 5) = 114
  set {char}($path_p_output + 6) = 107
  set {char}($path_p_output + 7) = 47
  set {char}($path_p_output + 8) = 115
  set {char}($path_p_output + 9) = 116
  set {char}($path_p_output + 10) = 97
  set {char}($path_p_output + 11) = 103
  set {char}($path_p_output + 12) = 101
  set {char}($path_p_output + 13) = 49
  set {char}($path_p_output + 14) = 54
  set {char}($path_p_output + 15) = 47
  set {char}($path_p_output + 16) = 117
  set {char}($path_p_output + 17) = 115
  set {char}($path_p_output + 18) = 101
  set {char}($path_p_output + 19) = 114
  set {char}($path_p_output + 20) = 100
  set {char}($path_p_output + 21) = 105
  set {char}($path_p_output + 22) = 99
  set {char}($path_p_output + 23) = 116
  set {char}($path_p_output + 24) = 45
  set {char}($path_p_output + 25) = 108
  set {char}($path_p_output + 26) = 105
  set {char}($path_p_output + 27) = 99
  set {char}($path_p_output + 28) = 101
  set {char}($path_p_output + 29) = 110
  set {char}($path_p_output + 30) = 115
  set {char}($path_p_output + 31) = 101
  set {char}($path_p_output + 32) = 45
  set {char}($path_p_output + 33) = 115
  set {char}($path_p_output + 34) = 116
  set {char}($path_p_output + 35) = 97
  set {char}($path_p_output + 36) = 116
  set {char}($path_p_output + 37) = 101
  set {char}($path_p_output + 38) = 50
  set {char}($path_p_output + 39) = 45
  set {char}($path_p_output + 40) = 112
  set {char}($path_p_output + 41) = 46
  set {char}($path_p_output + 42) = 119
  set {char}($path_p_output + 43) = 97
  set {char}($path_p_output + 44) = 118
  set {char}($path_p_output + 45) = 0
  set $path_a_output = ((char *(*)(unsigned int))0x1001d9c0)(46)
  set {char}($path_a_output + 0) = 90
  set {char}($path_a_output + 1) = 58
  set {char}($path_a_output + 2) = 47
  set {char}($path_a_output + 3) = 119
  set {char}($path_a_output + 4) = 111
  set {char}($path_a_output + 5) = 114
  set {char}($path_a_output + 6) = 107
  set {char}($path_a_output + 7) = 47
  set {char}($path_a_output + 8) = 115
  set {char}($path_a_output + 9) = 116
  set {char}($path_a_output + 10) = 97
  set {char}($path_a_output + 11) = 103
  set {char}($path_a_output + 12) = 101
  set {char}($path_a_output + 13) = 49
  set {char}($path_a_output + 14) = 54
  set {char}($path_a_output + 15) = 47
  set {char}($path_a_output + 16) = 117
  set {char}($path_a_output + 17) = 115
  set {char}($path_a_output + 18) = 101
  set {char}($path_a_output + 19) = 114
  set {char}($path_a_output + 20) = 100
  set {char}($path_a_output + 21) = 105
  set {char}($path_a_output + 22) = 99
  set {char}($path_a_output + 23) = 116
  set {char}($path_a_output + 24) = 45
  set {char}($path_a_output + 25) = 108
  set {char}($path_a_output + 26) = 105
  set {char}($path_a_output + 27) = 99
  set {char}($path_a_output + 28) = 101
  set {char}($path_a_output + 29) = 110
  set {char}($path_a_output + 30) = 115
  set {char}($path_a_output + 31) = 101
  set {char}($path_a_output + 32) = 45
  set {char}($path_a_output + 33) = 115
  set {char}($path_a_output + 34) = 116
  set {char}($path_a_output + 35) = 97
  set {char}($path_a_output + 36) = 116
  set {char}($path_a_output + 37) = 101
  set {char}($path_a_output + 38) = 50
  set {char}($path_a_output + 39) = 45
  set {char}($path_a_output + 40) = 97
  set {char}($path_a_output + 41) = 46
  set {char}($path_a_output + 42) = 119
  set {char}($path_a_output + 43) = 97
  set {char}($path_a_output + 44) = 118
  set {char}($path_a_output + 45) = 0
  set $control = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $path_control, 1, -1, -1, -1, -1, 0, 0)
  printf "USERDICT_SYNTH control dictidx=-1 status=%d\n", $control
  set $load_p = ((short (*)(int, char *))0x10027960)(27, $path_p)
  printf "USERDICT_LOAD type=P index=27 status=%d path=%s\n", $load_p, $path_p
  printf "USERDICT_STATE speaker1_gate=%u index27_pointer=%#x\n", *(unsigned char *)0x100a7489, *(unsigned int *)(0x100a647c + 27 * 4)
  if $load_p == 1
    set $synth_p = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $path_p_output, 1, -1, -1, -1, -1, 27, 0)
    printf "USERDICT_SYNTH type=P index=27 status=%d\n", $synth_p
    set $unload_p = ((short (*)(int))0x10027a80)(27)
    printf "USERDICT_UNLOAD type=P index=27 status=%d\n", $unload_p
  end
  set $load_a = ((short (*)(int, char *))0x10027960)(28, $path_a)
  printf "USERDICT_LOAD type=A index=28 status=%d path=%s\n", $load_a, $path_a
  if $load_a == 1
    set $synth_a = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $path_a_output, 1, -1, -1, -1, -1, 28, 0)
    printf "USERDICT_SYNTH type=A index=28 status=%d\n", $synth_a
    set $unload_a = ((short (*)(int))0x10027a80)(28)
    printf "USERDICT_UNLOAD type=A index=28 status=%d\n", $unload_a
  end
  continue
end

continue
