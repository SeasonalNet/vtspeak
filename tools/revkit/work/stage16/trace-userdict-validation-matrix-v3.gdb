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
  set $path_plain_p = ((char *(*)(unsigned int))0x1001d9c0)(48)
  set {char}($path_plain_p + 0) = 90
  set {char}($path_plain_p + 1) = 58
  set {char}($path_plain_p + 2) = 47
  set {char}($path_plain_p + 3) = 119
  set {char}($path_plain_p + 4) = 111
  set {char}($path_plain_p + 5) = 114
  set {char}($path_plain_p + 6) = 107
  set {char}($path_plain_p + 7) = 47
  set {char}($path_plain_p + 8) = 115
  set {char}($path_plain_p + 9) = 116
  set {char}($path_plain_p + 10) = 97
  set {char}($path_plain_p + 11) = 103
  set {char}($path_plain_p + 12) = 101
  set {char}($path_plain_p + 13) = 49
  set {char}($path_plain_p + 14) = 54
  set {char}($path_plain_p + 15) = 47
  set {char}($path_plain_p + 16) = 117
  set {char}($path_plain_p + 17) = 115
  set {char}($path_plain_p + 18) = 101
  set {char}($path_plain_p + 19) = 114
  set {char}($path_plain_p + 20) = 100
  set {char}($path_plain_p + 21) = 105
  set {char}($path_plain_p + 22) = 99
  set {char}($path_plain_p + 23) = 116
  set {char}($path_plain_p + 24) = 45
  set {char}($path_plain_p + 25) = 118
  set {char}($path_plain_p + 26) = 97
  set {char}($path_plain_p + 27) = 108
  set {char}($path_plain_p + 28) = 105
  set {char}($path_plain_p + 29) = 100
  set {char}($path_plain_p + 30) = 97
  set {char}($path_plain_p + 31) = 116
  set {char}($path_plain_p + 32) = 105
  set {char}($path_plain_p + 33) = 111
  set {char}($path_plain_p + 34) = 110
  set {char}($path_plain_p + 35) = 45
  set {char}($path_plain_p + 36) = 112
  set {char}($path_plain_p + 37) = 108
  set {char}($path_plain_p + 38) = 97
  set {char}($path_plain_p + 39) = 105
  set {char}($path_plain_p + 40) = 110
  set {char}($path_plain_p + 41) = 45
  set {char}($path_plain_p + 42) = 112
  set {char}($path_plain_p + 43) = 46
  set {char}($path_plain_p + 44) = 99
  set {char}($path_plain_p + 45) = 115
  set {char}($path_plain_p + 46) = 118
  set {char}($path_plain_p + 47) = 0
  set $load_plain_p = ((short (*)(int, char *))0x10027960)(100, $path_plain_p)
  printf "USERDICT_MATRIX case=plain-p index=100 load_ax=%d\n", $load_plain_p
  set $unload_plain_p = ((short (*)(int))0x10027a80)(100)
  printf "USERDICT_MATRIX case=plain-p index=100 unload_ax=%d\n", $unload_plain_p
  if $load_plain_p != 1
    set $recover_plain_p = ((short (*)(int, char *))0x10027960)(100, $valid_path)
    printf "USERDICT_MATRIX case=plain-p index=100 valid_recovery_ax=%d\n", $recover_plain_p
    if $recover_plain_p == 1
      set $recover_unload_plain_p = ((short (*)(int))0x10027a80)(100)
      printf "USERDICT_MATRIX case=plain-p index=100 recovery_unload_ax=%d\n", $recover_unload_plain_p
    end
  end
  set $path_four_p = ((char *(*)(unsigned int))0x1001d9c0)(47)
  set {char}($path_four_p + 0) = 90
  set {char}($path_four_p + 1) = 58
  set {char}($path_four_p + 2) = 47
  set {char}($path_four_p + 3) = 119
  set {char}($path_four_p + 4) = 111
  set {char}($path_four_p + 5) = 114
  set {char}($path_four_p + 6) = 107
  set {char}($path_four_p + 7) = 47
  set {char}($path_four_p + 8) = 115
  set {char}($path_four_p + 9) = 116
  set {char}($path_four_p + 10) = 97
  set {char}($path_four_p + 11) = 103
  set {char}($path_four_p + 12) = 101
  set {char}($path_four_p + 13) = 49
  set {char}($path_four_p + 14) = 54
  set {char}($path_four_p + 15) = 47
  set {char}($path_four_p + 16) = 117
  set {char}($path_four_p + 17) = 115
  set {char}($path_four_p + 18) = 101
  set {char}($path_four_p + 19) = 114
  set {char}($path_four_p + 20) = 100
  set {char}($path_four_p + 21) = 105
  set {char}($path_four_p + 22) = 99
  set {char}($path_four_p + 23) = 116
  set {char}($path_four_p + 24) = 45
  set {char}($path_four_p + 25) = 118
  set {char}($path_four_p + 26) = 97
  set {char}($path_four_p + 27) = 108
  set {char}($path_four_p + 28) = 105
  set {char}($path_four_p + 29) = 100
  set {char}($path_four_p + 30) = 97
  set {char}($path_four_p + 31) = 116
  set {char}($path_four_p + 32) = 105
  set {char}($path_four_p + 33) = 111
  set {char}($path_four_p + 34) = 110
  set {char}($path_four_p + 35) = 45
  set {char}($path_four_p + 36) = 102
  set {char}($path_four_p + 37) = 111
  set {char}($path_four_p + 38) = 117
  set {char}($path_four_p + 39) = 114
  set {char}($path_four_p + 40) = 45
  set {char}($path_four_p + 41) = 112
  set {char}($path_four_p + 42) = 46
  set {char}($path_four_p + 43) = 99
  set {char}($path_four_p + 44) = 115
  set {char}($path_four_p + 45) = 118
  set {char}($path_four_p + 46) = 0
  set $load_four_p = ((short (*)(int, char *))0x10027960)(101, $path_four_p)
  printf "USERDICT_MATRIX case=four-p index=101 load_ax=%d\n", $load_four_p
  set $unload_four_p = ((short (*)(int))0x10027a80)(101)
  printf "USERDICT_MATRIX case=four-p index=101 unload_ax=%d\n", $unload_four_p
  if $load_four_p != 1
    set $recover_four_p = ((short (*)(int, char *))0x10027960)(101, $valid_path)
    printf "USERDICT_MATRIX case=four-p index=101 valid_recovery_ax=%d\n", $recover_four_p
    if $recover_four_p == 1
      set $recover_unload_four_p = ((short (*)(int))0x10027a80)(101)
      printf "USERDICT_MATRIX case=four-p index=101 recovery_unload_ax=%d\n", $recover_unload_four_p
    end
  end
  set $path_four_a = ((char *(*)(unsigned int))0x1001d9c0)(47)
  set {char}($path_four_a + 0) = 90
  set {char}($path_four_a + 1) = 58
  set {char}($path_four_a + 2) = 47
  set {char}($path_four_a + 3) = 119
  set {char}($path_four_a + 4) = 111
  set {char}($path_four_a + 5) = 114
  set {char}($path_four_a + 6) = 107
  set {char}($path_four_a + 7) = 47
  set {char}($path_four_a + 8) = 115
  set {char}($path_four_a + 9) = 116
  set {char}($path_four_a + 10) = 97
  set {char}($path_four_a + 11) = 103
  set {char}($path_four_a + 12) = 101
  set {char}($path_four_a + 13) = 49
  set {char}($path_four_a + 14) = 54
  set {char}($path_four_a + 15) = 47
  set {char}($path_four_a + 16) = 117
  set {char}($path_four_a + 17) = 115
  set {char}($path_four_a + 18) = 101
  set {char}($path_four_a + 19) = 114
  set {char}($path_four_a + 20) = 100
  set {char}($path_four_a + 21) = 105
  set {char}($path_four_a + 22) = 99
  set {char}($path_four_a + 23) = 116
  set {char}($path_four_a + 24) = 45
  set {char}($path_four_a + 25) = 118
  set {char}($path_four_a + 26) = 97
  set {char}($path_four_a + 27) = 108
  set {char}($path_four_a + 28) = 105
  set {char}($path_four_a + 29) = 100
  set {char}($path_four_a + 30) = 97
  set {char}($path_four_a + 31) = 116
  set {char}($path_four_a + 32) = 105
  set {char}($path_four_a + 33) = 111
  set {char}($path_four_a + 34) = 110
  set {char}($path_four_a + 35) = 45
  set {char}($path_four_a + 36) = 102
  set {char}($path_four_a + 37) = 111
  set {char}($path_four_a + 38) = 117
  set {char}($path_four_a + 39) = 114
  set {char}($path_four_a + 40) = 45
  set {char}($path_four_a + 41) = 97
  set {char}($path_four_a + 42) = 46
  set {char}($path_four_a + 43) = 99
  set {char}($path_four_a + 44) = 115
  set {char}($path_four_a + 45) = 118
  set {char}($path_four_a + 46) = 0
  set $load_four_a = ((short (*)(int, char *))0x10027960)(102, $path_four_a)
  printf "USERDICT_MATRIX case=four-a index=102 load_ax=%d\n", $load_four_a
  set $unload_four_a = ((short (*)(int))0x10027a80)(102)
  printf "USERDICT_MATRIX case=four-a index=102 unload_ax=%d\n", $unload_four_a
  if $load_four_a != 1
    set $recover_four_a = ((short (*)(int, char *))0x10027960)(102, $valid_path)
    printf "USERDICT_MATRIX case=four-a index=102 valid_recovery_ax=%d\n", $recover_four_a
    if $recover_four_a == 1
      set $recover_unload_four_a = ((short (*)(int))0x10027a80)(102)
      printf "USERDICT_MATRIX case=four-a index=102 recovery_unload_ax=%d\n", $recover_unload_four_a
    end
  end
  set $path_quoted_p = ((char *(*)(unsigned int))0x1001d9c0)(49)
  set {char}($path_quoted_p + 0) = 90
  set {char}($path_quoted_p + 1) = 58
  set {char}($path_quoted_p + 2) = 47
  set {char}($path_quoted_p + 3) = 119
  set {char}($path_quoted_p + 4) = 111
  set {char}($path_quoted_p + 5) = 114
  set {char}($path_quoted_p + 6) = 107
  set {char}($path_quoted_p + 7) = 47
  set {char}($path_quoted_p + 8) = 115
  set {char}($path_quoted_p + 9) = 116
  set {char}($path_quoted_p + 10) = 97
  set {char}($path_quoted_p + 11) = 103
  set {char}($path_quoted_p + 12) = 101
  set {char}($path_quoted_p + 13) = 49
  set {char}($path_quoted_p + 14) = 54
  set {char}($path_quoted_p + 15) = 47
  set {char}($path_quoted_p + 16) = 117
  set {char}($path_quoted_p + 17) = 115
  set {char}($path_quoted_p + 18) = 101
  set {char}($path_quoted_p + 19) = 114
  set {char}($path_quoted_p + 20) = 100
  set {char}($path_quoted_p + 21) = 105
  set {char}($path_quoted_p + 22) = 99
  set {char}($path_quoted_p + 23) = 116
  set {char}($path_quoted_p + 24) = 45
  set {char}($path_quoted_p + 25) = 118
  set {char}($path_quoted_p + 26) = 97
  set {char}($path_quoted_p + 27) = 108
  set {char}($path_quoted_p + 28) = 105
  set {char}($path_quoted_p + 29) = 100
  set {char}($path_quoted_p + 30) = 97
  set {char}($path_quoted_p + 31) = 116
  set {char}($path_quoted_p + 32) = 105
  set {char}($path_quoted_p + 33) = 111
  set {char}($path_quoted_p + 34) = 110
  set {char}($path_quoted_p + 35) = 45
  set {char}($path_quoted_p + 36) = 113
  set {char}($path_quoted_p + 37) = 117
  set {char}($path_quoted_p + 38) = 111
  set {char}($path_quoted_p + 39) = 116
  set {char}($path_quoted_p + 40) = 101
  set {char}($path_quoted_p + 41) = 100
  set {char}($path_quoted_p + 42) = 45
  set {char}($path_quoted_p + 43) = 112
  set {char}($path_quoted_p + 44) = 46
  set {char}($path_quoted_p + 45) = 99
  set {char}($path_quoted_p + 46) = 115
  set {char}($path_quoted_p + 47) = 118
  set {char}($path_quoted_p + 48) = 0
  set $load_quoted_p = ((short (*)(int, char *))0x10027960)(103, $path_quoted_p)
  printf "USERDICT_MATRIX case=quoted-p index=103 load_ax=%d\n", $load_quoted_p
  set $unload_quoted_p = ((short (*)(int))0x10027a80)(103)
  printf "USERDICT_MATRIX case=quoted-p index=103 unload_ax=%d\n", $unload_quoted_p
  if $load_quoted_p != 1
    set $recover_quoted_p = ((short (*)(int, char *))0x10027960)(103, $valid_path)
    printf "USERDICT_MATRIX case=quoted-p index=103 valid_recovery_ax=%d\n", $recover_quoted_p
    if $recover_quoted_p == 1
      set $recover_unload_quoted_p = ((short (*)(int))0x10027a80)(103)
      printf "USERDICT_MATRIX case=quoted-p index=103 recovery_unload_ax=%d\n", $recover_unload_quoted_p
    end
  end
  set $path_two_fields = ((char *(*)(unsigned int))0x1001d9c0)(51)
  set {char}($path_two_fields + 0) = 90
  set {char}($path_two_fields + 1) = 58
  set {char}($path_two_fields + 2) = 47
  set {char}($path_two_fields + 3) = 119
  set {char}($path_two_fields + 4) = 111
  set {char}($path_two_fields + 5) = 114
  set {char}($path_two_fields + 6) = 107
  set {char}($path_two_fields + 7) = 47
  set {char}($path_two_fields + 8) = 115
  set {char}($path_two_fields + 9) = 116
  set {char}($path_two_fields + 10) = 97
  set {char}($path_two_fields + 11) = 103
  set {char}($path_two_fields + 12) = 101
  set {char}($path_two_fields + 13) = 49
  set {char}($path_two_fields + 14) = 54
  set {char}($path_two_fields + 15) = 47
  set {char}($path_two_fields + 16) = 117
  set {char}($path_two_fields + 17) = 115
  set {char}($path_two_fields + 18) = 101
  set {char}($path_two_fields + 19) = 114
  set {char}($path_two_fields + 20) = 100
  set {char}($path_two_fields + 21) = 105
  set {char}($path_two_fields + 22) = 99
  set {char}($path_two_fields + 23) = 116
  set {char}($path_two_fields + 24) = 45
  set {char}($path_two_fields + 25) = 118
  set {char}($path_two_fields + 26) = 97
  set {char}($path_two_fields + 27) = 108
  set {char}($path_two_fields + 28) = 105
  set {char}($path_two_fields + 29) = 100
  set {char}($path_two_fields + 30) = 97
  set {char}($path_two_fields + 31) = 116
  set {char}($path_two_fields + 32) = 105
  set {char}($path_two_fields + 33) = 111
  set {char}($path_two_fields + 34) = 110
  set {char}($path_two_fields + 35) = 45
  set {char}($path_two_fields + 36) = 116
  set {char}($path_two_fields + 37) = 119
  set {char}($path_two_fields + 38) = 111
  set {char}($path_two_fields + 39) = 45
  set {char}($path_two_fields + 40) = 102
  set {char}($path_two_fields + 41) = 105
  set {char}($path_two_fields + 42) = 101
  set {char}($path_two_fields + 43) = 108
  set {char}($path_two_fields + 44) = 100
  set {char}($path_two_fields + 45) = 115
  set {char}($path_two_fields + 46) = 46
  set {char}($path_two_fields + 47) = 99
  set {char}($path_two_fields + 48) = 115
  set {char}($path_two_fields + 49) = 118
  set {char}($path_two_fields + 50) = 0
  set $load_two_fields = ((short (*)(int, char *))0x10027960)(104, $path_two_fields)
  printf "USERDICT_MATRIX case=two-fields index=104 load_ax=%d\n", $load_two_fields
  set $unload_two_fields = ((short (*)(int))0x10027a80)(104)
  printf "USERDICT_MATRIX case=two-fields index=104 unload_ax=%d\n", $unload_two_fields
  if $load_two_fields != 1
    set $recover_two_fields = ((short (*)(int, char *))0x10027960)(104, $valid_path)
    printf "USERDICT_MATRIX case=two-fields index=104 valid_recovery_ax=%d\n", $recover_two_fields
    if $recover_two_fields == 1
      set $recover_unload_two_fields = ((short (*)(int))0x10027a80)(104)
      printf "USERDICT_MATRIX case=two-fields index=104 recovery_unload_ax=%d\n", $recover_unload_two_fields
    end
  end
  set $path_five_fields = ((char *(*)(unsigned int))0x1001d9c0)(52)
  set {char}($path_five_fields + 0) = 90
  set {char}($path_five_fields + 1) = 58
  set {char}($path_five_fields + 2) = 47
  set {char}($path_five_fields + 3) = 119
  set {char}($path_five_fields + 4) = 111
  set {char}($path_five_fields + 5) = 114
  set {char}($path_five_fields + 6) = 107
  set {char}($path_five_fields + 7) = 47
  set {char}($path_five_fields + 8) = 115
  set {char}($path_five_fields + 9) = 116
  set {char}($path_five_fields + 10) = 97
  set {char}($path_five_fields + 11) = 103
  set {char}($path_five_fields + 12) = 101
  set {char}($path_five_fields + 13) = 49
  set {char}($path_five_fields + 14) = 54
  set {char}($path_five_fields + 15) = 47
  set {char}($path_five_fields + 16) = 117
  set {char}($path_five_fields + 17) = 115
  set {char}($path_five_fields + 18) = 101
  set {char}($path_five_fields + 19) = 114
  set {char}($path_five_fields + 20) = 100
  set {char}($path_five_fields + 21) = 105
  set {char}($path_five_fields + 22) = 99
  set {char}($path_five_fields + 23) = 116
  set {char}($path_five_fields + 24) = 45
  set {char}($path_five_fields + 25) = 118
  set {char}($path_five_fields + 26) = 97
  set {char}($path_five_fields + 27) = 108
  set {char}($path_five_fields + 28) = 105
  set {char}($path_five_fields + 29) = 100
  set {char}($path_five_fields + 30) = 97
  set {char}($path_five_fields + 31) = 116
  set {char}($path_five_fields + 32) = 105
  set {char}($path_five_fields + 33) = 111
  set {char}($path_five_fields + 34) = 110
  set {char}($path_five_fields + 35) = 45
  set {char}($path_five_fields + 36) = 102
  set {char}($path_five_fields + 37) = 105
  set {char}($path_five_fields + 38) = 118
  set {char}($path_five_fields + 39) = 101
  set {char}($path_five_fields + 40) = 45
  set {char}($path_five_fields + 41) = 102
  set {char}($path_five_fields + 42) = 105
  set {char}($path_five_fields + 43) = 101
  set {char}($path_five_fields + 44) = 108
  set {char}($path_five_fields + 45) = 100
  set {char}($path_five_fields + 46) = 115
  set {char}($path_five_fields + 47) = 46
  set {char}($path_five_fields + 48) = 99
  set {char}($path_five_fields + 49) = 115
  set {char}($path_five_fields + 50) = 118
  set {char}($path_five_fields + 51) = 0
  set $load_five_fields = ((short (*)(int, char *))0x10027960)(105, $path_five_fields)
  printf "USERDICT_MATRIX case=five-fields index=105 load_ax=%d\n", $load_five_fields
  set $unload_five_fields = ((short (*)(int))0x10027a80)(105)
  printf "USERDICT_MATRIX case=five-fields index=105 unload_ax=%d\n", $unload_five_fields
  if $load_five_fields != 1
    set $recover_five_fields = ((short (*)(int, char *))0x10027960)(105, $valid_path)
    printf "USERDICT_MATRIX case=five-fields index=105 valid_recovery_ax=%d\n", $recover_five_fields
    if $recover_five_fields == 1
      set $recover_unload_five_fields = ((short (*)(int))0x10027a80)(105)
      printf "USERDICT_MATRIX case=five-fields index=105 recovery_unload_ax=%d\n", $recover_unload_five_fields
    end
  end
  set $path_empty_file = ((char *(*)(unsigned int))0x1001d9c0)(51)
  set {char}($path_empty_file + 0) = 90
  set {char}($path_empty_file + 1) = 58
  set {char}($path_empty_file + 2) = 47
  set {char}($path_empty_file + 3) = 119
  set {char}($path_empty_file + 4) = 111
  set {char}($path_empty_file + 5) = 114
  set {char}($path_empty_file + 6) = 107
  set {char}($path_empty_file + 7) = 47
  set {char}($path_empty_file + 8) = 115
  set {char}($path_empty_file + 9) = 116
  set {char}($path_empty_file + 10) = 97
  set {char}($path_empty_file + 11) = 103
  set {char}($path_empty_file + 12) = 101
  set {char}($path_empty_file + 13) = 49
  set {char}($path_empty_file + 14) = 54
  set {char}($path_empty_file + 15) = 47
  set {char}($path_empty_file + 16) = 117
  set {char}($path_empty_file + 17) = 115
  set {char}($path_empty_file + 18) = 101
  set {char}($path_empty_file + 19) = 114
  set {char}($path_empty_file + 20) = 100
  set {char}($path_empty_file + 21) = 105
  set {char}($path_empty_file + 22) = 99
  set {char}($path_empty_file + 23) = 116
  set {char}($path_empty_file + 24) = 45
  set {char}($path_empty_file + 25) = 118
  set {char}($path_empty_file + 26) = 97
  set {char}($path_empty_file + 27) = 108
  set {char}($path_empty_file + 28) = 105
  set {char}($path_empty_file + 29) = 100
  set {char}($path_empty_file + 30) = 97
  set {char}($path_empty_file + 31) = 116
  set {char}($path_empty_file + 32) = 105
  set {char}($path_empty_file + 33) = 111
  set {char}($path_empty_file + 34) = 110
  set {char}($path_empty_file + 35) = 45
  set {char}($path_empty_file + 36) = 101
  set {char}($path_empty_file + 37) = 109
  set {char}($path_empty_file + 38) = 112
  set {char}($path_empty_file + 39) = 116
  set {char}($path_empty_file + 40) = 121
  set {char}($path_empty_file + 41) = 45
  set {char}($path_empty_file + 42) = 102
  set {char}($path_empty_file + 43) = 105
  set {char}($path_empty_file + 44) = 108
  set {char}($path_empty_file + 45) = 101
  set {char}($path_empty_file + 46) = 46
  set {char}($path_empty_file + 47) = 99
  set {char}($path_empty_file + 48) = 115
  set {char}($path_empty_file + 49) = 118
  set {char}($path_empty_file + 50) = 0
  set $load_empty_file = ((short (*)(int, char *))0x10027960)(106, $path_empty_file)
  printf "USERDICT_MATRIX case=empty-file index=106 load_ax=%d\n", $load_empty_file
  set $unload_empty_file = ((short (*)(int))0x10027a80)(106)
  printf "USERDICT_MATRIX case=empty-file index=106 unload_ax=%d\n", $unload_empty_file
  if $load_empty_file != 1
    set $recover_empty_file = ((short (*)(int, char *))0x10027960)(106, $valid_path)
    printf "USERDICT_MATRIX case=empty-file index=106 valid_recovery_ax=%d\n", $recover_empty_file
    if $recover_empty_file == 1
      set $recover_unload_empty_file = ((short (*)(int))0x10027a80)(106)
      printf "USERDICT_MATRIX case=empty-file index=106 recovery_unload_ax=%d\n", $recover_unload_empty_file
    end
  end
  set $path_empty_source = ((char *(*)(unsigned int))0x1001d9c0)(53)
  set {char}($path_empty_source + 0) = 90
  set {char}($path_empty_source + 1) = 58
  set {char}($path_empty_source + 2) = 47
  set {char}($path_empty_source + 3) = 119
  set {char}($path_empty_source + 4) = 111
  set {char}($path_empty_source + 5) = 114
  set {char}($path_empty_source + 6) = 107
  set {char}($path_empty_source + 7) = 47
  set {char}($path_empty_source + 8) = 115
  set {char}($path_empty_source + 9) = 116
  set {char}($path_empty_source + 10) = 97
  set {char}($path_empty_source + 11) = 103
  set {char}($path_empty_source + 12) = 101
  set {char}($path_empty_source + 13) = 49
  set {char}($path_empty_source + 14) = 54
  set {char}($path_empty_source + 15) = 47
  set {char}($path_empty_source + 16) = 117
  set {char}($path_empty_source + 17) = 115
  set {char}($path_empty_source + 18) = 101
  set {char}($path_empty_source + 19) = 114
  set {char}($path_empty_source + 20) = 100
  set {char}($path_empty_source + 21) = 105
  set {char}($path_empty_source + 22) = 99
  set {char}($path_empty_source + 23) = 116
  set {char}($path_empty_source + 24) = 45
  set {char}($path_empty_source + 25) = 118
  set {char}($path_empty_source + 26) = 97
  set {char}($path_empty_source + 27) = 108
  set {char}($path_empty_source + 28) = 105
  set {char}($path_empty_source + 29) = 100
  set {char}($path_empty_source + 30) = 97
  set {char}($path_empty_source + 31) = 116
  set {char}($path_empty_source + 32) = 105
  set {char}($path_empty_source + 33) = 111
  set {char}($path_empty_source + 34) = 110
  set {char}($path_empty_source + 35) = 45
  set {char}($path_empty_source + 36) = 101
  set {char}($path_empty_source + 37) = 109
  set {char}($path_empty_source + 38) = 112
  set {char}($path_empty_source + 39) = 116
  set {char}($path_empty_source + 40) = 121
  set {char}($path_empty_source + 41) = 45
  set {char}($path_empty_source + 42) = 115
  set {char}($path_empty_source + 43) = 111
  set {char}($path_empty_source + 44) = 117
  set {char}($path_empty_source + 45) = 114
  set {char}($path_empty_source + 46) = 99
  set {char}($path_empty_source + 47) = 101
  set {char}($path_empty_source + 48) = 46
  set {char}($path_empty_source + 49) = 99
  set {char}($path_empty_source + 50) = 115
  set {char}($path_empty_source + 51) = 118
  set {char}($path_empty_source + 52) = 0
  set $load_empty_source = ((short (*)(int, char *))0x10027960)(107, $path_empty_source)
  printf "USERDICT_MATRIX case=empty-source index=107 load_ax=%d\n", $load_empty_source
  set $unload_empty_source = ((short (*)(int))0x10027a80)(107)
  printf "USERDICT_MATRIX case=empty-source index=107 unload_ax=%d\n", $unload_empty_source
  if $load_empty_source != 1
    set $recover_empty_source = ((short (*)(int, char *))0x10027960)(107, $valid_path)
    printf "USERDICT_MATRIX case=empty-source index=107 valid_recovery_ax=%d\n", $recover_empty_source
    if $recover_empty_source == 1
      set $recover_unload_empty_source = ((short (*)(int))0x10027a80)(107)
      printf "USERDICT_MATRIX case=empty-source index=107 recovery_unload_ax=%d\n", $recover_unload_empty_source
    end
  end
  set $path_empty_target_p = ((char *(*)(unsigned int))0x1001d9c0)(55)
  set {char}($path_empty_target_p + 0) = 90
  set {char}($path_empty_target_p + 1) = 58
  set {char}($path_empty_target_p + 2) = 47
  set {char}($path_empty_target_p + 3) = 119
  set {char}($path_empty_target_p + 4) = 111
  set {char}($path_empty_target_p + 5) = 114
  set {char}($path_empty_target_p + 6) = 107
  set {char}($path_empty_target_p + 7) = 47
  set {char}($path_empty_target_p + 8) = 115
  set {char}($path_empty_target_p + 9) = 116
  set {char}($path_empty_target_p + 10) = 97
  set {char}($path_empty_target_p + 11) = 103
  set {char}($path_empty_target_p + 12) = 101
  set {char}($path_empty_target_p + 13) = 49
  set {char}($path_empty_target_p + 14) = 54
  set {char}($path_empty_target_p + 15) = 47
  set {char}($path_empty_target_p + 16) = 117
  set {char}($path_empty_target_p + 17) = 115
  set {char}($path_empty_target_p + 18) = 101
  set {char}($path_empty_target_p + 19) = 114
  set {char}($path_empty_target_p + 20) = 100
  set {char}($path_empty_target_p + 21) = 105
  set {char}($path_empty_target_p + 22) = 99
  set {char}($path_empty_target_p + 23) = 116
  set {char}($path_empty_target_p + 24) = 45
  set {char}($path_empty_target_p + 25) = 118
  set {char}($path_empty_target_p + 26) = 97
  set {char}($path_empty_target_p + 27) = 108
  set {char}($path_empty_target_p + 28) = 105
  set {char}($path_empty_target_p + 29) = 100
  set {char}($path_empty_target_p + 30) = 97
  set {char}($path_empty_target_p + 31) = 116
  set {char}($path_empty_target_p + 32) = 105
  set {char}($path_empty_target_p + 33) = 111
  set {char}($path_empty_target_p + 34) = 110
  set {char}($path_empty_target_p + 35) = 45
  set {char}($path_empty_target_p + 36) = 101
  set {char}($path_empty_target_p + 37) = 109
  set {char}($path_empty_target_p + 38) = 112
  set {char}($path_empty_target_p + 39) = 116
  set {char}($path_empty_target_p + 40) = 121
  set {char}($path_empty_target_p + 41) = 45
  set {char}($path_empty_target_p + 42) = 116
  set {char}($path_empty_target_p + 43) = 97
  set {char}($path_empty_target_p + 44) = 114
  set {char}($path_empty_target_p + 45) = 103
  set {char}($path_empty_target_p + 46) = 101
  set {char}($path_empty_target_p + 47) = 116
  set {char}($path_empty_target_p + 48) = 45
  set {char}($path_empty_target_p + 49) = 112
  set {char}($path_empty_target_p + 50) = 46
  set {char}($path_empty_target_p + 51) = 99
  set {char}($path_empty_target_p + 52) = 115
  set {char}($path_empty_target_p + 53) = 118
  set {char}($path_empty_target_p + 54) = 0
  set $load_empty_target_p = ((short (*)(int, char *))0x10027960)(108, $path_empty_target_p)
  printf "USERDICT_MATRIX case=empty-target-p index=108 load_ax=%d\n", $load_empty_target_p
  set $unload_empty_target_p = ((short (*)(int))0x10027a80)(108)
  printf "USERDICT_MATRIX case=empty-target-p index=108 unload_ax=%d\n", $unload_empty_target_p
  if $load_empty_target_p != 1
    set $recover_empty_target_p = ((short (*)(int, char *))0x10027960)(108, $valid_path)
    printf "USERDICT_MATRIX case=empty-target-p index=108 valid_recovery_ax=%d\n", $recover_empty_target_p
    if $recover_empty_target_p == 1
      set $recover_unload_empty_target_p = ((short (*)(int))0x10027a80)(108)
      printf "USERDICT_MATRIX case=empty-target-p index=108 recovery_unload_ax=%d\n", $recover_unload_empty_target_p
    end
  end
  set $path_empty_target_a = ((char *(*)(unsigned int))0x1001d9c0)(55)
  set {char}($path_empty_target_a + 0) = 90
  set {char}($path_empty_target_a + 1) = 58
  set {char}($path_empty_target_a + 2) = 47
  set {char}($path_empty_target_a + 3) = 119
  set {char}($path_empty_target_a + 4) = 111
  set {char}($path_empty_target_a + 5) = 114
  set {char}($path_empty_target_a + 6) = 107
  set {char}($path_empty_target_a + 7) = 47
  set {char}($path_empty_target_a + 8) = 115
  set {char}($path_empty_target_a + 9) = 116
  set {char}($path_empty_target_a + 10) = 97
  set {char}($path_empty_target_a + 11) = 103
  set {char}($path_empty_target_a + 12) = 101
  set {char}($path_empty_target_a + 13) = 49
  set {char}($path_empty_target_a + 14) = 54
  set {char}($path_empty_target_a + 15) = 47
  set {char}($path_empty_target_a + 16) = 117
  set {char}($path_empty_target_a + 17) = 115
  set {char}($path_empty_target_a + 18) = 101
  set {char}($path_empty_target_a + 19) = 114
  set {char}($path_empty_target_a + 20) = 100
  set {char}($path_empty_target_a + 21) = 105
  set {char}($path_empty_target_a + 22) = 99
  set {char}($path_empty_target_a + 23) = 116
  set {char}($path_empty_target_a + 24) = 45
  set {char}($path_empty_target_a + 25) = 118
  set {char}($path_empty_target_a + 26) = 97
  set {char}($path_empty_target_a + 27) = 108
  set {char}($path_empty_target_a + 28) = 105
  set {char}($path_empty_target_a + 29) = 100
  set {char}($path_empty_target_a + 30) = 97
  set {char}($path_empty_target_a + 31) = 116
  set {char}($path_empty_target_a + 32) = 105
  set {char}($path_empty_target_a + 33) = 111
  set {char}($path_empty_target_a + 34) = 110
  set {char}($path_empty_target_a + 35) = 45
  set {char}($path_empty_target_a + 36) = 101
  set {char}($path_empty_target_a + 37) = 109
  set {char}($path_empty_target_a + 38) = 112
  set {char}($path_empty_target_a + 39) = 116
  set {char}($path_empty_target_a + 40) = 121
  set {char}($path_empty_target_a + 41) = 45
  set {char}($path_empty_target_a + 42) = 116
  set {char}($path_empty_target_a + 43) = 97
  set {char}($path_empty_target_a + 44) = 114
  set {char}($path_empty_target_a + 45) = 103
  set {char}($path_empty_target_a + 46) = 101
  set {char}($path_empty_target_a + 47) = 116
  set {char}($path_empty_target_a + 48) = 45
  set {char}($path_empty_target_a + 49) = 97
  set {char}($path_empty_target_a + 50) = 46
  set {char}($path_empty_target_a + 51) = 99
  set {char}($path_empty_target_a + 52) = 115
  set {char}($path_empty_target_a + 53) = 118
  set {char}($path_empty_target_a + 54) = 0
  set $load_empty_target_a = ((short (*)(int, char *))0x10027960)(109, $path_empty_target_a)
  printf "USERDICT_MATRIX case=empty-target-a index=109 load_ax=%d\n", $load_empty_target_a
  set $unload_empty_target_a = ((short (*)(int))0x10027a80)(109)
  printf "USERDICT_MATRIX case=empty-target-a index=109 unload_ax=%d\n", $unload_empty_target_a
  if $load_empty_target_a != 1
    set $recover_empty_target_a = ((short (*)(int, char *))0x10027960)(109, $valid_path)
    printf "USERDICT_MATRIX case=empty-target-a index=109 valid_recovery_ax=%d\n", $recover_empty_target_a
    if $recover_empty_target_a == 1
      set $recover_unload_empty_target_a = ((short (*)(int))0x10027a80)(109)
      printf "USERDICT_MATRIX case=empty-target-a index=109 recovery_unload_ax=%d\n", $recover_unload_empty_target_a
    end
  end
  set $path_invalid_type = ((char *(*)(unsigned int))0x1001d9c0)(53)
  set {char}($path_invalid_type + 0) = 90
  set {char}($path_invalid_type + 1) = 58
  set {char}($path_invalid_type + 2) = 47
  set {char}($path_invalid_type + 3) = 119
  set {char}($path_invalid_type + 4) = 111
  set {char}($path_invalid_type + 5) = 114
  set {char}($path_invalid_type + 6) = 107
  set {char}($path_invalid_type + 7) = 47
  set {char}($path_invalid_type + 8) = 115
  set {char}($path_invalid_type + 9) = 116
  set {char}($path_invalid_type + 10) = 97
  set {char}($path_invalid_type + 11) = 103
  set {char}($path_invalid_type + 12) = 101
  set {char}($path_invalid_type + 13) = 49
  set {char}($path_invalid_type + 14) = 54
  set {char}($path_invalid_type + 15) = 47
  set {char}($path_invalid_type + 16) = 117
  set {char}($path_invalid_type + 17) = 115
  set {char}($path_invalid_type + 18) = 101
  set {char}($path_invalid_type + 19) = 114
  set {char}($path_invalid_type + 20) = 100
  set {char}($path_invalid_type + 21) = 105
  set {char}($path_invalid_type + 22) = 99
  set {char}($path_invalid_type + 23) = 116
  set {char}($path_invalid_type + 24) = 45
  set {char}($path_invalid_type + 25) = 118
  set {char}($path_invalid_type + 26) = 97
  set {char}($path_invalid_type + 27) = 108
  set {char}($path_invalid_type + 28) = 105
  set {char}($path_invalid_type + 29) = 100
  set {char}($path_invalid_type + 30) = 97
  set {char}($path_invalid_type + 31) = 116
  set {char}($path_invalid_type + 32) = 105
  set {char}($path_invalid_type + 33) = 111
  set {char}($path_invalid_type + 34) = 110
  set {char}($path_invalid_type + 35) = 45
  set {char}($path_invalid_type + 36) = 105
  set {char}($path_invalid_type + 37) = 110
  set {char}($path_invalid_type + 38) = 118
  set {char}($path_invalid_type + 39) = 97
  set {char}($path_invalid_type + 40) = 108
  set {char}($path_invalid_type + 41) = 105
  set {char}($path_invalid_type + 42) = 100
  set {char}($path_invalid_type + 43) = 45
  set {char}($path_invalid_type + 44) = 116
  set {char}($path_invalid_type + 45) = 121
  set {char}($path_invalid_type + 46) = 112
  set {char}($path_invalid_type + 47) = 101
  set {char}($path_invalid_type + 48) = 46
  set {char}($path_invalid_type + 49) = 99
  set {char}($path_invalid_type + 50) = 115
  set {char}($path_invalid_type + 51) = 118
  set {char}($path_invalid_type + 52) = 0
  set $load_invalid_type = ((short (*)(int, char *))0x10027960)(110, $path_invalid_type)
  printf "USERDICT_MATRIX case=invalid-type index=110 load_ax=%d\n", $load_invalid_type
  set $unload_invalid_type = ((short (*)(int))0x10027a80)(110)
  printf "USERDICT_MATRIX case=invalid-type index=110 unload_ax=%d\n", $unload_invalid_type
  if $load_invalid_type != 1
    set $recover_invalid_type = ((short (*)(int, char *))0x10027960)(110, $valid_path)
    printf "USERDICT_MATRIX case=invalid-type index=110 valid_recovery_ax=%d\n", $recover_invalid_type
    if $recover_invalid_type == 1
      set $recover_unload_invalid_type = ((short (*)(int))0x10027a80)(110)
      printf "USERDICT_MATRIX case=invalid-type index=110 recovery_unload_ax=%d\n", $recover_unload_invalid_type
    end
  end
  set $path_long_type = ((char *(*)(unsigned int))0x1001d9c0)(50)
  set {char}($path_long_type + 0) = 90
  set {char}($path_long_type + 1) = 58
  set {char}($path_long_type + 2) = 47
  set {char}($path_long_type + 3) = 119
  set {char}($path_long_type + 4) = 111
  set {char}($path_long_type + 5) = 114
  set {char}($path_long_type + 6) = 107
  set {char}($path_long_type + 7) = 47
  set {char}($path_long_type + 8) = 115
  set {char}($path_long_type + 9) = 116
  set {char}($path_long_type + 10) = 97
  set {char}($path_long_type + 11) = 103
  set {char}($path_long_type + 12) = 101
  set {char}($path_long_type + 13) = 49
  set {char}($path_long_type + 14) = 54
  set {char}($path_long_type + 15) = 47
  set {char}($path_long_type + 16) = 117
  set {char}($path_long_type + 17) = 115
  set {char}($path_long_type + 18) = 101
  set {char}($path_long_type + 19) = 114
  set {char}($path_long_type + 20) = 100
  set {char}($path_long_type + 21) = 105
  set {char}($path_long_type + 22) = 99
  set {char}($path_long_type + 23) = 116
  set {char}($path_long_type + 24) = 45
  set {char}($path_long_type + 25) = 118
  set {char}($path_long_type + 26) = 97
  set {char}($path_long_type + 27) = 108
  set {char}($path_long_type + 28) = 105
  set {char}($path_long_type + 29) = 100
  set {char}($path_long_type + 30) = 97
  set {char}($path_long_type + 31) = 116
  set {char}($path_long_type + 32) = 105
  set {char}($path_long_type + 33) = 111
  set {char}($path_long_type + 34) = 110
  set {char}($path_long_type + 35) = 45
  set {char}($path_long_type + 36) = 108
  set {char}($path_long_type + 37) = 111
  set {char}($path_long_type + 38) = 110
  set {char}($path_long_type + 39) = 103
  set {char}($path_long_type + 40) = 45
  set {char}($path_long_type + 41) = 116
  set {char}($path_long_type + 42) = 121
  set {char}($path_long_type + 43) = 112
  set {char}($path_long_type + 44) = 101
  set {char}($path_long_type + 45) = 46
  set {char}($path_long_type + 46) = 99
  set {char}($path_long_type + 47) = 115
  set {char}($path_long_type + 48) = 118
  set {char}($path_long_type + 49) = 0
  set $load_long_type = ((short (*)(int, char *))0x10027960)(111, $path_long_type)
  printf "USERDICT_MATRIX case=long-type index=111 load_ax=%d\n", $load_long_type
  set $unload_long_type = ((short (*)(int))0x10027a80)(111)
  printf "USERDICT_MATRIX case=long-type index=111 unload_ax=%d\n", $unload_long_type
  if $load_long_type != 1
    set $recover_long_type = ((short (*)(int, char *))0x10027960)(111, $valid_path)
    printf "USERDICT_MATRIX case=long-type index=111 valid_recovery_ax=%d\n", $recover_long_type
    if $recover_long_type == 1
      set $recover_unload_long_type = ((short (*)(int))0x10027a80)(111)
      printf "USERDICT_MATRIX case=long-type index=111 recovery_unload_ax=%d\n", $recover_unload_long_type
    end
  end
  set $path_trailing_empty = ((char *(*)(unsigned int))0x1001d9c0)(55)
  set {char}($path_trailing_empty + 0) = 90
  set {char}($path_trailing_empty + 1) = 58
  set {char}($path_trailing_empty + 2) = 47
  set {char}($path_trailing_empty + 3) = 119
  set {char}($path_trailing_empty + 4) = 111
  set {char}($path_trailing_empty + 5) = 114
  set {char}($path_trailing_empty + 6) = 107
  set {char}($path_trailing_empty + 7) = 47
  set {char}($path_trailing_empty + 8) = 115
  set {char}($path_trailing_empty + 9) = 116
  set {char}($path_trailing_empty + 10) = 97
  set {char}($path_trailing_empty + 11) = 103
  set {char}($path_trailing_empty + 12) = 101
  set {char}($path_trailing_empty + 13) = 49
  set {char}($path_trailing_empty + 14) = 54
  set {char}($path_trailing_empty + 15) = 47
  set {char}($path_trailing_empty + 16) = 117
  set {char}($path_trailing_empty + 17) = 115
  set {char}($path_trailing_empty + 18) = 101
  set {char}($path_trailing_empty + 19) = 114
  set {char}($path_trailing_empty + 20) = 100
  set {char}($path_trailing_empty + 21) = 105
  set {char}($path_trailing_empty + 22) = 99
  set {char}($path_trailing_empty + 23) = 116
  set {char}($path_trailing_empty + 24) = 45
  set {char}($path_trailing_empty + 25) = 118
  set {char}($path_trailing_empty + 26) = 97
  set {char}($path_trailing_empty + 27) = 108
  set {char}($path_trailing_empty + 28) = 105
  set {char}($path_trailing_empty + 29) = 100
  set {char}($path_trailing_empty + 30) = 97
  set {char}($path_trailing_empty + 31) = 116
  set {char}($path_trailing_empty + 32) = 105
  set {char}($path_trailing_empty + 33) = 111
  set {char}($path_trailing_empty + 34) = 110
  set {char}($path_trailing_empty + 35) = 45
  set {char}($path_trailing_empty + 36) = 116
  set {char}($path_trailing_empty + 37) = 114
  set {char}($path_trailing_empty + 38) = 97
  set {char}($path_trailing_empty + 39) = 105
  set {char}($path_trailing_empty + 40) = 108
  set {char}($path_trailing_empty + 41) = 105
  set {char}($path_trailing_empty + 42) = 110
  set {char}($path_trailing_empty + 43) = 103
  set {char}($path_trailing_empty + 44) = 45
  set {char}($path_trailing_empty + 45) = 101
  set {char}($path_trailing_empty + 46) = 109
  set {char}($path_trailing_empty + 47) = 112
  set {char}($path_trailing_empty + 48) = 116
  set {char}($path_trailing_empty + 49) = 121
  set {char}($path_trailing_empty + 50) = 46
  set {char}($path_trailing_empty + 51) = 99
  set {char}($path_trailing_empty + 52) = 115
  set {char}($path_trailing_empty + 53) = 118
  set {char}($path_trailing_empty + 54) = 0
  set $load_trailing_empty = ((short (*)(int, char *))0x10027960)(112, $path_trailing_empty)
  printf "USERDICT_MATRIX case=trailing-empty index=112 load_ax=%d\n", $load_trailing_empty
  set $unload_trailing_empty = ((short (*)(int))0x10027a80)(112)
  printf "USERDICT_MATRIX case=trailing-empty index=112 unload_ax=%d\n", $unload_trailing_empty
  if $load_trailing_empty != 1
    set $recover_trailing_empty = ((short (*)(int, char *))0x10027960)(112, $valid_path)
    printf "USERDICT_MATRIX case=trailing-empty index=112 valid_recovery_ax=%d\n", $recover_trailing_empty
    if $recover_trailing_empty == 1
      set $recover_unload_trailing_empty = ((short (*)(int))0x10027a80)(112)
      printf "USERDICT_MATRIX case=trailing-empty index=112 recovery_unload_ax=%d\n", $recover_unload_trailing_empty
    end
  end
  set $path_crlf_p = ((char *(*)(unsigned int))0x1001d9c0)(47)
  set {char}($path_crlf_p + 0) = 90
  set {char}($path_crlf_p + 1) = 58
  set {char}($path_crlf_p + 2) = 47
  set {char}($path_crlf_p + 3) = 119
  set {char}($path_crlf_p + 4) = 111
  set {char}($path_crlf_p + 5) = 114
  set {char}($path_crlf_p + 6) = 107
  set {char}($path_crlf_p + 7) = 47
  set {char}($path_crlf_p + 8) = 115
  set {char}($path_crlf_p + 9) = 116
  set {char}($path_crlf_p + 10) = 97
  set {char}($path_crlf_p + 11) = 103
  set {char}($path_crlf_p + 12) = 101
  set {char}($path_crlf_p + 13) = 49
  set {char}($path_crlf_p + 14) = 54
  set {char}($path_crlf_p + 15) = 47
  set {char}($path_crlf_p + 16) = 117
  set {char}($path_crlf_p + 17) = 115
  set {char}($path_crlf_p + 18) = 101
  set {char}($path_crlf_p + 19) = 114
  set {char}($path_crlf_p + 20) = 100
  set {char}($path_crlf_p + 21) = 105
  set {char}($path_crlf_p + 22) = 99
  set {char}($path_crlf_p + 23) = 116
  set {char}($path_crlf_p + 24) = 45
  set {char}($path_crlf_p + 25) = 118
  set {char}($path_crlf_p + 26) = 97
  set {char}($path_crlf_p + 27) = 108
  set {char}($path_crlf_p + 28) = 105
  set {char}($path_crlf_p + 29) = 100
  set {char}($path_crlf_p + 30) = 97
  set {char}($path_crlf_p + 31) = 116
  set {char}($path_crlf_p + 32) = 105
  set {char}($path_crlf_p + 33) = 111
  set {char}($path_crlf_p + 34) = 110
  set {char}($path_crlf_p + 35) = 45
  set {char}($path_crlf_p + 36) = 99
  set {char}($path_crlf_p + 37) = 114
  set {char}($path_crlf_p + 38) = 108
  set {char}($path_crlf_p + 39) = 102
  set {char}($path_crlf_p + 40) = 45
  set {char}($path_crlf_p + 41) = 112
  set {char}($path_crlf_p + 42) = 46
  set {char}($path_crlf_p + 43) = 99
  set {char}($path_crlf_p + 44) = 115
  set {char}($path_crlf_p + 45) = 118
  set {char}($path_crlf_p + 46) = 0
  set $load_crlf_p = ((short (*)(int, char *))0x10027960)(113, $path_crlf_p)
  printf "USERDICT_MATRIX case=crlf-p index=113 load_ax=%d\n", $load_crlf_p
  set $unload_crlf_p = ((short (*)(int))0x10027a80)(113)
  printf "USERDICT_MATRIX case=crlf-p index=113 unload_ax=%d\n", $unload_crlf_p
  if $load_crlf_p != 1
    set $recover_crlf_p = ((short (*)(int, char *))0x10027960)(113, $valid_path)
    printf "USERDICT_MATRIX case=crlf-p index=113 valid_recovery_ax=%d\n", $recover_crlf_p
    if $recover_crlf_p == 1
      set $recover_unload_crlf_p = ((short (*)(int))0x10027a80)(113)
      printf "USERDICT_MATRIX case=crlf-p index=113 recovery_unload_ax=%d\n", $recover_unload_crlf_p
    end
  end
  continue
end

continue
