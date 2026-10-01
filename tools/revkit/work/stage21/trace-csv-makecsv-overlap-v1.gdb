set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
break *0x1001da50
commands 1
  silent
  disable 1
  set $fields = (char **)malloc(4)
  set $raw = (unsigned char *)malloc(18)
  set $out = $raw + 1
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 3)
  printf "CSV_ALIAS offset=0 cap=3 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 4)
  printf "CSV_ALIAS offset=0 cap=4 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 5)
  printf "CSV_ALIAS offset=0 cap=5 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 6)
  printf "CSV_ALIAS offset=0 cap=6 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 7)
  printf "CSV_ALIAS offset=0 cap=7 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
  printf "CSV_ALIAS offset=0 cap=8 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 9)
  printf "CSV_ALIAS offset=0 cap=9 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=0 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=0 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=0 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=0 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=0 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=0 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 0)
  set $out[0] = 65
  set $out[1] = 66
  set $out[2] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=0 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 0)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 4)
  printf "CSV_ALIAS offset=1 cap=4 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 5)
  printf "CSV_ALIAS offset=1 cap=5 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 6)
  printf "CSV_ALIAS offset=1 cap=6 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 7)
  printf "CSV_ALIAS offset=1 cap=7 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
  printf "CSV_ALIAS offset=1 cap=8 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 9)
  printf "CSV_ALIAS offset=1 cap=9 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=1 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=1 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=1 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=1 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=1 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=1 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 1)
  set $out[1] = 65
  set $out[2] = 66
  set $out[3] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=1 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 1)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 5)
  printf "CSV_ALIAS offset=2 cap=5 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 6)
  printf "CSV_ALIAS offset=2 cap=6 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 7)
  printf "CSV_ALIAS offset=2 cap=7 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
  printf "CSV_ALIAS offset=2 cap=8 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 9)
  printf "CSV_ALIAS offset=2 cap=9 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=2 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=2 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=2 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=2 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=2 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=2 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 2)
  set $out[2] = 65
  set $out[3] = 66
  set $out[4] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=2 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 2)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 6)
  printf "CSV_ALIAS offset=3 cap=6 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 7)
  printf "CSV_ALIAS offset=3 cap=7 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
  printf "CSV_ALIAS offset=3 cap=8 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 9)
  printf "CSV_ALIAS offset=3 cap=9 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=3 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=3 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=3 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=3 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=3 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=3 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 3)
  set $out[3] = 65
  set $out[4] = 66
  set $out[5] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=3 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 3)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 7)
  printf "CSV_ALIAS offset=4 cap=7 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
  printf "CSV_ALIAS offset=4 cap=8 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 9)
  printf "CSV_ALIAS offset=4 cap=9 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=4 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=4 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=4 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=4 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=4 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=4 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 4)
  set $out[4] = 65
  set $out[5] = 66
  set $out[6] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=4 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 4)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
  printf "CSV_ALIAS offset=5 cap=8 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 9)
  printf "CSV_ALIAS offset=5 cap=9 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=5 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=5 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=5 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=5 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=5 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=5 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 5)
  set $out[5] = 65
  set $out[6] = 66
  set $out[7] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=5 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 5)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 9)
  printf "CSV_ALIAS offset=6 cap=9 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=6 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=6 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=6 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=6 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=6 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=6 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 6)
  set $out[6] = 65
  set $out[7] = 66
  set $out[8] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=6 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 6)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 7)
  set $out[7] = 65
  set $out[8] = 66
  set $out[9] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 10)
  printf "CSV_ALIAS offset=7 cap=10 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 7)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 7)
  set $out[7] = 65
  set $out[8] = 66
  set $out[9] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=7 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 7)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 7)
  set $out[7] = 65
  set $out[8] = 66
  set $out[9] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=7 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 7)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 7)
  set $out[7] = 65
  set $out[8] = 66
  set $out[9] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=7 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 7)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 7)
  set $out[7] = 65
  set $out[8] = 66
  set $out[9] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=7 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 7)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 7)
  set $out[7] = 65
  set $out[8] = 66
  set $out[9] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=7 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 7)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 7)
  set $out[7] = 65
  set $out[8] = 66
  set $out[9] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=7 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 7)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 8)
  set $out[8] = 65
  set $out[9] = 66
  set $out[10] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 11)
  printf "CSV_ALIAS offset=8 cap=11 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 8)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 8)
  set $out[8] = 65
  set $out[9] = 66
  set $out[10] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=8 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 8)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 8)
  set $out[8] = 65
  set $out[9] = 66
  set $out[10] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=8 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 8)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 8)
  set $out[8] = 65
  set $out[9] = 66
  set $out[10] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=8 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 8)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 8)
  set $out[8] = 65
  set $out[9] = 66
  set $out[10] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=8 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 8)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 8)
  set $out[8] = 65
  set $out[9] = 66
  set $out[10] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=8 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 8)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 9)
  set $out[9] = 65
  set $out[10] = 66
  set $out[11] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 12)
  printf "CSV_ALIAS offset=9 cap=12 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 9)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 9)
  set $out[9] = 65
  set $out[10] = 66
  set $out[11] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 13)
  printf "CSV_ALIAS offset=9 cap=13 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 9)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 9)
  set $out[9] = 65
  set $out[10] = 66
  set $out[11] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 14)
  printf "CSV_ALIAS offset=9 cap=14 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 9)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 9)
  set $out[9] = 65
  set $out[10] = 66
  set $out[11] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 15)
  printf "CSV_ALIAS offset=9 cap=15 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 9)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 18)
  set $fields[0] = (char *)($out + 9)
  set $out[9] = 65
  set $out[10] = 66
  set $out[11] = 0
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 16)
  printf "CSV_ALIAS offset=9 cap=16 ret=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x pointer_same=%d\n", $ret, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15], $raw[16], $raw[17], $fields[0] == (char *)($out + 9)
  printf "CSV_ALIAS_DONE calls=95\n"
  continue
end
continue
