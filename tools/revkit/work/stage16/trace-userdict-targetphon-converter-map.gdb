set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands
  silent
  disable 1
  set $phone0 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone0 + 0) = 66
  set {char}($phone0 + 1) = 0
  set $phone1 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone1 + 0) = 68
  set {char}($phone1 + 1) = 0
  set $phone2 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone2 + 0) = 70
  set {char}($phone2 + 1) = 0
  set $phone3 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone3 + 0) = 71
  set {char}($phone3 + 1) = 0
  set $phone4 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone4 + 0) = 75
  set {char}($phone4 + 1) = 0
  set $phone5 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone5 + 0) = 76
  set {char}($phone5 + 1) = 0
  set $phone6 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone6 + 0) = 77
  set {char}($phone6 + 1) = 0
  set $phone7 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone7 + 0) = 78
  set {char}($phone7 + 1) = 0
  set $phone8 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone8 + 0) = 80
  set {char}($phone8 + 1) = 0
  set $phone9 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone9 + 0) = 82
  set {char}($phone9 + 1) = 0
  set $phone10 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone10 + 0) = 83
  set {char}($phone10 + 1) = 0
  set $phone11 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone11 + 0) = 84
  set {char}($phone11 + 1) = 0
  set $phone12 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone12 + 0) = 86
  set {char}($phone12 + 1) = 0
  set $phone13 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone13 + 0) = 87
  set {char}($phone13 + 1) = 0
  set $phone14 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone14 + 0) = 89
  set {char}($phone14 + 1) = 0
  set $phone15 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone15 + 0) = 90
  set {char}($phone15 + 1) = 0
  set $phone16 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone16 + 0) = 67
  set {char}($phone16 + 1) = 72
  set {char}($phone16 + 2) = 0
  set $phone17 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone17 + 0) = 68
  set {char}($phone17 + 1) = 72
  set {char}($phone17 + 2) = 0
  set $phone18 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone18 + 0) = 72
  set {char}($phone18 + 1) = 72
  set {char}($phone18 + 2) = 0
  set $phone19 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone19 + 0) = 74
  set {char}($phone19 + 1) = 72
  set {char}($phone19 + 2) = 0
  set $phone20 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone20 + 0) = 78
  set {char}($phone20 + 1) = 71
  set {char}($phone20 + 2) = 0
  set $phone21 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone21 + 0) = 83
  set {char}($phone21 + 1) = 72
  set {char}($phone21 + 2) = 0
  set $phone22 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone22 + 0) = 84
  set {char}($phone22 + 1) = 72
  set {char}($phone22 + 2) = 0
  set $phone23 = ((char *(*)(unsigned int))0x1001d9c0)(3)
  set {char}($phone23 + 0) = 90
  set {char}($phone23 + 1) = 72
  set {char}($phone23 + 2) = 0
  set $phone24 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone24 + 0) = 65
  set {char}($phone24 + 1) = 65
  set {char}($phone24 + 2) = 48
  set {char}($phone24 + 3) = 0
  set $phone25 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone25 + 0) = 65
  set {char}($phone25 + 1) = 65
  set {char}($phone25 + 2) = 49
  set {char}($phone25 + 3) = 0
  set $phone26 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone26 + 0) = 65
  set {char}($phone26 + 1) = 65
  set {char}($phone26 + 2) = 50
  set {char}($phone26 + 3) = 0
  set $phone27 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone27 + 0) = 65
  set {char}($phone27 + 1) = 69
  set {char}($phone27 + 2) = 48
  set {char}($phone27 + 3) = 0
  set $phone28 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone28 + 0) = 65
  set {char}($phone28 + 1) = 69
  set {char}($phone28 + 2) = 49
  set {char}($phone28 + 3) = 0
  set $phone29 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone29 + 0) = 65
  set {char}($phone29 + 1) = 69
  set {char}($phone29 + 2) = 50
  set {char}($phone29 + 3) = 0
  set $phone30 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone30 + 0) = 65
  set {char}($phone30 + 1) = 72
  set {char}($phone30 + 2) = 48
  set {char}($phone30 + 3) = 0
  set $phone31 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone31 + 0) = 65
  set {char}($phone31 + 1) = 72
  set {char}($phone31 + 2) = 49
  set {char}($phone31 + 3) = 0
  set $phone32 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone32 + 0) = 65
  set {char}($phone32 + 1) = 72
  set {char}($phone32 + 2) = 50
  set {char}($phone32 + 3) = 0
  set $phone33 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone33 + 0) = 65
  set {char}($phone33 + 1) = 79
  set {char}($phone33 + 2) = 48
  set {char}($phone33 + 3) = 0
  set $phone34 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone34 + 0) = 65
  set {char}($phone34 + 1) = 79
  set {char}($phone34 + 2) = 49
  set {char}($phone34 + 3) = 0
  set $phone35 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone35 + 0) = 65
  set {char}($phone35 + 1) = 79
  set {char}($phone35 + 2) = 50
  set {char}($phone35 + 3) = 0
  set $phone36 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone36 + 0) = 65
  set {char}($phone36 + 1) = 87
  set {char}($phone36 + 2) = 48
  set {char}($phone36 + 3) = 0
  set $phone37 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone37 + 0) = 65
  set {char}($phone37 + 1) = 87
  set {char}($phone37 + 2) = 49
  set {char}($phone37 + 3) = 0
  set $phone38 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone38 + 0) = 65
  set {char}($phone38 + 1) = 87
  set {char}($phone38 + 2) = 50
  set {char}($phone38 + 3) = 0
  set $phone39 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone39 + 0) = 65
  set {char}($phone39 + 1) = 89
  set {char}($phone39 + 2) = 48
  set {char}($phone39 + 3) = 0
  set $phone40 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone40 + 0) = 65
  set {char}($phone40 + 1) = 89
  set {char}($phone40 + 2) = 49
  set {char}($phone40 + 3) = 0
  set $phone41 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone41 + 0) = 65
  set {char}($phone41 + 1) = 89
  set {char}($phone41 + 2) = 50
  set {char}($phone41 + 3) = 0
  set $phone42 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone42 + 0) = 69
  set {char}($phone42 + 1) = 72
  set {char}($phone42 + 2) = 48
  set {char}($phone42 + 3) = 0
  set $phone43 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone43 + 0) = 69
  set {char}($phone43 + 1) = 72
  set {char}($phone43 + 2) = 49
  set {char}($phone43 + 3) = 0
  set $phone44 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone44 + 0) = 69
  set {char}($phone44 + 1) = 72
  set {char}($phone44 + 2) = 50
  set {char}($phone44 + 3) = 0
  set $phone45 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone45 + 0) = 69
  set {char}($phone45 + 1) = 82
  set {char}($phone45 + 2) = 48
  set {char}($phone45 + 3) = 0
  set $phone46 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone46 + 0) = 69
  set {char}($phone46 + 1) = 82
  set {char}($phone46 + 2) = 49
  set {char}($phone46 + 3) = 0
  set $phone47 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone47 + 0) = 69
  set {char}($phone47 + 1) = 82
  set {char}($phone47 + 2) = 50
  set {char}($phone47 + 3) = 0
  set $phone48 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone48 + 0) = 69
  set {char}($phone48 + 1) = 89
  set {char}($phone48 + 2) = 48
  set {char}($phone48 + 3) = 0
  set $phone49 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone49 + 0) = 69
  set {char}($phone49 + 1) = 89
  set {char}($phone49 + 2) = 49
  set {char}($phone49 + 3) = 0
  set $phone50 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone50 + 0) = 69
  set {char}($phone50 + 1) = 89
  set {char}($phone50 + 2) = 50
  set {char}($phone50 + 3) = 0
  set $phone51 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone51 + 0) = 73
  set {char}($phone51 + 1) = 72
  set {char}($phone51 + 2) = 48
  set {char}($phone51 + 3) = 0
  set $phone52 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone52 + 0) = 73
  set {char}($phone52 + 1) = 72
  set {char}($phone52 + 2) = 49
  set {char}($phone52 + 3) = 0
  set $phone53 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone53 + 0) = 73
  set {char}($phone53 + 1) = 72
  set {char}($phone53 + 2) = 50
  set {char}($phone53 + 3) = 0
  set $phone54 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone54 + 0) = 73
  set {char}($phone54 + 1) = 89
  set {char}($phone54 + 2) = 48
  set {char}($phone54 + 3) = 0
  set $phone55 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone55 + 0) = 73
  set {char}($phone55 + 1) = 89
  set {char}($phone55 + 2) = 49
  set {char}($phone55 + 3) = 0
  set $phone56 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone56 + 0) = 73
  set {char}($phone56 + 1) = 89
  set {char}($phone56 + 2) = 50
  set {char}($phone56 + 3) = 0
  set $phone57 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone57 + 0) = 79
  set {char}($phone57 + 1) = 87
  set {char}($phone57 + 2) = 48
  set {char}($phone57 + 3) = 0
  set $phone58 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone58 + 0) = 79
  set {char}($phone58 + 1) = 87
  set {char}($phone58 + 2) = 49
  set {char}($phone58 + 3) = 0
  set $phone59 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone59 + 0) = 79
  set {char}($phone59 + 1) = 87
  set {char}($phone59 + 2) = 50
  set {char}($phone59 + 3) = 0
  set $phone60 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone60 + 0) = 79
  set {char}($phone60 + 1) = 89
  set {char}($phone60 + 2) = 48
  set {char}($phone60 + 3) = 0
  set $phone61 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone61 + 0) = 79
  set {char}($phone61 + 1) = 89
  set {char}($phone61 + 2) = 49
  set {char}($phone61 + 3) = 0
  set $phone62 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone62 + 0) = 79
  set {char}($phone62 + 1) = 89
  set {char}($phone62 + 2) = 50
  set {char}($phone62 + 3) = 0
  set $phone63 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone63 + 0) = 85
  set {char}($phone63 + 1) = 72
  set {char}($phone63 + 2) = 48
  set {char}($phone63 + 3) = 0
  set $phone64 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone64 + 0) = 85
  set {char}($phone64 + 1) = 72
  set {char}($phone64 + 2) = 49
  set {char}($phone64 + 3) = 0
  set $phone65 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone65 + 0) = 85
  set {char}($phone65 + 1) = 72
  set {char}($phone65 + 2) = 50
  set {char}($phone65 + 3) = 0
  set $phone66 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone66 + 0) = 85
  set {char}($phone66 + 1) = 87
  set {char}($phone66 + 2) = 48
  set {char}($phone66 + 3) = 0
  set $phone67 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone67 + 0) = 85
  set {char}($phone67 + 1) = 87
  set {char}($phone67 + 2) = 49
  set {char}($phone67 + 3) = 0
  set $phone68 = ((char *(*)(unsigned int))0x1001d9c0)(4)
  set {char}($phone68 + 0) = 85
  set {char}($phone68 + 1) = 87
  set {char}($phone68 + 2) = 50
  set {char}($phone68 + 3) = 0
  set $phone69 = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set {char}($phone69 + 0) = 35
  set {char}($phone69 + 1) = 0
  set $output = ((char *(*)(unsigned int))0x1001d9c0)(66)
  set $i = 0
  while $i < 66
    set {char}($output + $i) = 0
    set $i = $i + 1
  end
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone0)
  printf "TARGETPHON converter token=B result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone1)
  printf "TARGETPHON converter token=D result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone2)
  printf "TARGETPHON converter token=F result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone3)
  printf "TARGETPHON converter token=G result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone4)
  printf "TARGETPHON converter token=K result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone5)
  printf "TARGETPHON converter token=L result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone6)
  printf "TARGETPHON converter token=M result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone7)
  printf "TARGETPHON converter token=N result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone8)
  printf "TARGETPHON converter token=P result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone9)
  printf "TARGETPHON converter token=R result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone10)
  printf "TARGETPHON converter token=S result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone11)
  printf "TARGETPHON converter token=T result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone12)
  printf "TARGETPHON converter token=V result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone13)
  printf "TARGETPHON converter token=W result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone14)
  printf "TARGETPHON converter token=Y result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone15)
  printf "TARGETPHON converter token=Z result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone16)
  printf "TARGETPHON converter token=CH result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone17)
  printf "TARGETPHON converter token=DH result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone18)
  printf "TARGETPHON converter token=HH result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone19)
  printf "TARGETPHON converter token=JH result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone20)
  printf "TARGETPHON converter token=NG result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone21)
  printf "TARGETPHON converter token=SH result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone22)
  printf "TARGETPHON converter token=TH result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone23)
  printf "TARGETPHON converter token=ZH result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone24)
  printf "TARGETPHON converter token=AA0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone25)
  printf "TARGETPHON converter token=AA1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone26)
  printf "TARGETPHON converter token=AA2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone27)
  printf "TARGETPHON converter token=AE0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone28)
  printf "TARGETPHON converter token=AE1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone29)
  printf "TARGETPHON converter token=AE2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone30)
  printf "TARGETPHON converter token=AH0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone31)
  printf "TARGETPHON converter token=AH1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone32)
  printf "TARGETPHON converter token=AH2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone33)
  printf "TARGETPHON converter token=AO0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone34)
  printf "TARGETPHON converter token=AO1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone35)
  printf "TARGETPHON converter token=AO2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone36)
  printf "TARGETPHON converter token=AW0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone37)
  printf "TARGETPHON converter token=AW1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone38)
  printf "TARGETPHON converter token=AW2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone39)
  printf "TARGETPHON converter token=AY0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone40)
  printf "TARGETPHON converter token=AY1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone41)
  printf "TARGETPHON converter token=AY2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone42)
  printf "TARGETPHON converter token=EH0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone43)
  printf "TARGETPHON converter token=EH1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone44)
  printf "TARGETPHON converter token=EH2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone45)
  printf "TARGETPHON converter token=ER0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone46)
  printf "TARGETPHON converter token=ER1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone47)
  printf "TARGETPHON converter token=ER2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone48)
  printf "TARGETPHON converter token=EY0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone49)
  printf "TARGETPHON converter token=EY1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone50)
  printf "TARGETPHON converter token=EY2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone51)
  printf "TARGETPHON converter token=IH0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone52)
  printf "TARGETPHON converter token=IH1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone53)
  printf "TARGETPHON converter token=IH2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone54)
  printf "TARGETPHON converter token=IY0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone55)
  printf "TARGETPHON converter token=IY1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone56)
  printf "TARGETPHON converter token=IY2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone57)
  printf "TARGETPHON converter token=OW0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone58)
  printf "TARGETPHON converter token=OW1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone59)
  printf "TARGETPHON converter token=OW2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone60)
  printf "TARGETPHON converter token=OY0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone61)
  printf "TARGETPHON converter token=OY1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone62)
  printf "TARGETPHON converter token=OY2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone63)
  printf "TARGETPHON converter token=UH0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone64)
  printf "TARGETPHON converter token=UH1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone65)
  printf "TARGETPHON converter token=UH2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone66)
  printf "TARGETPHON converter token=UW0 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone67)
  printf "TARGETPHON converter token=UW1 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone68)
  printf "TARGETPHON converter token=UW2 result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  set $result = ((short (*)(char *, char *))0x1005f710)($output, $phone69)
  printf "TARGETPHON converter token=# result=%d code=%02x next=%02x\n", $result, *(unsigned char *)$output, *((unsigned char *)$output + 1)
  call ((void (*)(void *))0x1001da30)($phone0)
  call ((void (*)(void *))0x1001da30)($phone1)
  call ((void (*)(void *))0x1001da30)($phone2)
  call ((void (*)(void *))0x1001da30)($phone3)
  call ((void (*)(void *))0x1001da30)($phone4)
  call ((void (*)(void *))0x1001da30)($phone5)
  call ((void (*)(void *))0x1001da30)($phone6)
  call ((void (*)(void *))0x1001da30)($phone7)
  call ((void (*)(void *))0x1001da30)($phone8)
  call ((void (*)(void *))0x1001da30)($phone9)
  call ((void (*)(void *))0x1001da30)($phone10)
  call ((void (*)(void *))0x1001da30)($phone11)
  call ((void (*)(void *))0x1001da30)($phone12)
  call ((void (*)(void *))0x1001da30)($phone13)
  call ((void (*)(void *))0x1001da30)($phone14)
  call ((void (*)(void *))0x1001da30)($phone15)
  call ((void (*)(void *))0x1001da30)($phone16)
  call ((void (*)(void *))0x1001da30)($phone17)
  call ((void (*)(void *))0x1001da30)($phone18)
  call ((void (*)(void *))0x1001da30)($phone19)
  call ((void (*)(void *))0x1001da30)($phone20)
  call ((void (*)(void *))0x1001da30)($phone21)
  call ((void (*)(void *))0x1001da30)($phone22)
  call ((void (*)(void *))0x1001da30)($phone23)
  call ((void (*)(void *))0x1001da30)($phone24)
  call ((void (*)(void *))0x1001da30)($phone25)
  call ((void (*)(void *))0x1001da30)($phone26)
  call ((void (*)(void *))0x1001da30)($phone27)
  call ((void (*)(void *))0x1001da30)($phone28)
  call ((void (*)(void *))0x1001da30)($phone29)
  call ((void (*)(void *))0x1001da30)($phone30)
  call ((void (*)(void *))0x1001da30)($phone31)
  call ((void (*)(void *))0x1001da30)($phone32)
  call ((void (*)(void *))0x1001da30)($phone33)
  call ((void (*)(void *))0x1001da30)($phone34)
  call ((void (*)(void *))0x1001da30)($phone35)
  call ((void (*)(void *))0x1001da30)($phone36)
  call ((void (*)(void *))0x1001da30)($phone37)
  call ((void (*)(void *))0x1001da30)($phone38)
  call ((void (*)(void *))0x1001da30)($phone39)
  call ((void (*)(void *))0x1001da30)($phone40)
  call ((void (*)(void *))0x1001da30)($phone41)
  call ((void (*)(void *))0x1001da30)($phone42)
  call ((void (*)(void *))0x1001da30)($phone43)
  call ((void (*)(void *))0x1001da30)($phone44)
  call ((void (*)(void *))0x1001da30)($phone45)
  call ((void (*)(void *))0x1001da30)($phone46)
  call ((void (*)(void *))0x1001da30)($phone47)
  call ((void (*)(void *))0x1001da30)($phone48)
  call ((void (*)(void *))0x1001da30)($phone49)
  call ((void (*)(void *))0x1001da30)($phone50)
  call ((void (*)(void *))0x1001da30)($phone51)
  call ((void (*)(void *))0x1001da30)($phone52)
  call ((void (*)(void *))0x1001da30)($phone53)
  call ((void (*)(void *))0x1001da30)($phone54)
  call ((void (*)(void *))0x1001da30)($phone55)
  call ((void (*)(void *))0x1001da30)($phone56)
  call ((void (*)(void *))0x1001da30)($phone57)
  call ((void (*)(void *))0x1001da30)($phone58)
  call ((void (*)(void *))0x1001da30)($phone59)
  call ((void (*)(void *))0x1001da30)($phone60)
  call ((void (*)(void *))0x1001da30)($phone61)
  call ((void (*)(void *))0x1001da30)($phone62)
  call ((void (*)(void *))0x1001da30)($phone63)
  call ((void (*)(void *))0x1001da30)($phone64)
  call ((void (*)(void *))0x1001da30)($phone65)
  call ((void (*)(void *))0x1001da30)($phone66)
  call ((void (*)(void *))0x1001da30)($phone67)
  call ((void (*)(void *))0x1001da30)($phone68)
  call ((void (*)(void *))0x1001da30)($phone69)
  call ((void (*)(void *))0x1001da30)($output)
  continue
end
continue
