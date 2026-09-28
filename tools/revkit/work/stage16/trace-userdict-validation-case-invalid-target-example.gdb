set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $valid_path = ((char *(*)(unsigned int))0x1001d9c0)(48)
  set {char}($valid_path + 0) = 90
  set {char}($valid_path + 1) = 58
  set {char}($valid_path + 2) = 47
  set {char}($valid_path + 3) = 119
  set {char}($valid_path + 4) = 111
  set {char}($valid_path + 5) = 114
  set {char}($valid_path + 6) = 107
  set {char}($valid_path + 7) = 47
  set {char}($valid_path + 8) = 115
  set {char}($valid_path + 9) = 116
  set {char}($valid_path + 10) = 97
  set {char}($valid_path + 11) = 103
  set {char}($valid_path + 12) = 101
  set {char}($valid_path + 13) = 49
  set {char}($valid_path + 14) = 54
  set {char}($valid_path + 15) = 47
  set {char}($valid_path + 16) = 117
  set {char}($valid_path + 17) = 115
  set {char}($valid_path + 18) = 101
  set {char}($valid_path + 19) = 114
  set {char}($valid_path + 20) = 100
  set {char}($valid_path + 21) = 105
  set {char}($valid_path + 22) = 99
  set {char}($valid_path + 23) = 116
  set {char}($valid_path + 24) = 45
  set {char}($valid_path + 25) = 118
  set {char}($valid_path + 26) = 97
  set {char}($valid_path + 27) = 108
  set {char}($valid_path + 28) = 105
  set {char}($valid_path + 29) = 100
  set {char}($valid_path + 30) = 97
  set {char}($valid_path + 31) = 116
  set {char}($valid_path + 32) = 105
  set {char}($valid_path + 33) = 111
  set {char}($valid_path + 34) = 110
  set {char}($valid_path + 35) = 45
  set {char}($valid_path + 36) = 112
  set {char}($valid_path + 37) = 108
  set {char}($valid_path + 38) = 97
  set {char}($valid_path + 39) = 105
  set {char}($valid_path + 40) = 110
  set {char}($valid_path + 41) = 45
  set {char}($valid_path + 42) = 112
  set {char}($valid_path + 43) = 46
  set {char}($valid_path + 44) = 99
  set {char}($valid_path + 45) = 115
  set {char}($valid_path + 46) = 118
  set {char}($valid_path + 47) = 0
  set $path_invalid_target_example = ((char *(*)(unsigned int))0x1001d9c0)(63)
  set {char}($path_invalid_target_example + 0) = 90
  set {char}($path_invalid_target_example + 1) = 58
  set {char}($path_invalid_target_example + 2) = 47
  set {char}($path_invalid_target_example + 3) = 119
  set {char}($path_invalid_target_example + 4) = 111
  set {char}($path_invalid_target_example + 5) = 114
  set {char}($path_invalid_target_example + 6) = 107
  set {char}($path_invalid_target_example + 7) = 47
  set {char}($path_invalid_target_example + 8) = 115
  set {char}($path_invalid_target_example + 9) = 116
  set {char}($path_invalid_target_example + 10) = 97
  set {char}($path_invalid_target_example + 11) = 103
  set {char}($path_invalid_target_example + 12) = 101
  set {char}($path_invalid_target_example + 13) = 49
  set {char}($path_invalid_target_example + 14) = 54
  set {char}($path_invalid_target_example + 15) = 47
  set {char}($path_invalid_target_example + 16) = 117
  set {char}($path_invalid_target_example + 17) = 115
  set {char}($path_invalid_target_example + 18) = 101
  set {char}($path_invalid_target_example + 19) = 114
  set {char}($path_invalid_target_example + 20) = 100
  set {char}($path_invalid_target_example + 21) = 105
  set {char}($path_invalid_target_example + 22) = 99
  set {char}($path_invalid_target_example + 23) = 116
  set {char}($path_invalid_target_example + 24) = 45
  set {char}($path_invalid_target_example + 25) = 118
  set {char}($path_invalid_target_example + 26) = 97
  set {char}($path_invalid_target_example + 27) = 108
  set {char}($path_invalid_target_example + 28) = 105
  set {char}($path_invalid_target_example + 29) = 100
  set {char}($path_invalid_target_example + 30) = 97
  set {char}($path_invalid_target_example + 31) = 116
  set {char}($path_invalid_target_example + 32) = 105
  set {char}($path_invalid_target_example + 33) = 111
  set {char}($path_invalid_target_example + 34) = 110
  set {char}($path_invalid_target_example + 35) = 45
  set {char}($path_invalid_target_example + 36) = 105
  set {char}($path_invalid_target_example + 37) = 110
  set {char}($path_invalid_target_example + 38) = 118
  set {char}($path_invalid_target_example + 39) = 97
  set {char}($path_invalid_target_example + 40) = 108
  set {char}($path_invalid_target_example + 41) = 105
  set {char}($path_invalid_target_example + 42) = 100
  set {char}($path_invalid_target_example + 43) = 45
  set {char}($path_invalid_target_example + 44) = 116
  set {char}($path_invalid_target_example + 45) = 97
  set {char}($path_invalid_target_example + 46) = 114
  set {char}($path_invalid_target_example + 47) = 103
  set {char}($path_invalid_target_example + 48) = 101
  set {char}($path_invalid_target_example + 49) = 116
  set {char}($path_invalid_target_example + 50) = 45
  set {char}($path_invalid_target_example + 51) = 101
  set {char}($path_invalid_target_example + 52) = 120
  set {char}($path_invalid_target_example + 53) = 97
  set {char}($path_invalid_target_example + 54) = 109
  set {char}($path_invalid_target_example + 55) = 112
  set {char}($path_invalid_target_example + 56) = 108
  set {char}($path_invalid_target_example + 57) = 101
  set {char}($path_invalid_target_example + 58) = 46
  set {char}($path_invalid_target_example + 59) = 99
  set {char}($path_invalid_target_example + 60) = 115
  set {char}($path_invalid_target_example + 61) = 118
  set {char}($path_invalid_target_example + 62) = 0
  set $load_invalid_target_example = ((short (*)(int, char *))0x10027960)(100, $path_invalid_target_example)
  printf "USERDICT_MATRIX case=invalid-target-example index=100 load_ax=%d\n", $load_invalid_target_example
  set $unload_invalid_target_example = ((short (*)(int))0x10027a80)(100)
  printf "USERDICT_MATRIX case=invalid-target-example index=100 unload_ax=%d\n", $unload_invalid_target_example
  if $load_invalid_target_example != 1
    set $recover_invalid_target_example = ((short (*)(int, char *))0x10027960)(100, $valid_path)
    printf "USERDICT_MATRIX case=invalid-target-example index=100 valid_recovery_ax=%d\n", $recover_invalid_target_example
    if $recover_invalid_target_example == 1
      set $recover_unload_invalid_target_example = ((short (*)(int))0x10027a80)(100)
      printf "USERDICT_MATRIX case=invalid-target-example index=100 recovery_unload_ax=%d\n", $recover_unload_invalid_target_example
    end
  end
  continue
end

continue
