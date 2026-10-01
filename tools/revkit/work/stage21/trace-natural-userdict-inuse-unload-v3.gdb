set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $dict_csv = ((char *(*)(unsigned int))0x1001d9c0)(37)
  set {char}($dict_csv + 0) = 90
  set {char}($dict_csv + 1) = 58
  set {char}($dict_csv + 2) = 47
  set {char}($dict_csv + 3) = 119
  set {char}($dict_csv + 4) = 111
  set {char}($dict_csv + 5) = 114
  set {char}($dict_csv + 6) = 107
  set {char}($dict_csv + 7) = 47
  set {char}($dict_csv + 8) = 115
  set {char}($dict_csv + 9) = 116
  set {char}($dict_csv + 10) = 97
  set {char}($dict_csv + 11) = 103
  set {char}($dict_csv + 12) = 101
  set {char}($dict_csv + 13) = 50
  set {char}($dict_csv + 14) = 49
  set {char}($dict_csv + 15) = 47
  set {char}($dict_csv + 16) = 106
  set {char}($dict_csv + 17) = 97
  set {char}($dict_csv + 18) = 109
  set {char}($dict_csv + 19) = 101
  set {char}($dict_csv + 20) = 115
  set {char}($dict_csv + 21) = 45
  set {char}($dict_csv + 22) = 117
  set {char}($dict_csv + 23) = 115
  set {char}($dict_csv + 24) = 101
  set {char}($dict_csv + 25) = 114
  set {char}($dict_csv + 26) = 100
  set {char}($dict_csv + 27) = 105
  set {char}($dict_csv + 28) = 99
  set {char}($dict_csv + 29) = 116
  set {char}($dict_csv + 30) = 45
  set {char}($dict_csv + 31) = 112
  set {char}($dict_csv + 32) = 46
  set {char}($dict_csv + 33) = 99
  set {char}($dict_csv + 34) = 115
  set {char}($dict_csv + 35) = 118
  set {char}($dict_csv + 36) = 0
  set $control = ((char *(*)(unsigned int))0x1001d9c0)(54)
  set {char}($control + 0) = 90
  set {char}($control + 1) = 58
  set {char}($control + 2) = 47
  set {char}($control + 3) = 119
  set {char}($control + 4) = 111
  set {char}($control + 5) = 114
  set {char}($control + 6) = 107
  set {char}($control + 7) = 47
  set {char}($control + 8) = 115
  set {char}($control + 9) = 116
  set {char}($control + 10) = 97
  set {char}($control + 11) = 103
  set {char}($control + 12) = 101
  set {char}($control + 13) = 50
  set {char}($control + 14) = 49
  set {char}($control + 15) = 47
  set {char}($control + 16) = 110
  set {char}($control + 17) = 97
  set {char}($control + 18) = 116
  set {char}($control + 19) = 117
  set {char}($control + 20) = 114
  set {char}($control + 21) = 97
  set {char}($control + 22) = 108
  set {char}($control + 23) = 45
  set {char}($control + 24) = 117
  set {char}($control + 25) = 115
  set {char}($control + 26) = 101
  set {char}($control + 27) = 114
  set {char}($control + 28) = 100
  set {char}($control + 29) = 105
  set {char}($control + 30) = 99
  set {char}($control + 31) = 116
  set {char}($control + 32) = 45
  set {char}($control + 33) = 105
  set {char}($control + 34) = 110
  set {char}($control + 35) = 117
  set {char}($control + 36) = 115
  set {char}($control + 37) = 101
  set {char}($control + 38) = 45
  set {char}($control + 39) = 99
  set {char}($control + 40) = 111
  set {char}($control + 41) = 110
  set {char}($control + 42) = 116
  set {char}($control + 43) = 114
  set {char}($control + 44) = 111
  set {char}($control + 45) = 108
  set {char}($control + 46) = 45
  set {char}($control + 47) = 118
  set {char}($control + 48) = 51
  set {char}($control + 49) = 46
  set {char}($control + 50) = 119
  set {char}($control + 51) = 97
  set {char}($control + 52) = 118
  set {char}($control + 53) = 0
  set $output = ((char *(*)(unsigned int))0x1001d9c0)(53)
  set {char}($output + 0) = 90
  set {char}($output + 1) = 58
  set {char}($output + 2) = 47
  set {char}($output + 3) = 119
  set {char}($output + 4) = 111
  set {char}($output + 5) = 114
  set {char}($output + 6) = 107
  set {char}($output + 7) = 47
  set {char}($output + 8) = 115
  set {char}($output + 9) = 116
  set {char}($output + 10) = 97
  set {char}($output + 11) = 103
  set {char}($output + 12) = 101
  set {char}($output + 13) = 50
  set {char}($output + 14) = 49
  set {char}($output + 15) = 47
  set {char}($output + 16) = 110
  set {char}($output + 17) = 97
  set {char}($output + 18) = 116
  set {char}($output + 19) = 117
  set {char}($output + 20) = 114
  set {char}($output + 21) = 97
  set {char}($output + 22) = 108
  set {char}($output + 23) = 45
  set {char}($output + 24) = 117
  set {char}($output + 25) = 115
  set {char}($output + 26) = 101
  set {char}($output + 27) = 114
  set {char}($output + 28) = 100
  set {char}($output + 29) = 105
  set {char}($output + 30) = 99
  set {char}($output + 31) = 116
  set {char}($output + 32) = 45
  set {char}($output + 33) = 105
  set {char}($output + 34) = 110
  set {char}($output + 35) = 117
  set {char}($output + 36) = 115
  set {char}($output + 37) = 101
  set {char}($output + 38) = 45
  set {char}($output + 39) = 97
  set {char}($output + 40) = 99
  set {char}($output + 41) = 116
  set {char}($output + 42) = 105
  set {char}($output + 43) = 118
  set {char}($output + 44) = 101
  set {char}($output + 45) = 45
  set {char}($output + 46) = 118
  set {char}($output + 47) = 51
  set {char}($output + 48) = 46
  set {char}($output + 49) = 119
  set {char}($output + 50) = 97
  set {char}($output + 51) = 118
  set {char}($output + 52) = 0
  set $restored = ((char *(*)(unsigned int))0x1001d9c0)(55)
  set {char}($restored + 0) = 90
  set {char}($restored + 1) = 58
  set {char}($restored + 2) = 47
  set {char}($restored + 3) = 119
  set {char}($restored + 4) = 111
  set {char}($restored + 5) = 114
  set {char}($restored + 6) = 107
  set {char}($restored + 7) = 47
  set {char}($restored + 8) = 115
  set {char}($restored + 9) = 116
  set {char}($restored + 10) = 97
  set {char}($restored + 11) = 103
  set {char}($restored + 12) = 101
  set {char}($restored + 13) = 50
  set {char}($restored + 14) = 49
  set {char}($restored + 15) = 47
  set {char}($restored + 16) = 110
  set {char}($restored + 17) = 97
  set {char}($restored + 18) = 116
  set {char}($restored + 19) = 117
  set {char}($restored + 20) = 114
  set {char}($restored + 21) = 97
  set {char}($restored + 22) = 108
  set {char}($restored + 23) = 45
  set {char}($restored + 24) = 117
  set {char}($restored + 25) = 115
  set {char}($restored + 26) = 101
  set {char}($restored + 27) = 114
  set {char}($restored + 28) = 100
  set {char}($restored + 29) = 105
  set {char}($restored + 30) = 99
  set {char}($restored + 31) = 116
  set {char}($restored + 32) = 45
  set {char}($restored + 33) = 105
  set {char}($restored + 34) = 110
  set {char}($restored + 35) = 117
  set {char}($restored + 36) = 115
  set {char}($restored + 37) = 101
  set {char}($restored + 38) = 45
  set {char}($restored + 39) = 114
  set {char}($restored + 40) = 101
  set {char}($restored + 41) = 115
  set {char}($restored + 42) = 116
  set {char}($restored + 43) = 111
  set {char}($restored + 44) = 114
  set {char}($restored + 45) = 101
  set {char}($restored + 46) = 100
  set {char}($restored + 47) = 45
  set {char}($restored + 48) = 118
  set {char}($restored + 49) = 51
  set {char}($restored + 50) = 46
  set {char}($restored + 51) = 119
  set {char}($restored + 52) = 97
  set {char}($restored + 53) = 118
  set {char}($restored + 54) = 0
  set $text = ((char *(*)(unsigned int))0x1001d9c0)(6)
  set {char}($text + 0) = 104
  set {char}($text + 1) = 101
  set {char}($text + 2) = 108
  set {char}($text + 3) = 108
  set {char}($text + 4) = 111
  set {char}($text + 5) = 0
  set $api_return = *(unsigned int *)$esp
  set $incoming_selector = *(int *)($esp + 4)
  set $gate_before = *(unsigned char *)0x100a7489
  set $control_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $control, 1, -1, -1, -1, -1, -1, 0)
  printf "NATURAL_INUSE control_synth=%d\n", $control_result
  set $dict_load = ((short (*)(int, char *))0x10027960)(0, $dict_csv)
  set {unsigned char}0x100a7489 = 1
  printf "NATURAL_INUSE setup selector=%d dict_load=%d gate_before=%d gate_forced=%d return=%p\n", $incoming_selector, $dict_load, $gate_before, *(unsigned char *)0x100a7489, $api_return
  set *(int *)($esp + 4) = 4
  set *(char **)($esp + 8) = $text
  set *(char **)($esp + 12) = $output
  set *(int *)($esp + 16) = 1
  set *(int *)($esp + 20) = -1
  set *(int *)($esp + 24) = -1
  set *(int *)($esp + 28) = -1
  set *(int *)($esp + 32) = -1
  set *(int *)($esp + 36) = 0
  set *(int *)($esp + 40) = 0
  break *0x100260c2
  commands 2
    silent
    set $active_context = *(unsigned int *)0x100a147c
    set $active_dictionary = *(unsigned int *)($active_context + 0x1312c0)
    set $dictionary_pointer = *(unsigned int *)0x100a647c
    printf "NATURAL_INUSE context=%p context_dictionary=%p dictionary=%p\n", $active_context, $active_dictionary, $dictionary_pointer
    set $active_unload = ((short (*)(int))0x10027a80)(0)
    printf "NATURAL_INUSE active_unload=%d\n", $active_unload
    disable 2
    continue
  end
  break *$api_return
  commands 3
    silent
    set $synth_status = $eax
    printf "NATURAL_INUSE synthesis_return=%d\n", $synth_status
    set $idle_unload = ((short (*)(int))0x10027a80)(0)
    printf "NATURAL_INUSE idle_unload=%d\n", $idle_unload
    set {unsigned char}0x100a7489 = $gate_before
    printf "NATURAL_INUSE gate_restored=%d\n", *(unsigned char *)0x100a7489
    set $restored_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $restored, 1, -1, -1, -1, -1, -1, 0)
    printf "NATURAL_INUSE restored_synth=%d\n", $restored_result
    disable 3
    continue
  end
  continue
end

continue
