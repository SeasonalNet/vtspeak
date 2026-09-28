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
  set $path_plain-p = ((char *(*)(unsigned int))0x1001d9c0)(48)
  set {char}($path_plain-p + 0) = 90
  set {char}($path_plain-p + 1) = 58
  set {char}($path_plain-p + 2) = 47
  set {char}($path_plain-p + 3) = 119
  set {char}($path_plain-p + 4) = 111
  set {char}($path_plain-p + 5) = 114
  set {char}($path_plain-p + 6) = 107
  set {char}($path_plain-p + 7) = 47
  set {char}($path_plain-p + 8) = 115
  set {char}($path_plain-p + 9) = 116
  set {char}($path_plain-p + 10) = 97
  set {char}($path_plain-p + 11) = 103
  set {char}($path_plain-p + 12) = 101
  set {char}($path_plain-p + 13) = 49
  set {char}($path_plain-p + 14) = 54
  set {char}($path_plain-p + 15) = 47
  set {char}($path_plain-p + 16) = 117
  set {char}($path_plain-p + 17) = 115
  set {char}($path_plain-p + 18) = 101
  set {char}($path_plain-p + 19) = 114
  set {char}($path_plain-p + 20) = 100
  set {char}($path_plain-p + 21) = 105
  set {char}($path_plain-p + 22) = 99
  set {char}($path_plain-p + 23) = 116
  set {char}($path_plain-p + 24) = 45
  set {char}($path_plain-p + 25) = 118
  set {char}($path_plain-p + 26) = 97
  set {char}($path_plain-p + 27) = 108
  set {char}($path_plain-p + 28) = 105
  set {char}($path_plain-p + 29) = 100
  set {char}($path_plain-p + 30) = 97
  set {char}($path_plain-p + 31) = 116
  set {char}($path_plain-p + 32) = 105
  set {char}($path_plain-p + 33) = 111
  set {char}($path_plain-p + 34) = 110
  set {char}($path_plain-p + 35) = 45
  set {char}($path_plain-p + 36) = 112
  set {char}($path_plain-p + 37) = 108
  set {char}($path_plain-p + 38) = 97
  set {char}($path_plain-p + 39) = 105
  set {char}($path_plain-p + 40) = 110
  set {char}($path_plain-p + 41) = 45
  set {char}($path_plain-p + 42) = 112
  set {char}($path_plain-p + 43) = 46
  set {char}($path_plain-p + 44) = 99
  set {char}($path_plain-p + 45) = 115
  set {char}($path_plain-p + 46) = 118
  set {char}($path_plain-p + 47) = 0
  set $load_plain-p = ((short (*)(int, char *))0x10027960)(100, $path_plain-p)
  printf "USERDICT_MATRIX case=plain-p index=100 load_ax=%d\n", $load_plain-p
  set $unload_plain-p = ((short (*)(int))0x10027a80)(100)
  printf "USERDICT_MATRIX case=plain-p index=100 unload_ax=%d\n", $unload_plain-p
  if $load_plain-p != 1
    set $recover_plain-p = ((short (*)(int, char *))0x10027960)(100, $valid_path)
    printf "USERDICT_MATRIX case=plain-p index=100 valid_recovery_ax=%d\n", $recover_plain-p
    if $recover_plain-p == 1
      set $recover_unload_plain-p = ((short (*)(int))0x10027a80)(100)
      printf "USERDICT_MATRIX case=plain-p index=100 recovery_unload_ax=%d\n", $recover_unload_plain-p
    end
  end
  set $path_four-p = ((char *(*)(unsigned int))0x1001d9c0)(47)
  set {char}($path_four-p + 0) = 90
  set {char}($path_four-p + 1) = 58
  set {char}($path_four-p + 2) = 47
  set {char}($path_four-p + 3) = 119
  set {char}($path_four-p + 4) = 111
  set {char}($path_four-p + 5) = 114
  set {char}($path_four-p + 6) = 107
  set {char}($path_four-p + 7) = 47
  set {char}($path_four-p + 8) = 115
  set {char}($path_four-p + 9) = 116
  set {char}($path_four-p + 10) = 97
  set {char}($path_four-p + 11) = 103
  set {char}($path_four-p + 12) = 101
  set {char}($path_four-p + 13) = 49
  set {char}($path_four-p + 14) = 54
  set {char}($path_four-p + 15) = 47
  set {char}($path_four-p + 16) = 117
  set {char}($path_four-p + 17) = 115
  set {char}($path_four-p + 18) = 101
  set {char}($path_four-p + 19) = 114
  set {char}($path_four-p + 20) = 100
  set {char}($path_four-p + 21) = 105
  set {char}($path_four-p + 22) = 99
  set {char}($path_four-p + 23) = 116
  set {char}($path_four-p + 24) = 45
  set {char}($path_four-p + 25) = 118
  set {char}($path_four-p + 26) = 97
  set {char}($path_four-p + 27) = 108
  set {char}($path_four-p + 28) = 105
  set {char}($path_four-p + 29) = 100
  set {char}($path_four-p + 30) = 97
  set {char}($path_four-p + 31) = 116
  set {char}($path_four-p + 32) = 105
  set {char}($path_four-p + 33) = 111
  set {char}($path_four-p + 34) = 110
  set {char}($path_four-p + 35) = 45
  set {char}($path_four-p + 36) = 102
  set {char}($path_four-p + 37) = 111
  set {char}($path_four-p + 38) = 117
  set {char}($path_four-p + 39) = 114
  set {char}($path_four-p + 40) = 45
  set {char}($path_four-p + 41) = 112
  set {char}($path_four-p + 42) = 46
  set {char}($path_four-p + 43) = 99
  set {char}($path_four-p + 44) = 115
  set {char}($path_four-p + 45) = 118
  set {char}($path_four-p + 46) = 0
  set $load_four-p = ((short (*)(int, char *))0x10027960)(101, $path_four-p)
  printf "USERDICT_MATRIX case=four-p index=101 load_ax=%d\n", $load_four-p
  set $unload_four-p = ((short (*)(int))0x10027a80)(101)
  printf "USERDICT_MATRIX case=four-p index=101 unload_ax=%d\n", $unload_four-p
  if $load_four-p != 1
    set $recover_four-p = ((short (*)(int, char *))0x10027960)(101, $valid_path)
    printf "USERDICT_MATRIX case=four-p index=101 valid_recovery_ax=%d\n", $recover_four-p
    if $recover_four-p == 1
      set $recover_unload_four-p = ((short (*)(int))0x10027a80)(101)
      printf "USERDICT_MATRIX case=four-p index=101 recovery_unload_ax=%d\n", $recover_unload_four-p
    end
  end
  set $path_four-a = ((char *(*)(unsigned int))0x1001d9c0)(47)
  set {char}($path_four-a + 0) = 90
  set {char}($path_four-a + 1) = 58
  set {char}($path_four-a + 2) = 47
  set {char}($path_four-a + 3) = 119
  set {char}($path_four-a + 4) = 111
  set {char}($path_four-a + 5) = 114
  set {char}($path_four-a + 6) = 107
  set {char}($path_four-a + 7) = 47
  set {char}($path_four-a + 8) = 115
  set {char}($path_four-a + 9) = 116
  set {char}($path_four-a + 10) = 97
  set {char}($path_four-a + 11) = 103
  set {char}($path_four-a + 12) = 101
  set {char}($path_four-a + 13) = 49
  set {char}($path_four-a + 14) = 54
  set {char}($path_four-a + 15) = 47
  set {char}($path_four-a + 16) = 117
  set {char}($path_four-a + 17) = 115
  set {char}($path_four-a + 18) = 101
  set {char}($path_four-a + 19) = 114
  set {char}($path_four-a + 20) = 100
  set {char}($path_four-a + 21) = 105
  set {char}($path_four-a + 22) = 99
  set {char}($path_four-a + 23) = 116
  set {char}($path_four-a + 24) = 45
  set {char}($path_four-a + 25) = 118
  set {char}($path_four-a + 26) = 97
  set {char}($path_four-a + 27) = 108
  set {char}($path_four-a + 28) = 105
  set {char}($path_four-a + 29) = 100
  set {char}($path_four-a + 30) = 97
  set {char}($path_four-a + 31) = 116
  set {char}($path_four-a + 32) = 105
  set {char}($path_four-a + 33) = 111
  set {char}($path_four-a + 34) = 110
  set {char}($path_four-a + 35) = 45
  set {char}($path_four-a + 36) = 102
  set {char}($path_four-a + 37) = 111
  set {char}($path_four-a + 38) = 117
  set {char}($path_four-a + 39) = 114
  set {char}($path_four-a + 40) = 45
  set {char}($path_four-a + 41) = 97
  set {char}($path_four-a + 42) = 46
  set {char}($path_four-a + 43) = 99
  set {char}($path_four-a + 44) = 115
  set {char}($path_four-a + 45) = 118
  set {char}($path_four-a + 46) = 0
  set $load_four-a = ((short (*)(int, char *))0x10027960)(102, $path_four-a)
  printf "USERDICT_MATRIX case=four-a index=102 load_ax=%d\n", $load_four-a
  set $unload_four-a = ((short (*)(int))0x10027a80)(102)
  printf "USERDICT_MATRIX case=four-a index=102 unload_ax=%d\n", $unload_four-a
  if $load_four-a != 1
    set $recover_four-a = ((short (*)(int, char *))0x10027960)(102, $valid_path)
    printf "USERDICT_MATRIX case=four-a index=102 valid_recovery_ax=%d\n", $recover_four-a
    if $recover_four-a == 1
      set $recover_unload_four-a = ((short (*)(int))0x10027a80)(102)
      printf "USERDICT_MATRIX case=four-a index=102 recovery_unload_ax=%d\n", $recover_unload_four-a
    end
  end
  set $path_quoted-p = ((char *(*)(unsigned int))0x1001d9c0)(49)
  set {char}($path_quoted-p + 0) = 90
  set {char}($path_quoted-p + 1) = 58
  set {char}($path_quoted-p + 2) = 47
  set {char}($path_quoted-p + 3) = 119
  set {char}($path_quoted-p + 4) = 111
  set {char}($path_quoted-p + 5) = 114
  set {char}($path_quoted-p + 6) = 107
  set {char}($path_quoted-p + 7) = 47
  set {char}($path_quoted-p + 8) = 115
  set {char}($path_quoted-p + 9) = 116
  set {char}($path_quoted-p + 10) = 97
  set {char}($path_quoted-p + 11) = 103
  set {char}($path_quoted-p + 12) = 101
  set {char}($path_quoted-p + 13) = 49
  set {char}($path_quoted-p + 14) = 54
  set {char}($path_quoted-p + 15) = 47
  set {char}($path_quoted-p + 16) = 117
  set {char}($path_quoted-p + 17) = 115
  set {char}($path_quoted-p + 18) = 101
  set {char}($path_quoted-p + 19) = 114
  set {char}($path_quoted-p + 20) = 100
  set {char}($path_quoted-p + 21) = 105
  set {char}($path_quoted-p + 22) = 99
  set {char}($path_quoted-p + 23) = 116
  set {char}($path_quoted-p + 24) = 45
  set {char}($path_quoted-p + 25) = 118
  set {char}($path_quoted-p + 26) = 97
  set {char}($path_quoted-p + 27) = 108
  set {char}($path_quoted-p + 28) = 105
  set {char}($path_quoted-p + 29) = 100
  set {char}($path_quoted-p + 30) = 97
  set {char}($path_quoted-p + 31) = 116
  set {char}($path_quoted-p + 32) = 105
  set {char}($path_quoted-p + 33) = 111
  set {char}($path_quoted-p + 34) = 110
  set {char}($path_quoted-p + 35) = 45
  set {char}($path_quoted-p + 36) = 113
  set {char}($path_quoted-p + 37) = 117
  set {char}($path_quoted-p + 38) = 111
  set {char}($path_quoted-p + 39) = 116
  set {char}($path_quoted-p + 40) = 101
  set {char}($path_quoted-p + 41) = 100
  set {char}($path_quoted-p + 42) = 45
  set {char}($path_quoted-p + 43) = 112
  set {char}($path_quoted-p + 44) = 46
  set {char}($path_quoted-p + 45) = 99
  set {char}($path_quoted-p + 46) = 115
  set {char}($path_quoted-p + 47) = 118
  set {char}($path_quoted-p + 48) = 0
  set $load_quoted-p = ((short (*)(int, char *))0x10027960)(103, $path_quoted-p)
  printf "USERDICT_MATRIX case=quoted-p index=103 load_ax=%d\n", $load_quoted-p
  set $unload_quoted-p = ((short (*)(int))0x10027a80)(103)
  printf "USERDICT_MATRIX case=quoted-p index=103 unload_ax=%d\n", $unload_quoted-p
  if $load_quoted-p != 1
    set $recover_quoted-p = ((short (*)(int, char *))0x10027960)(103, $valid_path)
    printf "USERDICT_MATRIX case=quoted-p index=103 valid_recovery_ax=%d\n", $recover_quoted-p
    if $recover_quoted-p == 1
      set $recover_unload_quoted-p = ((short (*)(int))0x10027a80)(103)
      printf "USERDICT_MATRIX case=quoted-p index=103 recovery_unload_ax=%d\n", $recover_unload_quoted-p
    end
  end
  set $path_two-fields = ((char *(*)(unsigned int))0x1001d9c0)(51)
  set {char}($path_two-fields + 0) = 90
  set {char}($path_two-fields + 1) = 58
  set {char}($path_two-fields + 2) = 47
  set {char}($path_two-fields + 3) = 119
  set {char}($path_two-fields + 4) = 111
  set {char}($path_two-fields + 5) = 114
  set {char}($path_two-fields + 6) = 107
  set {char}($path_two-fields + 7) = 47
  set {char}($path_two-fields + 8) = 115
  set {char}($path_two-fields + 9) = 116
  set {char}($path_two-fields + 10) = 97
  set {char}($path_two-fields + 11) = 103
  set {char}($path_two-fields + 12) = 101
  set {char}($path_two-fields + 13) = 49
  set {char}($path_two-fields + 14) = 54
  set {char}($path_two-fields + 15) = 47
  set {char}($path_two-fields + 16) = 117
  set {char}($path_two-fields + 17) = 115
  set {char}($path_two-fields + 18) = 101
  set {char}($path_two-fields + 19) = 114
  set {char}($path_two-fields + 20) = 100
  set {char}($path_two-fields + 21) = 105
  set {char}($path_two-fields + 22) = 99
  set {char}($path_two-fields + 23) = 116
  set {char}($path_two-fields + 24) = 45
  set {char}($path_two-fields + 25) = 118
  set {char}($path_two-fields + 26) = 97
  set {char}($path_two-fields + 27) = 108
  set {char}($path_two-fields + 28) = 105
  set {char}($path_two-fields + 29) = 100
  set {char}($path_two-fields + 30) = 97
  set {char}($path_two-fields + 31) = 116
  set {char}($path_two-fields + 32) = 105
  set {char}($path_two-fields + 33) = 111
  set {char}($path_two-fields + 34) = 110
  set {char}($path_two-fields + 35) = 45
  set {char}($path_two-fields + 36) = 116
  set {char}($path_two-fields + 37) = 119
  set {char}($path_two-fields + 38) = 111
  set {char}($path_two-fields + 39) = 45
  set {char}($path_two-fields + 40) = 102
  set {char}($path_two-fields + 41) = 105
  set {char}($path_two-fields + 42) = 101
  set {char}($path_two-fields + 43) = 108
  set {char}($path_two-fields + 44) = 100
  set {char}($path_two-fields + 45) = 115
  set {char}($path_two-fields + 46) = 46
  set {char}($path_two-fields + 47) = 99
  set {char}($path_two-fields + 48) = 115
  set {char}($path_two-fields + 49) = 118
  set {char}($path_two-fields + 50) = 0
  set $load_two-fields = ((short (*)(int, char *))0x10027960)(104, $path_two-fields)
  printf "USERDICT_MATRIX case=two-fields index=104 load_ax=%d\n", $load_two-fields
  set $unload_two-fields = ((short (*)(int))0x10027a80)(104)
  printf "USERDICT_MATRIX case=two-fields index=104 unload_ax=%d\n", $unload_two-fields
  if $load_two-fields != 1
    set $recover_two-fields = ((short (*)(int, char *))0x10027960)(104, $valid_path)
    printf "USERDICT_MATRIX case=two-fields index=104 valid_recovery_ax=%d\n", $recover_two-fields
    if $recover_two-fields == 1
      set $recover_unload_two-fields = ((short (*)(int))0x10027a80)(104)
      printf "USERDICT_MATRIX case=two-fields index=104 recovery_unload_ax=%d\n", $recover_unload_two-fields
    end
  end
  set $path_five-fields = ((char *(*)(unsigned int))0x1001d9c0)(52)
  set {char}($path_five-fields + 0) = 90
  set {char}($path_five-fields + 1) = 58
  set {char}($path_five-fields + 2) = 47
  set {char}($path_five-fields + 3) = 119
  set {char}($path_five-fields + 4) = 111
  set {char}($path_five-fields + 5) = 114
  set {char}($path_five-fields + 6) = 107
  set {char}($path_five-fields + 7) = 47
  set {char}($path_five-fields + 8) = 115
  set {char}($path_five-fields + 9) = 116
  set {char}($path_five-fields + 10) = 97
  set {char}($path_five-fields + 11) = 103
  set {char}($path_five-fields + 12) = 101
  set {char}($path_five-fields + 13) = 49
  set {char}($path_five-fields + 14) = 54
  set {char}($path_five-fields + 15) = 47
  set {char}($path_five-fields + 16) = 117
  set {char}($path_five-fields + 17) = 115
  set {char}($path_five-fields + 18) = 101
  set {char}($path_five-fields + 19) = 114
  set {char}($path_five-fields + 20) = 100
  set {char}($path_five-fields + 21) = 105
  set {char}($path_five-fields + 22) = 99
  set {char}($path_five-fields + 23) = 116
  set {char}($path_five-fields + 24) = 45
  set {char}($path_five-fields + 25) = 118
  set {char}($path_five-fields + 26) = 97
  set {char}($path_five-fields + 27) = 108
  set {char}($path_five-fields + 28) = 105
  set {char}($path_five-fields + 29) = 100
  set {char}($path_five-fields + 30) = 97
  set {char}($path_five-fields + 31) = 116
  set {char}($path_five-fields + 32) = 105
  set {char}($path_five-fields + 33) = 111
  set {char}($path_five-fields + 34) = 110
  set {char}($path_five-fields + 35) = 45
  set {char}($path_five-fields + 36) = 102
  set {char}($path_five-fields + 37) = 105
  set {char}($path_five-fields + 38) = 118
  set {char}($path_five-fields + 39) = 101
  set {char}($path_five-fields + 40) = 45
  set {char}($path_five-fields + 41) = 102
  set {char}($path_five-fields + 42) = 105
  set {char}($path_five-fields + 43) = 101
  set {char}($path_five-fields + 44) = 108
  set {char}($path_five-fields + 45) = 100
  set {char}($path_five-fields + 46) = 115
  set {char}($path_five-fields + 47) = 46
  set {char}($path_five-fields + 48) = 99
  set {char}($path_five-fields + 49) = 115
  set {char}($path_five-fields + 50) = 118
  set {char}($path_five-fields + 51) = 0
  set $load_five-fields = ((short (*)(int, char *))0x10027960)(105, $path_five-fields)
  printf "USERDICT_MATRIX case=five-fields index=105 load_ax=%d\n", $load_five-fields
  set $unload_five-fields = ((short (*)(int))0x10027a80)(105)
  printf "USERDICT_MATRIX case=five-fields index=105 unload_ax=%d\n", $unload_five-fields
  if $load_five-fields != 1
    set $recover_five-fields = ((short (*)(int, char *))0x10027960)(105, $valid_path)
    printf "USERDICT_MATRIX case=five-fields index=105 valid_recovery_ax=%d\n", $recover_five-fields
    if $recover_five-fields == 1
      set $recover_unload_five-fields = ((short (*)(int))0x10027a80)(105)
      printf "USERDICT_MATRIX case=five-fields index=105 recovery_unload_ax=%d\n", $recover_unload_five-fields
    end
  end
  set $path_empty-file = ((char *(*)(unsigned int))0x1001d9c0)(51)
  set {char}($path_empty-file + 0) = 90
  set {char}($path_empty-file + 1) = 58
  set {char}($path_empty-file + 2) = 47
  set {char}($path_empty-file + 3) = 119
  set {char}($path_empty-file + 4) = 111
  set {char}($path_empty-file + 5) = 114
  set {char}($path_empty-file + 6) = 107
  set {char}($path_empty-file + 7) = 47
  set {char}($path_empty-file + 8) = 115
  set {char}($path_empty-file + 9) = 116
  set {char}($path_empty-file + 10) = 97
  set {char}($path_empty-file + 11) = 103
  set {char}($path_empty-file + 12) = 101
  set {char}($path_empty-file + 13) = 49
  set {char}($path_empty-file + 14) = 54
  set {char}($path_empty-file + 15) = 47
  set {char}($path_empty-file + 16) = 117
  set {char}($path_empty-file + 17) = 115
  set {char}($path_empty-file + 18) = 101
  set {char}($path_empty-file + 19) = 114
  set {char}($path_empty-file + 20) = 100
  set {char}($path_empty-file + 21) = 105
  set {char}($path_empty-file + 22) = 99
  set {char}($path_empty-file + 23) = 116
  set {char}($path_empty-file + 24) = 45
  set {char}($path_empty-file + 25) = 118
  set {char}($path_empty-file + 26) = 97
  set {char}($path_empty-file + 27) = 108
  set {char}($path_empty-file + 28) = 105
  set {char}($path_empty-file + 29) = 100
  set {char}($path_empty-file + 30) = 97
  set {char}($path_empty-file + 31) = 116
  set {char}($path_empty-file + 32) = 105
  set {char}($path_empty-file + 33) = 111
  set {char}($path_empty-file + 34) = 110
  set {char}($path_empty-file + 35) = 45
  set {char}($path_empty-file + 36) = 101
  set {char}($path_empty-file + 37) = 109
  set {char}($path_empty-file + 38) = 112
  set {char}($path_empty-file + 39) = 116
  set {char}($path_empty-file + 40) = 121
  set {char}($path_empty-file + 41) = 45
  set {char}($path_empty-file + 42) = 102
  set {char}($path_empty-file + 43) = 105
  set {char}($path_empty-file + 44) = 108
  set {char}($path_empty-file + 45) = 101
  set {char}($path_empty-file + 46) = 46
  set {char}($path_empty-file + 47) = 99
  set {char}($path_empty-file + 48) = 115
  set {char}($path_empty-file + 49) = 118
  set {char}($path_empty-file + 50) = 0
  set $load_empty-file = ((short (*)(int, char *))0x10027960)(106, $path_empty-file)
  printf "USERDICT_MATRIX case=empty-file index=106 load_ax=%d\n", $load_empty-file
  set $unload_empty-file = ((short (*)(int))0x10027a80)(106)
  printf "USERDICT_MATRIX case=empty-file index=106 unload_ax=%d\n", $unload_empty-file
  if $load_empty-file != 1
    set $recover_empty-file = ((short (*)(int, char *))0x10027960)(106, $valid_path)
    printf "USERDICT_MATRIX case=empty-file index=106 valid_recovery_ax=%d\n", $recover_empty-file
    if $recover_empty-file == 1
      set $recover_unload_empty-file = ((short (*)(int))0x10027a80)(106)
      printf "USERDICT_MATRIX case=empty-file index=106 recovery_unload_ax=%d\n", $recover_unload_empty-file
    end
  end
  set $path_empty-source = ((char *(*)(unsigned int))0x1001d9c0)(53)
  set {char}($path_empty-source + 0) = 90
  set {char}($path_empty-source + 1) = 58
  set {char}($path_empty-source + 2) = 47
  set {char}($path_empty-source + 3) = 119
  set {char}($path_empty-source + 4) = 111
  set {char}($path_empty-source + 5) = 114
  set {char}($path_empty-source + 6) = 107
  set {char}($path_empty-source + 7) = 47
  set {char}($path_empty-source + 8) = 115
  set {char}($path_empty-source + 9) = 116
  set {char}($path_empty-source + 10) = 97
  set {char}($path_empty-source + 11) = 103
  set {char}($path_empty-source + 12) = 101
  set {char}($path_empty-source + 13) = 49
  set {char}($path_empty-source + 14) = 54
  set {char}($path_empty-source + 15) = 47
  set {char}($path_empty-source + 16) = 117
  set {char}($path_empty-source + 17) = 115
  set {char}($path_empty-source + 18) = 101
  set {char}($path_empty-source + 19) = 114
  set {char}($path_empty-source + 20) = 100
  set {char}($path_empty-source + 21) = 105
  set {char}($path_empty-source + 22) = 99
  set {char}($path_empty-source + 23) = 116
  set {char}($path_empty-source + 24) = 45
  set {char}($path_empty-source + 25) = 118
  set {char}($path_empty-source + 26) = 97
  set {char}($path_empty-source + 27) = 108
  set {char}($path_empty-source + 28) = 105
  set {char}($path_empty-source + 29) = 100
  set {char}($path_empty-source + 30) = 97
  set {char}($path_empty-source + 31) = 116
  set {char}($path_empty-source + 32) = 105
  set {char}($path_empty-source + 33) = 111
  set {char}($path_empty-source + 34) = 110
  set {char}($path_empty-source + 35) = 45
  set {char}($path_empty-source + 36) = 101
  set {char}($path_empty-source + 37) = 109
  set {char}($path_empty-source + 38) = 112
  set {char}($path_empty-source + 39) = 116
  set {char}($path_empty-source + 40) = 121
  set {char}($path_empty-source + 41) = 45
  set {char}($path_empty-source + 42) = 115
  set {char}($path_empty-source + 43) = 111
  set {char}($path_empty-source + 44) = 117
  set {char}($path_empty-source + 45) = 114
  set {char}($path_empty-source + 46) = 99
  set {char}($path_empty-source + 47) = 101
  set {char}($path_empty-source + 48) = 46
  set {char}($path_empty-source + 49) = 99
  set {char}($path_empty-source + 50) = 115
  set {char}($path_empty-source + 51) = 118
  set {char}($path_empty-source + 52) = 0
  set $load_empty-source = ((short (*)(int, char *))0x10027960)(107, $path_empty-source)
  printf "USERDICT_MATRIX case=empty-source index=107 load_ax=%d\n", $load_empty-source
  set $unload_empty-source = ((short (*)(int))0x10027a80)(107)
  printf "USERDICT_MATRIX case=empty-source index=107 unload_ax=%d\n", $unload_empty-source
  if $load_empty-source != 1
    set $recover_empty-source = ((short (*)(int, char *))0x10027960)(107, $valid_path)
    printf "USERDICT_MATRIX case=empty-source index=107 valid_recovery_ax=%d\n", $recover_empty-source
    if $recover_empty-source == 1
      set $recover_unload_empty-source = ((short (*)(int))0x10027a80)(107)
      printf "USERDICT_MATRIX case=empty-source index=107 recovery_unload_ax=%d\n", $recover_unload_empty-source
    end
  end
  set $path_empty-target-p = ((char *(*)(unsigned int))0x1001d9c0)(55)
  set {char}($path_empty-target-p + 0) = 90
  set {char}($path_empty-target-p + 1) = 58
  set {char}($path_empty-target-p + 2) = 47
  set {char}($path_empty-target-p + 3) = 119
  set {char}($path_empty-target-p + 4) = 111
  set {char}($path_empty-target-p + 5) = 114
  set {char}($path_empty-target-p + 6) = 107
  set {char}($path_empty-target-p + 7) = 47
  set {char}($path_empty-target-p + 8) = 115
  set {char}($path_empty-target-p + 9) = 116
  set {char}($path_empty-target-p + 10) = 97
  set {char}($path_empty-target-p + 11) = 103
  set {char}($path_empty-target-p + 12) = 101
  set {char}($path_empty-target-p + 13) = 49
  set {char}($path_empty-target-p + 14) = 54
  set {char}($path_empty-target-p + 15) = 47
  set {char}($path_empty-target-p + 16) = 117
  set {char}($path_empty-target-p + 17) = 115
  set {char}($path_empty-target-p + 18) = 101
  set {char}($path_empty-target-p + 19) = 114
  set {char}($path_empty-target-p + 20) = 100
  set {char}($path_empty-target-p + 21) = 105
  set {char}($path_empty-target-p + 22) = 99
  set {char}($path_empty-target-p + 23) = 116
  set {char}($path_empty-target-p + 24) = 45
  set {char}($path_empty-target-p + 25) = 118
  set {char}($path_empty-target-p + 26) = 97
  set {char}($path_empty-target-p + 27) = 108
  set {char}($path_empty-target-p + 28) = 105
  set {char}($path_empty-target-p + 29) = 100
  set {char}($path_empty-target-p + 30) = 97
  set {char}($path_empty-target-p + 31) = 116
  set {char}($path_empty-target-p + 32) = 105
  set {char}($path_empty-target-p + 33) = 111
  set {char}($path_empty-target-p + 34) = 110
  set {char}($path_empty-target-p + 35) = 45
  set {char}($path_empty-target-p + 36) = 101
  set {char}($path_empty-target-p + 37) = 109
  set {char}($path_empty-target-p + 38) = 112
  set {char}($path_empty-target-p + 39) = 116
  set {char}($path_empty-target-p + 40) = 121
  set {char}($path_empty-target-p + 41) = 45
  set {char}($path_empty-target-p + 42) = 116
  set {char}($path_empty-target-p + 43) = 97
  set {char}($path_empty-target-p + 44) = 114
  set {char}($path_empty-target-p + 45) = 103
  set {char}($path_empty-target-p + 46) = 101
  set {char}($path_empty-target-p + 47) = 116
  set {char}($path_empty-target-p + 48) = 45
  set {char}($path_empty-target-p + 49) = 112
  set {char}($path_empty-target-p + 50) = 46
  set {char}($path_empty-target-p + 51) = 99
  set {char}($path_empty-target-p + 52) = 115
  set {char}($path_empty-target-p + 53) = 118
  set {char}($path_empty-target-p + 54) = 0
  set $load_empty-target-p = ((short (*)(int, char *))0x10027960)(108, $path_empty-target-p)
  printf "USERDICT_MATRIX case=empty-target-p index=108 load_ax=%d\n", $load_empty-target-p
  set $unload_empty-target-p = ((short (*)(int))0x10027a80)(108)
  printf "USERDICT_MATRIX case=empty-target-p index=108 unload_ax=%d\n", $unload_empty-target-p
  if $load_empty-target-p != 1
    set $recover_empty-target-p = ((short (*)(int, char *))0x10027960)(108, $valid_path)
    printf "USERDICT_MATRIX case=empty-target-p index=108 valid_recovery_ax=%d\n", $recover_empty-target-p
    if $recover_empty-target-p == 1
      set $recover_unload_empty-target-p = ((short (*)(int))0x10027a80)(108)
      printf "USERDICT_MATRIX case=empty-target-p index=108 recovery_unload_ax=%d\n", $recover_unload_empty-target-p
    end
  end
  set $path_empty-target-a = ((char *(*)(unsigned int))0x1001d9c0)(55)
  set {char}($path_empty-target-a + 0) = 90
  set {char}($path_empty-target-a + 1) = 58
  set {char}($path_empty-target-a + 2) = 47
  set {char}($path_empty-target-a + 3) = 119
  set {char}($path_empty-target-a + 4) = 111
  set {char}($path_empty-target-a + 5) = 114
  set {char}($path_empty-target-a + 6) = 107
  set {char}($path_empty-target-a + 7) = 47
  set {char}($path_empty-target-a + 8) = 115
  set {char}($path_empty-target-a + 9) = 116
  set {char}($path_empty-target-a + 10) = 97
  set {char}($path_empty-target-a + 11) = 103
  set {char}($path_empty-target-a + 12) = 101
  set {char}($path_empty-target-a + 13) = 49
  set {char}($path_empty-target-a + 14) = 54
  set {char}($path_empty-target-a + 15) = 47
  set {char}($path_empty-target-a + 16) = 117
  set {char}($path_empty-target-a + 17) = 115
  set {char}($path_empty-target-a + 18) = 101
  set {char}($path_empty-target-a + 19) = 114
  set {char}($path_empty-target-a + 20) = 100
  set {char}($path_empty-target-a + 21) = 105
  set {char}($path_empty-target-a + 22) = 99
  set {char}($path_empty-target-a + 23) = 116
  set {char}($path_empty-target-a + 24) = 45
  set {char}($path_empty-target-a + 25) = 118
  set {char}($path_empty-target-a + 26) = 97
  set {char}($path_empty-target-a + 27) = 108
  set {char}($path_empty-target-a + 28) = 105
  set {char}($path_empty-target-a + 29) = 100
  set {char}($path_empty-target-a + 30) = 97
  set {char}($path_empty-target-a + 31) = 116
  set {char}($path_empty-target-a + 32) = 105
  set {char}($path_empty-target-a + 33) = 111
  set {char}($path_empty-target-a + 34) = 110
  set {char}($path_empty-target-a + 35) = 45
  set {char}($path_empty-target-a + 36) = 101
  set {char}($path_empty-target-a + 37) = 109
  set {char}($path_empty-target-a + 38) = 112
  set {char}($path_empty-target-a + 39) = 116
  set {char}($path_empty-target-a + 40) = 121
  set {char}($path_empty-target-a + 41) = 45
  set {char}($path_empty-target-a + 42) = 116
  set {char}($path_empty-target-a + 43) = 97
  set {char}($path_empty-target-a + 44) = 114
  set {char}($path_empty-target-a + 45) = 103
  set {char}($path_empty-target-a + 46) = 101
  set {char}($path_empty-target-a + 47) = 116
  set {char}($path_empty-target-a + 48) = 45
  set {char}($path_empty-target-a + 49) = 97
  set {char}($path_empty-target-a + 50) = 46
  set {char}($path_empty-target-a + 51) = 99
  set {char}($path_empty-target-a + 52) = 115
  set {char}($path_empty-target-a + 53) = 118
  set {char}($path_empty-target-a + 54) = 0
  set $load_empty-target-a = ((short (*)(int, char *))0x10027960)(109, $path_empty-target-a)
  printf "USERDICT_MATRIX case=empty-target-a index=109 load_ax=%d\n", $load_empty-target-a
  set $unload_empty-target-a = ((short (*)(int))0x10027a80)(109)
  printf "USERDICT_MATRIX case=empty-target-a index=109 unload_ax=%d\n", $unload_empty-target-a
  if $load_empty-target-a != 1
    set $recover_empty-target-a = ((short (*)(int, char *))0x10027960)(109, $valid_path)
    printf "USERDICT_MATRIX case=empty-target-a index=109 valid_recovery_ax=%d\n", $recover_empty-target-a
    if $recover_empty-target-a == 1
      set $recover_unload_empty-target-a = ((short (*)(int))0x10027a80)(109)
      printf "USERDICT_MATRIX case=empty-target-a index=109 recovery_unload_ax=%d\n", $recover_unload_empty-target-a
    end
  end
  set $path_invalid-type = ((char *(*)(unsigned int))0x1001d9c0)(53)
  set {char}($path_invalid-type + 0) = 90
  set {char}($path_invalid-type + 1) = 58
  set {char}($path_invalid-type + 2) = 47
  set {char}($path_invalid-type + 3) = 119
  set {char}($path_invalid-type + 4) = 111
  set {char}($path_invalid-type + 5) = 114
  set {char}($path_invalid-type + 6) = 107
  set {char}($path_invalid-type + 7) = 47
  set {char}($path_invalid-type + 8) = 115
  set {char}($path_invalid-type + 9) = 116
  set {char}($path_invalid-type + 10) = 97
  set {char}($path_invalid-type + 11) = 103
  set {char}($path_invalid-type + 12) = 101
  set {char}($path_invalid-type + 13) = 49
  set {char}($path_invalid-type + 14) = 54
  set {char}($path_invalid-type + 15) = 47
  set {char}($path_invalid-type + 16) = 117
  set {char}($path_invalid-type + 17) = 115
  set {char}($path_invalid-type + 18) = 101
  set {char}($path_invalid-type + 19) = 114
  set {char}($path_invalid-type + 20) = 100
  set {char}($path_invalid-type + 21) = 105
  set {char}($path_invalid-type + 22) = 99
  set {char}($path_invalid-type + 23) = 116
  set {char}($path_invalid-type + 24) = 45
  set {char}($path_invalid-type + 25) = 118
  set {char}($path_invalid-type + 26) = 97
  set {char}($path_invalid-type + 27) = 108
  set {char}($path_invalid-type + 28) = 105
  set {char}($path_invalid-type + 29) = 100
  set {char}($path_invalid-type + 30) = 97
  set {char}($path_invalid-type + 31) = 116
  set {char}($path_invalid-type + 32) = 105
  set {char}($path_invalid-type + 33) = 111
  set {char}($path_invalid-type + 34) = 110
  set {char}($path_invalid-type + 35) = 45
  set {char}($path_invalid-type + 36) = 105
  set {char}($path_invalid-type + 37) = 110
  set {char}($path_invalid-type + 38) = 118
  set {char}($path_invalid-type + 39) = 97
  set {char}($path_invalid-type + 40) = 108
  set {char}($path_invalid-type + 41) = 105
  set {char}($path_invalid-type + 42) = 100
  set {char}($path_invalid-type + 43) = 45
  set {char}($path_invalid-type + 44) = 116
  set {char}($path_invalid-type + 45) = 121
  set {char}($path_invalid-type + 46) = 112
  set {char}($path_invalid-type + 47) = 101
  set {char}($path_invalid-type + 48) = 46
  set {char}($path_invalid-type + 49) = 99
  set {char}($path_invalid-type + 50) = 115
  set {char}($path_invalid-type + 51) = 118
  set {char}($path_invalid-type + 52) = 0
  set $load_invalid-type = ((short (*)(int, char *))0x10027960)(110, $path_invalid-type)
  printf "USERDICT_MATRIX case=invalid-type index=110 load_ax=%d\n", $load_invalid-type
  set $unload_invalid-type = ((short (*)(int))0x10027a80)(110)
  printf "USERDICT_MATRIX case=invalid-type index=110 unload_ax=%d\n", $unload_invalid-type
  if $load_invalid-type != 1
    set $recover_invalid-type = ((short (*)(int, char *))0x10027960)(110, $valid_path)
    printf "USERDICT_MATRIX case=invalid-type index=110 valid_recovery_ax=%d\n", $recover_invalid-type
    if $recover_invalid-type == 1
      set $recover_unload_invalid-type = ((short (*)(int))0x10027a80)(110)
      printf "USERDICT_MATRIX case=invalid-type index=110 recovery_unload_ax=%d\n", $recover_unload_invalid-type
    end
  end
  set $path_long-type = ((char *(*)(unsigned int))0x1001d9c0)(50)
  set {char}($path_long-type + 0) = 90
  set {char}($path_long-type + 1) = 58
  set {char}($path_long-type + 2) = 47
  set {char}($path_long-type + 3) = 119
  set {char}($path_long-type + 4) = 111
  set {char}($path_long-type + 5) = 114
  set {char}($path_long-type + 6) = 107
  set {char}($path_long-type + 7) = 47
  set {char}($path_long-type + 8) = 115
  set {char}($path_long-type + 9) = 116
  set {char}($path_long-type + 10) = 97
  set {char}($path_long-type + 11) = 103
  set {char}($path_long-type + 12) = 101
  set {char}($path_long-type + 13) = 49
  set {char}($path_long-type + 14) = 54
  set {char}($path_long-type + 15) = 47
  set {char}($path_long-type + 16) = 117
  set {char}($path_long-type + 17) = 115
  set {char}($path_long-type + 18) = 101
  set {char}($path_long-type + 19) = 114
  set {char}($path_long-type + 20) = 100
  set {char}($path_long-type + 21) = 105
  set {char}($path_long-type + 22) = 99
  set {char}($path_long-type + 23) = 116
  set {char}($path_long-type + 24) = 45
  set {char}($path_long-type + 25) = 118
  set {char}($path_long-type + 26) = 97
  set {char}($path_long-type + 27) = 108
  set {char}($path_long-type + 28) = 105
  set {char}($path_long-type + 29) = 100
  set {char}($path_long-type + 30) = 97
  set {char}($path_long-type + 31) = 116
  set {char}($path_long-type + 32) = 105
  set {char}($path_long-type + 33) = 111
  set {char}($path_long-type + 34) = 110
  set {char}($path_long-type + 35) = 45
  set {char}($path_long-type + 36) = 108
  set {char}($path_long-type + 37) = 111
  set {char}($path_long-type + 38) = 110
  set {char}($path_long-type + 39) = 103
  set {char}($path_long-type + 40) = 45
  set {char}($path_long-type + 41) = 116
  set {char}($path_long-type + 42) = 121
  set {char}($path_long-type + 43) = 112
  set {char}($path_long-type + 44) = 101
  set {char}($path_long-type + 45) = 46
  set {char}($path_long-type + 46) = 99
  set {char}($path_long-type + 47) = 115
  set {char}($path_long-type + 48) = 118
  set {char}($path_long-type + 49) = 0
  set $load_long-type = ((short (*)(int, char *))0x10027960)(111, $path_long-type)
  printf "USERDICT_MATRIX case=long-type index=111 load_ax=%d\n", $load_long-type
  set $unload_long-type = ((short (*)(int))0x10027a80)(111)
  printf "USERDICT_MATRIX case=long-type index=111 unload_ax=%d\n", $unload_long-type
  if $load_long-type != 1
    set $recover_long-type = ((short (*)(int, char *))0x10027960)(111, $valid_path)
    printf "USERDICT_MATRIX case=long-type index=111 valid_recovery_ax=%d\n", $recover_long-type
    if $recover_long-type == 1
      set $recover_unload_long-type = ((short (*)(int))0x10027a80)(111)
      printf "USERDICT_MATRIX case=long-type index=111 recovery_unload_ax=%d\n", $recover_unload_long-type
    end
  end
  set $path_trailing-empty = ((char *(*)(unsigned int))0x1001d9c0)(55)
  set {char}($path_trailing-empty + 0) = 90
  set {char}($path_trailing-empty + 1) = 58
  set {char}($path_trailing-empty + 2) = 47
  set {char}($path_trailing-empty + 3) = 119
  set {char}($path_trailing-empty + 4) = 111
  set {char}($path_trailing-empty + 5) = 114
  set {char}($path_trailing-empty + 6) = 107
  set {char}($path_trailing-empty + 7) = 47
  set {char}($path_trailing-empty + 8) = 115
  set {char}($path_trailing-empty + 9) = 116
  set {char}($path_trailing-empty + 10) = 97
  set {char}($path_trailing-empty + 11) = 103
  set {char}($path_trailing-empty + 12) = 101
  set {char}($path_trailing-empty + 13) = 49
  set {char}($path_trailing-empty + 14) = 54
  set {char}($path_trailing-empty + 15) = 47
  set {char}($path_trailing-empty + 16) = 117
  set {char}($path_trailing-empty + 17) = 115
  set {char}($path_trailing-empty + 18) = 101
  set {char}($path_trailing-empty + 19) = 114
  set {char}($path_trailing-empty + 20) = 100
  set {char}($path_trailing-empty + 21) = 105
  set {char}($path_trailing-empty + 22) = 99
  set {char}($path_trailing-empty + 23) = 116
  set {char}($path_trailing-empty + 24) = 45
  set {char}($path_trailing-empty + 25) = 118
  set {char}($path_trailing-empty + 26) = 97
  set {char}($path_trailing-empty + 27) = 108
  set {char}($path_trailing-empty + 28) = 105
  set {char}($path_trailing-empty + 29) = 100
  set {char}($path_trailing-empty + 30) = 97
  set {char}($path_trailing-empty + 31) = 116
  set {char}($path_trailing-empty + 32) = 105
  set {char}($path_trailing-empty + 33) = 111
  set {char}($path_trailing-empty + 34) = 110
  set {char}($path_trailing-empty + 35) = 45
  set {char}($path_trailing-empty + 36) = 116
  set {char}($path_trailing-empty + 37) = 114
  set {char}($path_trailing-empty + 38) = 97
  set {char}($path_trailing-empty + 39) = 105
  set {char}($path_trailing-empty + 40) = 108
  set {char}($path_trailing-empty + 41) = 105
  set {char}($path_trailing-empty + 42) = 110
  set {char}($path_trailing-empty + 43) = 103
  set {char}($path_trailing-empty + 44) = 45
  set {char}($path_trailing-empty + 45) = 101
  set {char}($path_trailing-empty + 46) = 109
  set {char}($path_trailing-empty + 47) = 112
  set {char}($path_trailing-empty + 48) = 116
  set {char}($path_trailing-empty + 49) = 121
  set {char}($path_trailing-empty + 50) = 46
  set {char}($path_trailing-empty + 51) = 99
  set {char}($path_trailing-empty + 52) = 115
  set {char}($path_trailing-empty + 53) = 118
  set {char}($path_trailing-empty + 54) = 0
  set $load_trailing-empty = ((short (*)(int, char *))0x10027960)(112, $path_trailing-empty)
  printf "USERDICT_MATRIX case=trailing-empty index=112 load_ax=%d\n", $load_trailing-empty
  set $unload_trailing-empty = ((short (*)(int))0x10027a80)(112)
  printf "USERDICT_MATRIX case=trailing-empty index=112 unload_ax=%d\n", $unload_trailing-empty
  if $load_trailing-empty != 1
    set $recover_trailing-empty = ((short (*)(int, char *))0x10027960)(112, $valid_path)
    printf "USERDICT_MATRIX case=trailing-empty index=112 valid_recovery_ax=%d\n", $recover_trailing-empty
    if $recover_trailing-empty == 1
      set $recover_unload_trailing-empty = ((short (*)(int))0x10027a80)(112)
      printf "USERDICT_MATRIX case=trailing-empty index=112 recovery_unload_ax=%d\n", $recover_unload_trailing-empty
    end
  end
  set $path_crlf-p = ((char *(*)(unsigned int))0x1001d9c0)(47)
  set {char}($path_crlf-p + 0) = 90
  set {char}($path_crlf-p + 1) = 58
  set {char}($path_crlf-p + 2) = 47
  set {char}($path_crlf-p + 3) = 119
  set {char}($path_crlf-p + 4) = 111
  set {char}($path_crlf-p + 5) = 114
  set {char}($path_crlf-p + 6) = 107
  set {char}($path_crlf-p + 7) = 47
  set {char}($path_crlf-p + 8) = 115
  set {char}($path_crlf-p + 9) = 116
  set {char}($path_crlf-p + 10) = 97
  set {char}($path_crlf-p + 11) = 103
  set {char}($path_crlf-p + 12) = 101
  set {char}($path_crlf-p + 13) = 49
  set {char}($path_crlf-p + 14) = 54
  set {char}($path_crlf-p + 15) = 47
  set {char}($path_crlf-p + 16) = 117
  set {char}($path_crlf-p + 17) = 115
  set {char}($path_crlf-p + 18) = 101
  set {char}($path_crlf-p + 19) = 114
  set {char}($path_crlf-p + 20) = 100
  set {char}($path_crlf-p + 21) = 105
  set {char}($path_crlf-p + 22) = 99
  set {char}($path_crlf-p + 23) = 116
  set {char}($path_crlf-p + 24) = 45
  set {char}($path_crlf-p + 25) = 118
  set {char}($path_crlf-p + 26) = 97
  set {char}($path_crlf-p + 27) = 108
  set {char}($path_crlf-p + 28) = 105
  set {char}($path_crlf-p + 29) = 100
  set {char}($path_crlf-p + 30) = 97
  set {char}($path_crlf-p + 31) = 116
  set {char}($path_crlf-p + 32) = 105
  set {char}($path_crlf-p + 33) = 111
  set {char}($path_crlf-p + 34) = 110
  set {char}($path_crlf-p + 35) = 45
  set {char}($path_crlf-p + 36) = 99
  set {char}($path_crlf-p + 37) = 114
  set {char}($path_crlf-p + 38) = 108
  set {char}($path_crlf-p + 39) = 102
  set {char}($path_crlf-p + 40) = 45
  set {char}($path_crlf-p + 41) = 112
  set {char}($path_crlf-p + 42) = 46
  set {char}($path_crlf-p + 43) = 99
  set {char}($path_crlf-p + 44) = 115
  set {char}($path_crlf-p + 45) = 118
  set {char}($path_crlf-p + 46) = 0
  set $load_crlf-p = ((short (*)(int, char *))0x10027960)(113, $path_crlf-p)
  printf "USERDICT_MATRIX case=crlf-p index=113 load_ax=%d\n", $load_crlf-p
  set $unload_crlf-p = ((short (*)(int))0x10027a80)(113)
  printf "USERDICT_MATRIX case=crlf-p index=113 unload_ax=%d\n", $unload_crlf-p
  if $load_crlf-p != 1
    set $recover_crlf-p = ((short (*)(int, char *))0x10027960)(113, $valid_path)
    printf "USERDICT_MATRIX case=crlf-p index=113 valid_recovery_ax=%d\n", $recover_crlf-p
    if $recover_crlf-p == 1
      set $recover_unload_crlf-p = ((short (*)(int))0x10027a80)(113)
      printf "USERDICT_MATRIX case=crlf-p index=113 recovery_unload_ax=%d\n", $recover_unload_crlf-p
    end
  end
  set $path_invalid-target-p = ((char *(*)(unsigned int))0x1001d9c0)(57)
  set {char}($path_invalid-target-p + 0) = 90
  set {char}($path_invalid-target-p + 1) = 58
  set {char}($path_invalid-target-p + 2) = 47
  set {char}($path_invalid-target-p + 3) = 119
  set {char}($path_invalid-target-p + 4) = 111
  set {char}($path_invalid-target-p + 5) = 114
  set {char}($path_invalid-target-p + 6) = 107
  set {char}($path_invalid-target-p + 7) = 47
  set {char}($path_invalid-target-p + 8) = 115
  set {char}($path_invalid-target-p + 9) = 116
  set {char}($path_invalid-target-p + 10) = 97
  set {char}($path_invalid-target-p + 11) = 103
  set {char}($path_invalid-target-p + 12) = 101
  set {char}($path_invalid-target-p + 13) = 49
  set {char}($path_invalid-target-p + 14) = 54
  set {char}($path_invalid-target-p + 15) = 47
  set {char}($path_invalid-target-p + 16) = 117
  set {char}($path_invalid-target-p + 17) = 115
  set {char}($path_invalid-target-p + 18) = 101
  set {char}($path_invalid-target-p + 19) = 114
  set {char}($path_invalid-target-p + 20) = 100
  set {char}($path_invalid-target-p + 21) = 105
  set {char}($path_invalid-target-p + 22) = 99
  set {char}($path_invalid-target-p + 23) = 116
  set {char}($path_invalid-target-p + 24) = 45
  set {char}($path_invalid-target-p + 25) = 118
  set {char}($path_invalid-target-p + 26) = 97
  set {char}($path_invalid-target-p + 27) = 108
  set {char}($path_invalid-target-p + 28) = 105
  set {char}($path_invalid-target-p + 29) = 100
  set {char}($path_invalid-target-p + 30) = 97
  set {char}($path_invalid-target-p + 31) = 116
  set {char}($path_invalid-target-p + 32) = 105
  set {char}($path_invalid-target-p + 33) = 111
  set {char}($path_invalid-target-p + 34) = 110
  set {char}($path_invalid-target-p + 35) = 45
  set {char}($path_invalid-target-p + 36) = 105
  set {char}($path_invalid-target-p + 37) = 110
  set {char}($path_invalid-target-p + 38) = 118
  set {char}($path_invalid-target-p + 39) = 97
  set {char}($path_invalid-target-p + 40) = 108
  set {char}($path_invalid-target-p + 41) = 105
  set {char}($path_invalid-target-p + 42) = 100
  set {char}($path_invalid-target-p + 43) = 45
  set {char}($path_invalid-target-p + 44) = 116
  set {char}($path_invalid-target-p + 45) = 97
  set {char}($path_invalid-target-p + 46) = 114
  set {char}($path_invalid-target-p + 47) = 103
  set {char}($path_invalid-target-p + 48) = 101
  set {char}($path_invalid-target-p + 49) = 116
  set {char}($path_invalid-target-p + 50) = 45
  set {char}($path_invalid-target-p + 51) = 112
  set {char}($path_invalid-target-p + 52) = 46
  set {char}($path_invalid-target-p + 53) = 99
  set {char}($path_invalid-target-p + 54) = 115
  set {char}($path_invalid-target-p + 55) = 118
  set {char}($path_invalid-target-p + 56) = 0
  set $load_invalid-target-p = ((short (*)(int, char *))0x10027960)(114, $path_invalid-target-p)
  printf "USERDICT_MATRIX case=invalid-target-p index=114 load_ax=%d\n", $load_invalid-target-p
  set $unload_invalid-target-p = ((short (*)(int))0x10027a80)(114)
  printf "USERDICT_MATRIX case=invalid-target-p index=114 unload_ax=%d\n", $unload_invalid-target-p
  if $load_invalid-target-p != 1
    set $recover_invalid-target-p = ((short (*)(int, char *))0x10027960)(114, $valid_path)
    printf "USERDICT_MATRIX case=invalid-target-p index=114 valid_recovery_ax=%d\n", $recover_invalid-target-p
    if $recover_invalid-target-p == 1
      set $recover_unload_invalid-target-p = ((short (*)(int))0x10027a80)(114)
      printf "USERDICT_MATRIX case=invalid-target-p index=114 recovery_unload_ax=%d\n", $recover_unload_invalid-target-p
    end
  end
  continue
end

continue
