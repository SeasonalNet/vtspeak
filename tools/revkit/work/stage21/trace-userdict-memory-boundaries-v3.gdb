set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $buf_1 = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {unsigned char}($buf_1 + 0) = 104
  set {unsigned char}($buf_1 + 1) = 101
  set {unsigned char}($buf_1 + 2) = 108
  set {unsigned char}($buf_1 + 3) = 108
  set {unsigned char}($buf_1 + 4) = 111
  set {unsigned char}($buf_1 + 5) = 44
  set {unsigned char}($buf_1 + 6) = 72
  set {unsigned char}($buf_1 + 7) = 72
  set {unsigned char}($buf_1 + 8) = 44
  set {unsigned char}($buf_1 + 9) = 80
  set {unsigned char}($buf_1 + 10) = 165
  set {unsigned char}($buf_1 + 11) = 0
  set $load_1 = ((short (*)(int, char *, char *, int))0x10027880)(400, 0, $buf_1, 10)
  printf "USERDICT_MEMORY_BOUNDARY case=ten_span_trailing_a5 index=400 span=10 physical=12 load_ax=%d\n", $load_1
  if $load_1 == 1
    set $unload_1 = ((short (*)(int))0x10027980)(400)
    printf "USERDICT_MEMORY_BOUNDARY case=ten_span_trailing_a5 index=400 unload_ax=%d\n", $unload_1
  else
  set $empty_1 = ((short (*)(int))0x10027980)(400)
    printf "USERDICT_MEMORY_BOUNDARY case=ten_span_trailing_a5 index=400 empty_slot_unload_ax=%d\n", $empty_1
  end
  set $buf_2 = ((char *(*)(unsigned int))0x1001d9c0)(12)
  set {unsigned char}($buf_2 + 0) = 104
  set {unsigned char}($buf_2 + 1) = 101
  set {unsigned char}($buf_2 + 2) = 108
  set {unsigned char}($buf_2 + 3) = 108
  set {unsigned char}($buf_2 + 4) = 111
  set {unsigned char}($buf_2 + 5) = 44
  set {unsigned char}($buf_2 + 6) = 72
  set {unsigned char}($buf_2 + 7) = 72
  set {unsigned char}($buf_2 + 8) = 44
  set {unsigned char}($buf_2 + 9) = 80
  set {unsigned char}($buf_2 + 10) = 165
  set {unsigned char}($buf_2 + 11) = 0
  set $load_2 = ((short (*)(int, char *, char *, int))0x10027880)(401, 0, $buf_2, 11)
  printf "USERDICT_MEMORY_BOUNDARY case=eleven_span_trailing_a5 index=401 span=11 physical=12 load_ax=%d\n", $load_2
  if $load_2 == 1
    set $unload_2 = ((short (*)(int))0x10027980)(401)
    printf "USERDICT_MEMORY_BOUNDARY case=eleven_span_trailing_a5 index=401 unload_ax=%d\n", $unload_2
  else
  set $empty_2 = ((short (*)(int))0x10027980)(401)
    printf "USERDICT_MEMORY_BOUNDARY case=eleven_span_trailing_a5 index=401 empty_slot_unload_ax=%d\n", $empty_2
  end
  set $buf_3 = ((char *(*)(unsigned int))0x1001d9c0)(14)
  set {unsigned char}($buf_3 + 0) = 104
  set {unsigned char}($buf_3 + 1) = 101
  set {unsigned char}($buf_3 + 2) = 108
  set {unsigned char}($buf_3 + 3) = 108
  set {unsigned char}($buf_3 + 4) = 111
  set {unsigned char}($buf_3 + 5) = 44
  set {unsigned char}($buf_3 + 6) = 72
  set {unsigned char}($buf_3 + 7) = 72
  set {unsigned char}($buf_3 + 8) = 44
  set {unsigned char}($buf_3 + 9) = 80
  set {unsigned char}($buf_3 + 10) = 0
  set {unsigned char}($buf_3 + 11) = 44
  set {unsigned char}($buf_3 + 12) = 88
  set {unsigned char}($buf_3 + 13) = 0
  set $load_3 = ((short (*)(int, char *, char *, int))0x10027880)(402, 0, $buf_3, 14)
  printf "USERDICT_MEMORY_BOUNDARY case=embedded_nul_tail index=402 span=14 physical=14 load_ax=%d\n", $load_3
  if $load_3 == 1
    set $unload_3 = ((short (*)(int))0x10027980)(402)
    printf "USERDICT_MEMORY_BOUNDARY case=embedded_nul_tail index=402 unload_ax=%d\n", $unload_3
  else
  set $empty_3 = ((short (*)(int))0x10027980)(402)
    printf "USERDICT_MEMORY_BOUNDARY case=embedded_nul_tail index=402 empty_slot_unload_ax=%d\n", $empty_3
  end
  set $buf_4 = ((char *(*)(unsigned int))0x1001d9c0)(18)
  set {unsigned char}($buf_4 + 0) = 104
  set {unsigned char}($buf_4 + 1) = 101
  set {unsigned char}($buf_4 + 2) = 108
  set {unsigned char}($buf_4 + 3) = 108
  set {unsigned char}($buf_4 + 4) = 111
  set {unsigned char}($buf_4 + 5) = 44
  set {unsigned char}($buf_4 + 6) = 72
  set {unsigned char}($buf_4 + 7) = 72
  set {unsigned char}($buf_4 + 8) = 44
  set {unsigned char}($buf_4 + 9) = 80
  set {unsigned char}($buf_4 + 10) = 10
  set {unsigned char}($buf_4 + 11) = 120
  set {unsigned char}($buf_4 + 12) = 44
  set {unsigned char}($buf_4 + 13) = 72
  set {unsigned char}($buf_4 + 14) = 72
  set {unsigned char}($buf_4 + 15) = 44
  set {unsigned char}($buf_4 + 16) = 88
  set {unsigned char}($buf_4 + 17) = 0
  set $load_4 = ((short (*)(int, char *, char *, int))0x10027880)(403, 0, $buf_4, 17)
  printf "USERDICT_MEMORY_BOUNDARY case=malformed_second_row index=403 span=17 physical=18 load_ax=%d\n", $load_4
  if $load_4 == 1
    set $unload_4 = ((short (*)(int))0x10027980)(403)
    printf "USERDICT_MEMORY_BOUNDARY case=malformed_second_row index=403 unload_ax=%d\n", $unload_4
  else
  set $empty_4 = ((short (*)(int))0x10027980)(403)
    printf "USERDICT_MEMORY_BOUNDARY case=malformed_second_row index=403 empty_slot_unload_ax=%d\n", $empty_4
  end
  continue
end

continue
