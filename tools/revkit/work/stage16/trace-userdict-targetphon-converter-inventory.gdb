set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $set_a = ((char *(*)(unsigned int))0x1001d9c0)(104)
  set {char}($set_a + 0) = 66
  set {char}($set_a + 1) = 32
  set {char}($set_a + 2) = 68
  set {char}($set_a + 3) = 32
  set {char}($set_a + 4) = 70
  set {char}($set_a + 5) = 32
  set {char}($set_a + 6) = 71
  set {char}($set_a + 7) = 32
  set {char}($set_a + 8) = 75
  set {char}($set_a + 9) = 32
  set {char}($set_a + 10) = 76
  set {char}($set_a + 11) = 32
  set {char}($set_a + 12) = 77
  set {char}($set_a + 13) = 32
  set {char}($set_a + 14) = 78
  set {char}($set_a + 15) = 32
  set {char}($set_a + 16) = 80
  set {char}($set_a + 17) = 32
  set {char}($set_a + 18) = 82
  set {char}($set_a + 19) = 32
  set {char}($set_a + 20) = 83
  set {char}($set_a + 21) = 32
  set {char}($set_a + 22) = 84
  set {char}($set_a + 23) = 32
  set {char}($set_a + 24) = 86
  set {char}($set_a + 25) = 32
  set {char}($set_a + 26) = 87
  set {char}($set_a + 27) = 32
  set {char}($set_a + 28) = 89
  set {char}($set_a + 29) = 32
  set {char}($set_a + 30) = 90
  set {char}($set_a + 31) = 32
  set {char}($set_a + 32) = 67
  set {char}($set_a + 33) = 72
  set {char}($set_a + 34) = 32
  set {char}($set_a + 35) = 68
  set {char}($set_a + 36) = 72
  set {char}($set_a + 37) = 32
  set {char}($set_a + 38) = 72
  set {char}($set_a + 39) = 72
  set {char}($set_a + 40) = 32
  set {char}($set_a + 41) = 74
  set {char}($set_a + 42) = 72
  set {char}($set_a + 43) = 32
  set {char}($set_a + 44) = 78
  set {char}($set_a + 45) = 71
  set {char}($set_a + 46) = 32
  set {char}($set_a + 47) = 83
  set {char}($set_a + 48) = 72
  set {char}($set_a + 49) = 32
  set {char}($set_a + 50) = 84
  set {char}($set_a + 51) = 72
  set {char}($set_a + 52) = 32
  set {char}($set_a + 53) = 90
  set {char}($set_a + 54) = 72
  set {char}($set_a + 55) = 32
  set {char}($set_a + 56) = 65
  set {char}($set_a + 57) = 65
  set {char}($set_a + 58) = 48
  set {char}($set_a + 59) = 32
  set {char}($set_a + 60) = 65
  set {char}($set_a + 61) = 69
  set {char}($set_a + 62) = 48
  set {char}($set_a + 63) = 32
  set {char}($set_a + 64) = 65
  set {char}($set_a + 65) = 72
  set {char}($set_a + 66) = 48
  set {char}($set_a + 67) = 32
  set {char}($set_a + 68) = 65
  set {char}($set_a + 69) = 79
  set {char}($set_a + 70) = 48
  set {char}($set_a + 71) = 32
  set {char}($set_a + 72) = 65
  set {char}($set_a + 73) = 87
  set {char}($set_a + 74) = 48
  set {char}($set_a + 75) = 32
  set {char}($set_a + 76) = 65
  set {char}($set_a + 77) = 89
  set {char}($set_a + 78) = 48
  set {char}($set_a + 79) = 32
  set {char}($set_a + 80) = 69
  set {char}($set_a + 81) = 72
  set {char}($set_a + 82) = 48
  set {char}($set_a + 83) = 32
  set {char}($set_a + 84) = 69
  set {char}($set_a + 85) = 82
  set {char}($set_a + 86) = 48
  set {char}($set_a + 87) = 32
  set {char}($set_a + 88) = 69
  set {char}($set_a + 89) = 89
  set {char}($set_a + 90) = 48
  set {char}($set_a + 91) = 32
  set {char}($set_a + 92) = 73
  set {char}($set_a + 93) = 72
  set {char}($set_a + 94) = 48
  set {char}($set_a + 95) = 32
  set {char}($set_a + 96) = 73
  set {char}($set_a + 97) = 89
  set {char}($set_a + 98) = 48
  set {char}($set_a + 99) = 32
  set {char}($set_a + 100) = 79
  set {char}($set_a + 101) = 87
  set {char}($set_a + 102) = 48
  set {char}($set_a + 103) = 0
  set $set_b = ((char *(*)(unsigned int))0x1001d9c0)(134)
  set {char}($set_b + 0) = 65
  set {char}($set_b + 1) = 65
  set {char}($set_b + 2) = 49
  set {char}($set_b + 3) = 32
  set {char}($set_b + 4) = 65
  set {char}($set_b + 5) = 69
  set {char}($set_b + 6) = 49
  set {char}($set_b + 7) = 32
  set {char}($set_b + 8) = 65
  set {char}($set_b + 9) = 72
  set {char}($set_b + 10) = 49
  set {char}($set_b + 11) = 32
  set {char}($set_b + 12) = 65
  set {char}($set_b + 13) = 79
  set {char}($set_b + 14) = 49
  set {char}($set_b + 15) = 32
  set {char}($set_b + 16) = 65
  set {char}($set_b + 17) = 87
  set {char}($set_b + 18) = 49
  set {char}($set_b + 19) = 32
  set {char}($set_b + 20) = 65
  set {char}($set_b + 21) = 89
  set {char}($set_b + 22) = 49
  set {char}($set_b + 23) = 32
  set {char}($set_b + 24) = 69
  set {char}($set_b + 25) = 72
  set {char}($set_b + 26) = 49
  set {char}($set_b + 27) = 32
  set {char}($set_b + 28) = 69
  set {char}($set_b + 29) = 82
  set {char}($set_b + 30) = 49
  set {char}($set_b + 31) = 32
  set {char}($set_b + 32) = 69
  set {char}($set_b + 33) = 89
  set {char}($set_b + 34) = 49
  set {char}($set_b + 35) = 32
  set {char}($set_b + 36) = 73
  set {char}($set_b + 37) = 72
  set {char}($set_b + 38) = 49
  set {char}($set_b + 39) = 32
  set {char}($set_b + 40) = 73
  set {char}($set_b + 41) = 89
  set {char}($set_b + 42) = 49
  set {char}($set_b + 43) = 32
  set {char}($set_b + 44) = 79
  set {char}($set_b + 45) = 87
  set {char}($set_b + 46) = 49
  set {char}($set_b + 47) = 32
  set {char}($set_b + 48) = 79
  set {char}($set_b + 49) = 89
  set {char}($set_b + 50) = 48
  set {char}($set_b + 51) = 32
  set {char}($set_b + 52) = 85
  set {char}($set_b + 53) = 72
  set {char}($set_b + 54) = 48
  set {char}($set_b + 55) = 32
  set {char}($set_b + 56) = 85
  set {char}($set_b + 57) = 87
  set {char}($set_b + 58) = 48
  set {char}($set_b + 59) = 32
  set {char}($set_b + 60) = 65
  set {char}($set_b + 61) = 65
  set {char}($set_b + 62) = 50
  set {char}($set_b + 63) = 32
  set {char}($set_b + 64) = 65
  set {char}($set_b + 65) = 69
  set {char}($set_b + 66) = 50
  set {char}($set_b + 67) = 32
  set {char}($set_b + 68) = 65
  set {char}($set_b + 69) = 72
  set {char}($set_b + 70) = 50
  set {char}($set_b + 71) = 32
  set {char}($set_b + 72) = 65
  set {char}($set_b + 73) = 79
  set {char}($set_b + 74) = 50
  set {char}($set_b + 75) = 32
  set {char}($set_b + 76) = 65
  set {char}($set_b + 77) = 87
  set {char}($set_b + 78) = 50
  set {char}($set_b + 79) = 32
  set {char}($set_b + 80) = 65
  set {char}($set_b + 81) = 89
  set {char}($set_b + 82) = 50
  set {char}($set_b + 83) = 32
  set {char}($set_b + 84) = 69
  set {char}($set_b + 85) = 72
  set {char}($set_b + 86) = 50
  set {char}($set_b + 87) = 32
  set {char}($set_b + 88) = 69
  set {char}($set_b + 89) = 82
  set {char}($set_b + 90) = 50
  set {char}($set_b + 91) = 32
  set {char}($set_b + 92) = 69
  set {char}($set_b + 93) = 89
  set {char}($set_b + 94) = 50
  set {char}($set_b + 95) = 32
  set {char}($set_b + 96) = 73
  set {char}($set_b + 97) = 72
  set {char}($set_b + 98) = 50
  set {char}($set_b + 99) = 32
  set {char}($set_b + 100) = 73
  set {char}($set_b + 101) = 89
  set {char}($set_b + 102) = 50
  set {char}($set_b + 103) = 32
  set {char}($set_b + 104) = 79
  set {char}($set_b + 105) = 87
  set {char}($set_b + 106) = 50
  set {char}($set_b + 107) = 32
  set {char}($set_b + 108) = 79
  set {char}($set_b + 109) = 89
  set {char}($set_b + 110) = 49
  set {char}($set_b + 111) = 32
  set {char}($set_b + 112) = 85
  set {char}($set_b + 113) = 72
  set {char}($set_b + 114) = 49
  set {char}($set_b + 115) = 32
  set {char}($set_b + 116) = 85
  set {char}($set_b + 117) = 87
  set {char}($set_b + 118) = 49
  set {char}($set_b + 119) = 32
  set {char}($set_b + 120) = 79
  set {char}($set_b + 121) = 89
  set {char}($set_b + 122) = 50
  set {char}($set_b + 123) = 32
  set {char}($set_b + 124) = 85
  set {char}($set_b + 125) = 72
  set {char}($set_b + 126) = 50
  set {char}($set_b + 127) = 32
  set {char}($set_b + 128) = 85
  set {char}($set_b + 129) = 87
  set {char}($set_b + 130) = 50
  set {char}($set_b + 131) = 32
  set {char}($set_b + 132) = 35
  set {char}($set_b + 133) = 0
  set $output = ((char *(*)(unsigned int))0x1001d9c0)(66)
  set $i = 0
  while $i < 66
    set {char}($output + $i) = 0
    set $i = $i + 1
  end
  set $result_a = ((short (*)(char *, char *))0x1005f710)($output, $set_a)
  printf "TARGETPHON converter=set-a result=%d bytes=", $result_a
  set $i = 0
  while *(unsigned char *)($output + $i) != 0
    printf "%02x", *(unsigned char *)($output + $i)
    set $i = $i + 1
    if *(unsigned char *)($output + $i) != 0
      printf ","
    end
  end
  printf "\n"
  set $result_b = ((short (*)(char *, char *))0x1005f710)($output, $set_b)
  printf "TARGETPHON converter=set-b result=%d bytes=", $result_b
  set $i = 0
  while *(unsigned char *)($output + $i) != 0
    printf "%02x", *(unsigned char *)($output + $i)
    set $i = $i + 1
    if *(unsigned char *)($output + $i) != 0
      printf ","
    end
  end
  printf "\n"
  call ((void (*)(void *))0x1001da30)($set_a)
  call ((void (*)(void *))0x1001da30)($set_b)
  call ((void (*)(void *))0x1001da30)($output)
  continue
end
continue
