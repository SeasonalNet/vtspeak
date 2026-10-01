set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 66
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=B input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 98
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=B input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 68
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=D input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 100
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=D input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 70
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=F input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 102
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=F input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 71
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=G input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 103
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=G input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 75
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=K input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 107
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=K input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 76
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=L input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 108
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=L input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 77
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=M input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 109
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=M input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 78
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=N input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 110
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=N input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 80
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=P input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 112
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=P input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 82
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=R input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 114
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=R input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 83
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=S input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 115
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=S input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 84
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=T input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 116
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=T input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 86
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=V input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 118
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=V input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 87
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=W input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 119
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=W input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 89
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=Y input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 121
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=Y input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 90
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=Z input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($p + 0) = 122
  set {char}($p + 1) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=single base=Z input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 67
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=CH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 99
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=CH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 67
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=CH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 99
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=CH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 68
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=DH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 100
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=DH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 68
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=DH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 100
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=DH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 72
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=HH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 104
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=HH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 72
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=HH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 104
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=HH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 74
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=JH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 106
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=JH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 74
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=JH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 106
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=JH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 78
  set {char}($p + 1) = 71
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=NG input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 110
  set {char}($p + 1) = 71
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=NG input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 78
  set {char}($p + 1) = 103
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=NG input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 110
  set {char}($p + 1) = 103
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=NG input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 83
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=SH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 115
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=SH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 83
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=SH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 115
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=SH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 84
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=TH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 116
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=TH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 84
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=TH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 116
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=TH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 90
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=ZH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 122
  set {char}($p + 1) = 72
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=ZH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 90
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=ZH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($p + 0) = 122
  set {char}($p + 1) = 104
  set {char}($p + 2) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=pair base=ZH input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 65
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 65
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 97
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 97
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 65
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 65
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 97
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 97
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 65
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 65
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 97
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 97
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AA2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 69
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 69
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 101
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 101
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 69
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 69
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 101
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 101
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 69
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 69
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 101
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 101
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AE2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 79
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 79
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 111
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 111
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 79
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 79
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 111
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 111
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 79
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 79
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 111
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 111
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AO2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 87
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 87
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 119
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 119
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 87
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 87
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 119
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 119
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 87
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 87
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 119
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 119
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 65
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 97
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=AY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 82
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 82
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 114
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 114
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 82
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 82
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 114
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 114
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 82
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 82
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 114
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 114
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=ER2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 69
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 101
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=EY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 73
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 105
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=IY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 87
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 87
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 119
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 119
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 87
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 87
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 119
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 119
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 87
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 87
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 119
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 119
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 89
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 121
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 89
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 121
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 89
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 79
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 111
  set {char}($p + 1) = 121
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=OY2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 72
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 104
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 72
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 104
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 72
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 104
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UH2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 87
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 87
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 119
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 119
  set {char}($p + 2) = 48
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW0 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 87
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 87
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 119
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 119
  set {char}($p + 2) = 49
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW1 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 87
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 87
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 85
  set {char}($p + 1) = 119
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW2 input=%s ax=%u\n", $p, $r & 65535
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($p + 0) = 117
  set {char}($p + 1) = 119
  set {char}($p + 2) = 50
  set {char}($p + 3) = 0
  set $r = ((int (*)(char *))0x1002a590)($p)
  printf "TARGET_PHON_CASE family=vowel base=UW2 input=%s ax=%u\n", $p, $r & 65535
  continue
end

continue

