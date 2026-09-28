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
  set $path_escaped_source_quote = ((char *(*)(unsigned int))0x1001d9c0)(61)
  set {char}($path_escaped_source_quote + 0) = 90
  set {char}($path_escaped_source_quote + 1) = 58
  set {char}($path_escaped_source_quote + 2) = 47
  set {char}($path_escaped_source_quote + 3) = 119
  set {char}($path_escaped_source_quote + 4) = 111
  set {char}($path_escaped_source_quote + 5) = 114
  set {char}($path_escaped_source_quote + 6) = 107
  set {char}($path_escaped_source_quote + 7) = 47
  set {char}($path_escaped_source_quote + 8) = 115
  set {char}($path_escaped_source_quote + 9) = 116
  set {char}($path_escaped_source_quote + 10) = 97
  set {char}($path_escaped_source_quote + 11) = 103
  set {char}($path_escaped_source_quote + 12) = 101
  set {char}($path_escaped_source_quote + 13) = 49
  set {char}($path_escaped_source_quote + 14) = 54
  set {char}($path_escaped_source_quote + 15) = 47
  set {char}($path_escaped_source_quote + 16) = 117
  set {char}($path_escaped_source_quote + 17) = 115
  set {char}($path_escaped_source_quote + 18) = 101
  set {char}($path_escaped_source_quote + 19) = 114
  set {char}($path_escaped_source_quote + 20) = 100
  set {char}($path_escaped_source_quote + 21) = 105
  set {char}($path_escaped_source_quote + 22) = 99
  set {char}($path_escaped_source_quote + 23) = 116
  set {char}($path_escaped_source_quote + 24) = 45
  set {char}($path_escaped_source_quote + 25) = 118
  set {char}($path_escaped_source_quote + 26) = 97
  set {char}($path_escaped_source_quote + 27) = 108
  set {char}($path_escaped_source_quote + 28) = 105
  set {char}($path_escaped_source_quote + 29) = 100
  set {char}($path_escaped_source_quote + 30) = 97
  set {char}($path_escaped_source_quote + 31) = 116
  set {char}($path_escaped_source_quote + 32) = 105
  set {char}($path_escaped_source_quote + 33) = 111
  set {char}($path_escaped_source_quote + 34) = 110
  set {char}($path_escaped_source_quote + 35) = 45
  set {char}($path_escaped_source_quote + 36) = 101
  set {char}($path_escaped_source_quote + 37) = 115
  set {char}($path_escaped_source_quote + 38) = 99
  set {char}($path_escaped_source_quote + 39) = 97
  set {char}($path_escaped_source_quote + 40) = 112
  set {char}($path_escaped_source_quote + 41) = 101
  set {char}($path_escaped_source_quote + 42) = 100
  set {char}($path_escaped_source_quote + 43) = 45
  set {char}($path_escaped_source_quote + 44) = 115
  set {char}($path_escaped_source_quote + 45) = 111
  set {char}($path_escaped_source_quote + 46) = 117
  set {char}($path_escaped_source_quote + 47) = 114
  set {char}($path_escaped_source_quote + 48) = 99
  set {char}($path_escaped_source_quote + 49) = 101
  set {char}($path_escaped_source_quote + 50) = 45
  set {char}($path_escaped_source_quote + 51) = 113
  set {char}($path_escaped_source_quote + 52) = 117
  set {char}($path_escaped_source_quote + 53) = 111
  set {char}($path_escaped_source_quote + 54) = 116
  set {char}($path_escaped_source_quote + 55) = 101
  set {char}($path_escaped_source_quote + 56) = 46
  set {char}($path_escaped_source_quote + 57) = 99
  set {char}($path_escaped_source_quote + 58) = 115
  set {char}($path_escaped_source_quote + 59) = 118
  set {char}($path_escaped_source_quote + 60) = 0
  set $load_escaped_source_quote = ((short (*)(int, char *))0x10027960)(100, $path_escaped_source_quote)
  printf "USERDICT_MATRIX case=escaped-source-quote index=100 load_ax=%d\n", $load_escaped_source_quote
  set $unload_escaped_source_quote = ((short (*)(int))0x10027a80)(100)
  printf "USERDICT_MATRIX case=escaped-source-quote index=100 unload_ax=%d\n", $unload_escaped_source_quote
  if $load_escaped_source_quote != 1
    set $recover_escaped_source_quote = ((short (*)(int, char *))0x10027960)(100, $valid_path)
    printf "USERDICT_MATRIX case=escaped-source-quote index=100 valid_recovery_ax=%d\n", $recover_escaped_source_quote
    if $recover_escaped_source_quote == 1
      set $recover_unload_escaped_source_quote = ((short (*)(int))0x10027a80)(100)
      printf "USERDICT_MATRIX case=escaped-source-quote index=100 recovery_unload_ax=%d\n", $recover_unload_escaped_source_quote
    end
  end
  continue
end

continue
