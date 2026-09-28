set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $first12 = ((char *(*)(unsigned int))0x1001d9c0)(24)
  set {char}($first12 + 0) = 66
  set {char}($first12 + 1) = 32
  set {char}($first12 + 2) = 68
  set {char}($first12 + 3) = 32
  set {char}($first12 + 4) = 70
  set {char}($first12 + 5) = 32
  set {char}($first12 + 6) = 71
  set {char}($first12 + 7) = 32
  set {char}($first12 + 8) = 75
  set {char}($first12 + 9) = 32
  set {char}($first12 + 10) = 76
  set {char}($first12 + 11) = 32
  set {char}($first12 + 12) = 77
  set {char}($first12 + 13) = 32
  set {char}($first12 + 14) = 78
  set {char}($first12 + 15) = 32
  set {char}($first12 + 16) = 80
  set {char}($first12 + 17) = 32
  set {char}($first12 + 18) = 82
  set {char}($first12 + 19) = 32
  set {char}($first12 + 20) = 83
  set {char}($first12 + 21) = 32
  set {char}($first12 + 22) = 84
  set {char}($first12 + 23) = 0
  set $rest24 = ((char *(*)(unsigned int))0x1001d9c0)(80)
  set {char}($rest24 + 0) = 86
  set {char}($rest24 + 1) = 32
  set {char}($rest24 + 2) = 87
  set {char}($rest24 + 3) = 32
  set {char}($rest24 + 4) = 89
  set {char}($rest24 + 5) = 32
  set {char}($rest24 + 6) = 90
  set {char}($rest24 + 7) = 32
  set {char}($rest24 + 8) = 67
  set {char}($rest24 + 9) = 72
  set {char}($rest24 + 10) = 32
  set {char}($rest24 + 11) = 68
  set {char}($rest24 + 12) = 72
  set {char}($rest24 + 13) = 32
  set {char}($rest24 + 14) = 72
  set {char}($rest24 + 15) = 72
  set {char}($rest24 + 16) = 32
  set {char}($rest24 + 17) = 74
  set {char}($rest24 + 18) = 72
  set {char}($rest24 + 19) = 32
  set {char}($rest24 + 20) = 78
  set {char}($rest24 + 21) = 71
  set {char}($rest24 + 22) = 32
  set {char}($rest24 + 23) = 83
  set {char}($rest24 + 24) = 72
  set {char}($rest24 + 25) = 32
  set {char}($rest24 + 26) = 84
  set {char}($rest24 + 27) = 72
  set {char}($rest24 + 28) = 32
  set {char}($rest24 + 29) = 90
  set {char}($rest24 + 30) = 72
  set {char}($rest24 + 31) = 32
  set {char}($rest24 + 32) = 65
  set {char}($rest24 + 33) = 65
  set {char}($rest24 + 34) = 48
  set {char}($rest24 + 35) = 32
  set {char}($rest24 + 36) = 65
  set {char}($rest24 + 37) = 69
  set {char}($rest24 + 38) = 48
  set {char}($rest24 + 39) = 32
  set {char}($rest24 + 40) = 65
  set {char}($rest24 + 41) = 72
  set {char}($rest24 + 42) = 48
  set {char}($rest24 + 43) = 32
  set {char}($rest24 + 44) = 65
  set {char}($rest24 + 45) = 79
  set {char}($rest24 + 46) = 48
  set {char}($rest24 + 47) = 32
  set {char}($rest24 + 48) = 65
  set {char}($rest24 + 49) = 87
  set {char}($rest24 + 50) = 48
  set {char}($rest24 + 51) = 32
  set {char}($rest24 + 52) = 65
  set {char}($rest24 + 53) = 89
  set {char}($rest24 + 54) = 48
  set {char}($rest24 + 55) = 32
  set {char}($rest24 + 56) = 69
  set {char}($rest24 + 57) = 72
  set {char}($rest24 + 58) = 48
  set {char}($rest24 + 59) = 32
  set {char}($rest24 + 60) = 69
  set {char}($rest24 + 61) = 82
  set {char}($rest24 + 62) = 48
  set {char}($rest24 + 63) = 32
  set {char}($rest24 + 64) = 69
  set {char}($rest24 + 65) = 89
  set {char}($rest24 + 66) = 48
  set {char}($rest24 + 67) = 32
  set {char}($rest24 + 68) = 73
  set {char}($rest24 + 69) = 72
  set {char}($rest24 + 70) = 48
  set {char}($rest24 + 71) = 32
  set {char}($rest24 + 72) = 73
  set {char}($rest24 + 73) = 89
  set {char}($rest24 + 74) = 48
  set {char}($rest24 + 75) = 32
  set {char}($rest24 + 76) = 79
  set {char}($rest24 + 77) = 87
  set {char}($rest24 + 78) = 48
  set {char}($rest24 + 79) = 0
  set $all36 = ((char *(*)(unsigned int))0x1001d9c0)(104)
  set {char}($all36 + 0) = 66
  set {char}($all36 + 1) = 32
  set {char}($all36 + 2) = 68
  set {char}($all36 + 3) = 32
  set {char}($all36 + 4) = 70
  set {char}($all36 + 5) = 32
  set {char}($all36 + 6) = 71
  set {char}($all36 + 7) = 32
  set {char}($all36 + 8) = 75
  set {char}($all36 + 9) = 32
  set {char}($all36 + 10) = 76
  set {char}($all36 + 11) = 32
  set {char}($all36 + 12) = 77
  set {char}($all36 + 13) = 32
  set {char}($all36 + 14) = 78
  set {char}($all36 + 15) = 32
  set {char}($all36 + 16) = 80
  set {char}($all36 + 17) = 32
  set {char}($all36 + 18) = 82
  set {char}($all36 + 19) = 32
  set {char}($all36 + 20) = 83
  set {char}($all36 + 21) = 32
  set {char}($all36 + 22) = 84
  set {char}($all36 + 23) = 32
  set {char}($all36 + 24) = 86
  set {char}($all36 + 25) = 32
  set {char}($all36 + 26) = 87
  set {char}($all36 + 27) = 32
  set {char}($all36 + 28) = 89
  set {char}($all36 + 29) = 32
  set {char}($all36 + 30) = 90
  set {char}($all36 + 31) = 32
  set {char}($all36 + 32) = 67
  set {char}($all36 + 33) = 72
  set {char}($all36 + 34) = 32
  set {char}($all36 + 35) = 68
  set {char}($all36 + 36) = 72
  set {char}($all36 + 37) = 32
  set {char}($all36 + 38) = 72
  set {char}($all36 + 39) = 72
  set {char}($all36 + 40) = 32
  set {char}($all36 + 41) = 74
  set {char}($all36 + 42) = 72
  set {char}($all36 + 43) = 32
  set {char}($all36 + 44) = 78
  set {char}($all36 + 45) = 71
  set {char}($all36 + 46) = 32
  set {char}($all36 + 47) = 83
  set {char}($all36 + 48) = 72
  set {char}($all36 + 49) = 32
  set {char}($all36 + 50) = 84
  set {char}($all36 + 51) = 72
  set {char}($all36 + 52) = 32
  set {char}($all36 + 53) = 90
  set {char}($all36 + 54) = 72
  set {char}($all36 + 55) = 32
  set {char}($all36 + 56) = 65
  set {char}($all36 + 57) = 65
  set {char}($all36 + 58) = 48
  set {char}($all36 + 59) = 32
  set {char}($all36 + 60) = 65
  set {char}($all36 + 61) = 69
  set {char}($all36 + 62) = 48
  set {char}($all36 + 63) = 32
  set {char}($all36 + 64) = 65
  set {char}($all36 + 65) = 72
  set {char}($all36 + 66) = 48
  set {char}($all36 + 67) = 32
  set {char}($all36 + 68) = 65
  set {char}($all36 + 69) = 79
  set {char}($all36 + 70) = 48
  set {char}($all36 + 71) = 32
  set {char}($all36 + 72) = 65
  set {char}($all36 + 73) = 87
  set {char}($all36 + 74) = 48
  set {char}($all36 + 75) = 32
  set {char}($all36 + 76) = 65
  set {char}($all36 + 77) = 89
  set {char}($all36 + 78) = 48
  set {char}($all36 + 79) = 32
  set {char}($all36 + 80) = 69
  set {char}($all36 + 81) = 72
  set {char}($all36 + 82) = 48
  set {char}($all36 + 83) = 32
  set {char}($all36 + 84) = 69
  set {char}($all36 + 85) = 82
  set {char}($all36 + 86) = 48
  set {char}($all36 + 87) = 32
  set {char}($all36 + 88) = 69
  set {char}($all36 + 89) = 89
  set {char}($all36 + 90) = 48
  set {char}($all36 + 91) = 32
  set {char}($all36 + 92) = 73
  set {char}($all36 + 93) = 72
  set {char}($all36 + 94) = 48
  set {char}($all36 + 95) = 32
  set {char}($all36 + 96) = 73
  set {char}($all36 + 97) = 89
  set {char}($all36 + 98) = 48
  set {char}($all36 + 99) = 32
  set {char}($all36 + 100) = 79
  set {char}($all36 + 101) = 87
  set {char}($all36 + 102) = 48
  set {char}($all36 + 103) = 0
  set $output = ((char *(*)(unsigned int))0x1001d9c0)(66)
  set $i = 0
  while $i < 66
    set {char}($output + $i) = 0
    set $i = $i + 1
  end
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $first12)
  printf "TARGETPHON sequence=first12 result=%d bytes=", $result
  set $i = 0
  while $i < 12
    printf "%02x", *(unsigned char *)($output + $i)
    set $i = $i + 1
  end
  printf "\n"
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $rest24)
  printf "TARGETPHON sequence=rest24 result=%d bytes=", $result
  set $i = 0
  while $i < 24
    printf "%02x", *(unsigned char *)($output + $i)
    set $i = $i + 1
  end
  printf "\n"
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $all36)
  printf "TARGETPHON sequence=all36 result=%d bytes=", $result
  set $i = 0
  while $i < 36
    printf "%02x", *(unsigned char *)($output + $i)
    set $i = $i + 1
  end
  printf "\n"
  call ((void (*)(void *))0x1001da30)($first12)
  call ((void (*)(void *))0x1001da30)($rest24)
  call ((void (*)(void *))0x1001da30)($all36)
  call ((void (*)(void *))0x1001da30)($output)
  continue
end
continue
