set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $path = ((char *(*)(unsigned int))0x1001d9c0)(35)
  set {char}($path + 0) = 90
  set {char}($path + 1) = 58
  set {char}($path + 2) = 47
  set {char}($path + 3) = 119
  set {char}($path + 4) = 111
  set {char}($path + 5) = 114
  set {char}($path + 6) = 107
  set {char}($path + 7) = 47
  set {char}($path + 8) = 115
  set {char}($path + 9) = 116
  set {char}($path + 10) = 97
  set {char}($path + 11) = 103
  set {char}($path + 12) = 101
  set {char}($path + 13) = 49
  set {char}($path + 14) = 54
  set {char}($path + 15) = 47
  set {char}($path + 16) = 117
  set {char}($path + 17) = 115
  set {char}($path + 18) = 101
  set {char}($path + 19) = 114
  set {char}($path + 20) = 100
  set {char}($path + 21) = 105
  set {char}($path + 22) = 99
  set {char}($path + 23) = 116
  set {char}($path + 24) = 45
  set {char}($path + 25) = 112
  set {char}($path + 26) = 114
  set {char}($path + 27) = 111
  set {char}($path + 28) = 98
  set {char}($path + 29) = 101
  set {char}($path + 30) = 46
  set {char}($path + 31) = 99
  set {char}($path + 32) = 115
  set {char}($path + 33) = 118
  set {char}($path + 34) = 0
  set $csv = ((char *(*)(unsigned int))0x1001d9c0)(11)
  set {char}($csv + 0) = 104
  set {char}($csv + 1) = 101
  set {char}($csv + 2) = 108
  set {char}($csv + 3) = 108
  set {char}($csv + 4) = 111
  set {char}($csv + 5) = 44
  set {char}($csv + 6) = 72
  set {char}($csv + 7) = 72
  set {char}($csv + 8) = 44
  set {char}($csv + 9) = 80
  set {char}($csv + 10) = 0
  set $result_1 = ((short (*)(int, char *, char *, int))0x10027880)(300, 0, 0, 0)
  printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_0 index=300 ax=%d\n", $result_1
  if $result_1 == 1
    set $unload_1 = ((short (*)(int))0x10027980)(300)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_0 index=300 unload_ax=%d\n", $unload_1
  else
    set $failed_unload_1 = ((short (*)(int))0x10027980)(300)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_0 index=300 empty_slot_unload_ax=%%d\n", $failed_unload_1
  end
  set $result_2 = ((short (*)(int, char *, char *, int))0x10027880)(301, 0, 0, 1)
  printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_1 index=301 ax=%d\n", $result_2
  if $result_2 == 1
    set $unload_2 = ((short (*)(int))0x10027980)(301)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_1 index=301 unload_ax=%d\n", $unload_2
  else
    set $failed_unload_2 = ((short (*)(int))0x10027980)(301)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_1 index=301 empty_slot_unload_ax=%%d\n", $failed_unload_2
  end
  set $result_3 = ((short (*)(int, char *, char *, int))0x10027880)(302, 0, 0, 9)
  printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_9 index=302 ax=%d\n", $result_3
  if $result_3 == 1
    set $unload_3 = ((short (*)(int))0x10027980)(302)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_9 index=302 unload_ax=%d\n", $unload_3
  else
    set $failed_unload_3 = ((short (*)(int))0x10027980)(302)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_9 index=302 empty_slot_unload_ax=%%d\n", $failed_unload_3
  end
  set $result_4 = ((short (*)(int, char *, char *, int))0x10027880)(303, 0, 0, 10)
  printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_10 index=303 ax=%d\n", $result_4
  if $result_4 == 1
    set $unload_4 = ((short (*)(int))0x10027980)(303)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_10 index=303 unload_ax=%d\n", $unload_4
  else
    set $failed_unload_4 = ((short (*)(int))0x10027980)(303)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_10 index=303 empty_slot_unload_ax=%%d\n", $failed_unload_4
  end
  set $result_5 = ((short (*)(int, char *, char *, int))0x10027880)(304, 0, 0, 11)
  printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_11 index=304 ax=%d\n", $result_5
  if $result_5 == 1
    set $unload_5 = ((short (*)(int))0x10027980)(304)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_11 index=304 unload_ax=%d\n", $unload_5
  else
    set $failed_unload_5 = ((short (*)(int))0x10027980)(304)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_11 index=304 empty_slot_unload_ax=%%d\n", $failed_unload_5
  end
  set $result_6 = ((short (*)(int, char *, char *, int))0x10027880)(305, 0, 0, -1)
  printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_1_neg index=305 ax=%d\n", $result_6
  if $result_6 == 1
    set $unload_6 = ((short (*)(int))0x10027980)(305)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_1_neg index=305 unload_ax=%d\n", $unload_6
  else
    set $failed_unload_6 = ((short (*)(int))0x10027980)(305)
    printf "USERDICT_EXT_MATRIX case=p_null_b_null_n_1_neg index=305 empty_slot_unload_ax=%%d\n", $failed_unload_6
  end
  set $result_7 = ((short (*)(int, char *, char *, int))0x10027880)(306, 0, $csv, 0)
  printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_0 index=306 ax=%d\n", $result_7
  if $result_7 == 1
    set $unload_7 = ((short (*)(int))0x10027980)(306)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_0 index=306 unload_ax=%d\n", $unload_7
  else
    set $failed_unload_7 = ((short (*)(int))0x10027980)(306)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_0 index=306 empty_slot_unload_ax=%%d\n", $failed_unload_7
  end
  set $result_8 = ((short (*)(int, char *, char *, int))0x10027880)(307, 0, $csv, 1)
  printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_1 index=307 ax=%d\n", $result_8
  if $result_8 == 1
    set $unload_8 = ((short (*)(int))0x10027980)(307)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_1 index=307 unload_ax=%d\n", $unload_8
  else
    set $failed_unload_8 = ((short (*)(int))0x10027980)(307)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_1 index=307 empty_slot_unload_ax=%%d\n", $failed_unload_8
  end
  set $result_9 = ((short (*)(int, char *, char *, int))0x10027880)(308, 0, $csv, 9)
  printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_9 index=308 ax=%d\n", $result_9
  if $result_9 == 1
    set $unload_9 = ((short (*)(int))0x10027980)(308)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_9 index=308 unload_ax=%d\n", $unload_9
  else
    set $failed_unload_9 = ((short (*)(int))0x10027980)(308)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_9 index=308 empty_slot_unload_ax=%%d\n", $failed_unload_9
  end
  set $result_10 = ((short (*)(int, char *, char *, int))0x10027880)(309, 0, $csv, 10)
  printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_10 index=309 ax=%d\n", $result_10
  if $result_10 == 1
    set $unload_10 = ((short (*)(int))0x10027980)(309)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_10 index=309 unload_ax=%d\n", $unload_10
  else
    set $failed_unload_10 = ((short (*)(int))0x10027980)(309)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_10 index=309 empty_slot_unload_ax=%%d\n", $failed_unload_10
  end
  set $result_11 = ((short (*)(int, char *, char *, int))0x10027880)(310, 0, $csv, 11)
  printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_11 index=310 ax=%d\n", $result_11
  if $result_11 == 1
    set $unload_11 = ((short (*)(int))0x10027980)(310)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_11 index=310 unload_ax=%d\n", $unload_11
  else
    set $failed_unload_11 = ((short (*)(int))0x10027980)(310)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_11 index=310 empty_slot_unload_ax=%%d\n", $failed_unload_11
  end
  set $result_12 = ((short (*)(int, char *, char *, int))0x10027880)(311, 0, $csv, -1)
  printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_1_neg index=311 ax=%d\n", $result_12
  if $result_12 == 1
    set $unload_12 = ((short (*)(int))0x10027980)(311)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_1_neg index=311 unload_ax=%d\n", $unload_12
  else
    set $failed_unload_12 = ((short (*)(int))0x10027980)(311)
    printf "USERDICT_EXT_MATRIX case=p_null_b_data_n_1_neg index=311 empty_slot_unload_ax=%%d\n", $failed_unload_12
  end
  set $result_13 = ((short (*)(int, char *, char *, int))0x10027880)(312, $path, 0, 0)
  printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_0 index=312 ax=%d\n", $result_13
  if $result_13 == 1
    set $unload_13 = ((short (*)(int))0x10027980)(312)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_0 index=312 unload_ax=%d\n", $unload_13
  else
    set $failed_unload_13 = ((short (*)(int))0x10027980)(312)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_0 index=312 empty_slot_unload_ax=%%d\n", $failed_unload_13
  end
  set $result_14 = ((short (*)(int, char *, char *, int))0x10027880)(313, $path, 0, 1)
  printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_1 index=313 ax=%d\n", $result_14
  if $result_14 == 1
    set $unload_14 = ((short (*)(int))0x10027980)(313)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_1 index=313 unload_ax=%d\n", $unload_14
  else
    set $failed_unload_14 = ((short (*)(int))0x10027980)(313)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_1 index=313 empty_slot_unload_ax=%%d\n", $failed_unload_14
  end
  set $result_15 = ((short (*)(int, char *, char *, int))0x10027880)(314, $path, 0, 9)
  printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_9 index=314 ax=%d\n", $result_15
  if $result_15 == 1
    set $unload_15 = ((short (*)(int))0x10027980)(314)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_9 index=314 unload_ax=%d\n", $unload_15
  else
    set $failed_unload_15 = ((short (*)(int))0x10027980)(314)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_9 index=314 empty_slot_unload_ax=%%d\n", $failed_unload_15
  end
  set $result_16 = ((short (*)(int, char *, char *, int))0x10027880)(315, $path, 0, 10)
  printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_10 index=315 ax=%d\n", $result_16
  if $result_16 == 1
    set $unload_16 = ((short (*)(int))0x10027980)(315)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_10 index=315 unload_ax=%d\n", $unload_16
  else
    set $failed_unload_16 = ((short (*)(int))0x10027980)(315)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_10 index=315 empty_slot_unload_ax=%%d\n", $failed_unload_16
  end
  set $result_17 = ((short (*)(int, char *, char *, int))0x10027880)(316, $path, 0, 11)
  printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_11 index=316 ax=%d\n", $result_17
  if $result_17 == 1
    set $unload_17 = ((short (*)(int))0x10027980)(316)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_11 index=316 unload_ax=%d\n", $unload_17
  else
    set $failed_unload_17 = ((short (*)(int))0x10027980)(316)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_11 index=316 empty_slot_unload_ax=%%d\n", $failed_unload_17
  end
  set $result_18 = ((short (*)(int, char *, char *, int))0x10027880)(317, $path, 0, -1)
  printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_1_neg index=317 ax=%d\n", $result_18
  if $result_18 == 1
    set $unload_18 = ((short (*)(int))0x10027980)(317)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_1_neg index=317 unload_ax=%d\n", $unload_18
  else
    set $failed_unload_18 = ((short (*)(int))0x10027980)(317)
    printf "USERDICT_EXT_MATRIX case=p_file_b_null_n_1_neg index=317 empty_slot_unload_ax=%%d\n", $failed_unload_18
  end
  set $result_19 = ((short (*)(int, char *, char *, int))0x10027880)(318, $path, $csv, 0)
  printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_0 index=318 ax=%d\n", $result_19
  if $result_19 == 1
    set $unload_19 = ((short (*)(int))0x10027980)(318)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_0 index=318 unload_ax=%d\n", $unload_19
  else
    set $failed_unload_19 = ((short (*)(int))0x10027980)(318)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_0 index=318 empty_slot_unload_ax=%%d\n", $failed_unload_19
  end
  set $result_20 = ((short (*)(int, char *, char *, int))0x10027880)(319, $path, $csv, 1)
  printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_1 index=319 ax=%d\n", $result_20
  if $result_20 == 1
    set $unload_20 = ((short (*)(int))0x10027980)(319)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_1 index=319 unload_ax=%d\n", $unload_20
  else
    set $failed_unload_20 = ((short (*)(int))0x10027980)(319)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_1 index=319 empty_slot_unload_ax=%%d\n", $failed_unload_20
  end
  set $result_21 = ((short (*)(int, char *, char *, int))0x10027880)(320, $path, $csv, 9)
  printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_9 index=320 ax=%d\n", $result_21
  if $result_21 == 1
    set $unload_21 = ((short (*)(int))0x10027980)(320)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_9 index=320 unload_ax=%d\n", $unload_21
  else
    set $failed_unload_21 = ((short (*)(int))0x10027980)(320)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_9 index=320 empty_slot_unload_ax=%%d\n", $failed_unload_21
  end
  set $result_22 = ((short (*)(int, char *, char *, int))0x10027880)(321, $path, $csv, 10)
  printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_10 index=321 ax=%d\n", $result_22
  if $result_22 == 1
    set $unload_22 = ((short (*)(int))0x10027980)(321)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_10 index=321 unload_ax=%d\n", $unload_22
  else
    set $failed_unload_22 = ((short (*)(int))0x10027980)(321)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_10 index=321 empty_slot_unload_ax=%%d\n", $failed_unload_22
  end
  set $result_23 = ((short (*)(int, char *, char *, int))0x10027880)(322, $path, $csv, 11)
  printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_11 index=322 ax=%d\n", $result_23
  if $result_23 == 1
    set $unload_23 = ((short (*)(int))0x10027980)(322)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_11 index=322 unload_ax=%d\n", $unload_23
  else
    set $failed_unload_23 = ((short (*)(int))0x10027980)(322)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_11 index=322 empty_slot_unload_ax=%%d\n", $failed_unload_23
  end
  set $result_24 = ((short (*)(int, char *, char *, int))0x10027880)(323, $path, $csv, -1)
  printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_1_neg index=323 ax=%d\n", $result_24
  if $result_24 == 1
    set $unload_24 = ((short (*)(int))0x10027980)(323)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_1_neg index=323 unload_ax=%d\n", $unload_24
  else
    set $failed_unload_24 = ((short (*)(int))0x10027980)(323)
    printf "USERDICT_EXT_MATRIX case=p_file_b_data_n_1_neg index=323 empty_slot_unload_ax=%%d\n", $failed_unload_24
  end
  continue
end

continue
