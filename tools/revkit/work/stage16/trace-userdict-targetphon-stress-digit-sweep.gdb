set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $phon = ((char *(*)(unsigned int))0x1001d9c0)(5)
  set {unsigned char}($phon) = 65
  set {unsigned char}($phon + 1) = 65
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=AA%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 65
  set {unsigned char}($phon + 1) = 69
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=AE%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 65
  set {unsigned char}($phon + 1) = 72
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=AH%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 65
  set {unsigned char}($phon + 1) = 79
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=AO%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 65
  set {unsigned char}($phon + 1) = 87
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=AW%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 65
  set {unsigned char}($phon + 1) = 89
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=AY%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 69
  set {unsigned char}($phon + 1) = 72
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=EH%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 69
  set {unsigned char}($phon + 1) = 82
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=ER%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 69
  set {unsigned char}($phon + 1) = 89
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=EY%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 73
  set {unsigned char}($phon + 1) = 72
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=IH%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 73
  set {unsigned char}($phon + 1) = 89
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=IY%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 79
  set {unsigned char}($phon + 1) = 87
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=OW%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 79
  set {unsigned char}($phon + 1) = 89
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=OY%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 85
  set {unsigned char}($phon + 1) = 72
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=UH%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  set {unsigned char}($phon) = 85
  set {unsigned char}($phon + 1) = 87
  set $stress = 48
  while $stress <= 57
    set {unsigned char}($phon + 2) = $stress
    set {char}($phon + 3) = 0
    set $result = ((short (*)(char *))0x1002a590)($phon)
    printf "TARGETPHON digit=UW%c result=%d\n", $stress, $result
    set $stress = $stress + 1
  end
  call ((void (*)(void *))0x1001da30)($phon)
  continue
end
continue
