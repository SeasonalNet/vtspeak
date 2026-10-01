set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(char **)($esp + 8)
  set $speaker = *(int *)($esp + 16)
  disable 1
  set $probe_path = (char *)malloc(512)
  set $probe_path = (char *)malloc(47)
  set {char}($probe_path + 0) = 90
  set {char}($probe_path + 1) = 58
  set {char}($probe_path + 2) = 47
  set {char}($probe_path + 3) = 119
  set {char}($probe_path + 4) = 111
  set {char}($probe_path + 5) = 114
  set {char}($probe_path + 6) = 107
  set {char}($probe_path + 7) = 47
  set {char}($probe_path + 8) = 115
  set {char}($probe_path + 9) = 116
  set {char}($probe_path + 10) = 97
  set {char}($probe_path + 11) = 103
  set {char}($probe_path + 12) = 101
  set {char}($probe_path + 13) = 50
  set {char}($probe_path + 14) = 49
  set {char}($probe_path + 15) = 47
  set {char}($probe_path + 16) = 116
  set {char}($probe_path + 17) = 101
  set {char}($probe_path + 18) = 120
  set {char}($probe_path + 19) = 116
  set {char}($probe_path + 20) = 45
  set {char}($probe_path + 21) = 102
  set {char}($probe_path + 22) = 105
  set {char}($probe_path + 23) = 108
  set {char}($probe_path + 24) = 101
  set {char}($probe_path + 25) = 45
  set {char}($probe_path + 26) = 112
  set {char}($probe_path + 27) = 97
  set {char}($probe_path + 28) = 116
  set {char}($probe_path + 29) = 104
  set {char}($probe_path + 30) = 45
  set {char}($probe_path + 31) = 112
  set {char}($probe_path + 32) = 114
  set {char}($probe_path + 33) = 111
  set {char}($probe_path + 34) = 98
  set {char}($probe_path + 35) = 101
  set {char}($probe_path + 36) = 45
  set {char}($probe_path + 37) = 118
  set {char}($probe_path + 38) = 97
  set {char}($probe_path + 39) = 108
  set {char}($probe_path + 40) = 105
  set {char}($probe_path + 41) = 100
  set {char}($probe_path + 42) = 46
  set {char}($probe_path + 43) = 119
  set {char}($probe_path + 44) = 97
  set {char}($probe_path + 45) = 118
  set {char}($probe_path + 46) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $probe_path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "TEXT_FILE_PATH name=valid result=%d\n", $result
  set $probe_path = (char *)malloc(31)
  set {char}($probe_path + 0) = 90
  set {char}($probe_path + 1) = 58
  set {char}($probe_path + 2) = 47
  set {char}($probe_path + 3) = 119
  set {char}($probe_path + 4) = 111
  set {char}($probe_path + 5) = 114
  set {char}($probe_path + 6) = 107
  set {char}($probe_path + 7) = 47
  set {char}($probe_path + 8) = 115
  set {char}($probe_path + 9) = 116
  set {char}($probe_path + 10) = 97
  set {char}($probe_path + 11) = 103
  set {char}($probe_path + 12) = 101
  set {char}($probe_path + 13) = 50
  set {char}($probe_path + 14) = 49
  set {char}($probe_path + 15) = 47
  set {char}($probe_path + 16) = 115
  set {char}($probe_path + 17) = 97
  set {char}($probe_path + 18) = 110
  set {char}($probe_path + 19) = 100
  set {char}($probe_path + 20) = 98
  set {char}($probe_path + 21) = 111
  set {char}($probe_path + 22) = 120
  set {char}($probe_path + 23) = 47
  set {char}($probe_path + 24) = 115
  set {char}($probe_path + 25) = 116
  set {char}($probe_path + 26) = 97
  set {char}($probe_path + 27) = 103
  set {char}($probe_path + 28) = 101
  set {char}($probe_path + 29) = 53
  set {char}($probe_path + 30) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $probe_path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "TEXT_FILE_PATH name=directory result=%d\n", $result
  set $probe_path = (char *)malloc(60)
  set {char}($probe_path + 0) = 90
  set {char}($probe_path + 1) = 58
  set {char}($probe_path + 2) = 47
  set {char}($probe_path + 3) = 119
  set {char}($probe_path + 4) = 111
  set {char}($probe_path + 5) = 114
  set {char}($probe_path + 6) = 107
  set {char}($probe_path + 7) = 47
  set {char}($probe_path + 8) = 115
  set {char}($probe_path + 9) = 116
  set {char}($probe_path + 10) = 97
  set {char}($probe_path + 11) = 103
  set {char}($probe_path + 12) = 101
  set {char}($probe_path + 13) = 50
  set {char}($probe_path + 14) = 49
  set {char}($probe_path + 15) = 47
  set {char}($probe_path + 16) = 116
  set {char}($probe_path + 17) = 101
  set {char}($probe_path + 18) = 120
  set {char}($probe_path + 19) = 116
  set {char}($probe_path + 20) = 45
  set {char}($probe_path + 21) = 102
  set {char}($probe_path + 22) = 105
  set {char}($probe_path + 23) = 108
  set {char}($probe_path + 24) = 101
  set {char}($probe_path + 25) = 45
  set {char}($probe_path + 26) = 112
  set {char}($probe_path + 27) = 97
  set {char}($probe_path + 28) = 116
  set {char}($probe_path + 29) = 104
  set {char}($probe_path + 30) = 45
  set {char}($probe_path + 31) = 109
  set {char}($probe_path + 32) = 105
  set {char}($probe_path + 33) = 115
  set {char}($probe_path + 34) = 115
  set {char}($probe_path + 35) = 105
  set {char}($probe_path + 36) = 110
  set {char}($probe_path + 37) = 103
  set {char}($probe_path + 38) = 45
  set {char}($probe_path + 39) = 112
  set {char}($probe_path + 40) = 97
  set {char}($probe_path + 41) = 114
  set {char}($probe_path + 42) = 101
  set {char}($probe_path + 43) = 110
  set {char}($probe_path + 44) = 116
  set {char}($probe_path + 45) = 45
  set {char}($probe_path + 46) = 118
  set {char}($probe_path + 47) = 49
  set {char}($probe_path + 48) = 47
  set {char}($probe_path + 49) = 111
  set {char}($probe_path + 50) = 117
  set {char}($probe_path + 51) = 116
  set {char}($probe_path + 52) = 112
  set {char}($probe_path + 53) = 117
  set {char}($probe_path + 54) = 116
  set {char}($probe_path + 55) = 46
  set {char}($probe_path + 56) = 119
  set {char}($probe_path + 57) = 97
  set {char}($probe_path + 58) = 118
  set {char}($probe_path + 59) = 0
  set $result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $text, $probe_path, $speaker, -1, -1, -1, -1, -1, -1)
  printf "TEXT_FILE_PATH name=missing_parent result=%d\n", $result
  kill
  quit
end

continue
