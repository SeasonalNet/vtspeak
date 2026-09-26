set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $heap_export = *(unsigned int *)0x100ff11c
  set $ver1 = *(unsigned int *)0x1007d6f8
  set $ver2 = *(unsigned int *)0x1007d6fc
  set $ver3 = *(unsigned int *)0x1007d700
  set $ver4 = *(unsigned int *)0x1007d704
  printf "EXPORTED_DATA heap_start_address=0x%08x version=(%u,%u,%u,%u)\n", $heap_export, $ver1, $ver2, $ver3, $ver4

  set $csv = ((void *(*)(void))0x10016bc0)()
  printf "CSV_INIT pointer=%p\n", $csv
  set $csvline = (char *)malloc(128)
  set {char[25]}$csvline = "alpha,\"beta,gamma\",delta"
  set $csv_parse_result = ((int (*)(void *, char *, int))0x10016bd0)($csv, $csvline, 0)
  set $csv_fields = ((int (*)(void *))0x10016c10)($csv)
  printf "CSV_PARSE result=%d field_count=%d\n", $csv_parse_result, $csv_fields
  set $csv_index = 0
  while $csv_index < $csv_fields
    set $csv_field = ((char *(*)(void *, int))0x10016c30)($csv, $csv_index)
    printf "CSV_FIELD index=%d value=%s\n", $csv_index, $csv_field
    set $csv_index = $csv_index + 1
  end
  set $csv_oob = ((char *(*)(void *, int))0x10016c30)($csv, $csv_fields)
  printf "CSV_FIELD out_of_range=%p\n", $csv_oob
  call ((void (*)(void *))0x10016bf0)($csv)
  printf "CSV_EXIT completed\n"

  set $csv_edge = ((void *(*)(void))0x10016bc0)()
  set $csv_edge_input = (char *)malloc(128)
  set {char[32]}$csv_edge_input = "a,\"unterminated"
  set $csv_unclosed_result = ((int (*)(void *, char *, int))0x10016bd0)($csv_edge, $csv_edge_input, 0)
  set $csv_unclosed_fields = ((int (*)(void *))0x10016c10)($csv_edge)
  set $csv_unclosed0 = ((char *(*)(void *, int))0x10016c30)($csv_edge, 0)
  set $csv_unclosed1 = ((char *(*)(void *, int))0x10016c30)($csv_edge, 1)
  printf "CSV_EDGE unclosed_quote result=%d field_count=%d fields=\"%s\",\"%s\"\n", $csv_unclosed_result, $csv_unclosed_fields, $csv_unclosed0, $csv_unclosed1
  call ((void (*)(void *))0x10016bf0)($csv_edge)

  set $csv_edge = ((void *(*)(void))0x10016bc0)()
  set {char[32]}$csv_edge_input = "a,,c"
  set $csv_empty_field_result = ((int (*)(void *, char *, int))0x10016bd0)($csv_edge, $csv_edge_input, 0)
  set $csv_empty_field_count = ((int (*)(void *))0x10016c10)($csv_edge)
  set $csv_empty_field = ((char *(*)(void *, int))0x10016c30)($csv_edge, 1)
  printf "CSV_EDGE empty_middle_field result=%d field_count=%d middle=\"%s\"\n", $csv_empty_field_result, $csv_empty_field_count, $csv_empty_field
  call ((void (*)(void *))0x10016bf0)($csv_edge)

  set $csv_edge = ((void *(*)(void))0x10016bc0)()
  set {char[32]}$csv_edge_input = "a,\nb"
  set $csv_multiline_result = ((int (*)(void *, char *, int))0x10016bd0)($csv_edge, $csv_edge_input, 0)
  set $csv_multiline_fields = ((int (*)(void *))0x10016c10)($csv_edge)
  set $csv_multiline0 = ((char *(*)(void *, int))0x10016c30)($csv_edge, 0)
  set $csv_multiline1 = ((char *(*)(void *, int))0x10016c30)($csv_edge, 1)
  printf "CSV_EDGE embedded_newline result=%d field_count=%d field0=\"%s\" field1_first=0x%02x field1_second=0x%02x\n", $csv_multiline_result, $csv_multiline_fields, $csv_multiline0, *(unsigned char *)$csv_multiline1, *(unsigned char *)($csv_multiline1 + 1)
  call ((void (*)(void *))0x10016bf0)($csv_edge)

  set $path_slot = 0
  while $path_slot < 6
    set $path_key = ((char *(*)(int))0x1001c690)($path_slot)
    if ($path_key == 0)
      printf "PATH_KEY slot=%d value=NULL\n", $path_slot
    else
      printf "PATH_KEY slot=%d value=%s\n", $path_slot, $path_key
    end
    set $path_slot = $path_slot + 1
  end

  call ((void (*)(void))0x1001def0)()
  printf "USERDICT_DEFAULT_ABSOLUTE value=%s\n", (char *)0x1009fc54
  call ((void (*)(void))0x1001df00)()
  printf "USERDICT_DEFAULT_RELATIVE value=%s\n", (char *)0x1009fc54

  set $csvfield1 = (char *)malloc(16)
  set {char[2]}$csvfield1 = "A"
  set $csvfield2 = (char *)malloc(16)
  set {char[4]}$csvfield2 = "b,c"
  set $csvfield_array = (char **)malloc(8)
  set $csvfield_array[0] = $csvfield1
  set $csvfield_array[1] = $csvfield2
  set $csvout = (char *)malloc(128)
  set $csv_make_result = ((int (*)(char **, int, char *, unsigned int))0x10016c50)($csvfield_array, 2, $csvout, 128)
  printf "CSV_MAKE result=%d value=%s\n", $csv_make_result, $csvout
  set $iscsv_match = ((int (*)(char *, int, int))0x10016c70)($csvline, 3, 128)
  set $iscsv_mismatch = ((int (*)(char *, int, int))0x10016c70)($csvline, 4, 128)
  printf "CSV_IS expected_fields_3=%d expected_fields_4=%d\n", $iscsv_match, $iscsv_mismatch

  set $sync_a = ((int *(*)(void))0x10026360)()
  set $sync_b = ((int *(*)(void))0x10026360)()
  printf "SYNC_ALLOC first=%p second=%p\n", $sync_a, $sync_b
  if ($sync_a != 0) && ($sync_b != 0)
    set $sync_a[4] = 0x13579bdf
    set $sync_a[6] = 0x2468ace0
    call ((void (*)(int *))0x10026450)($sync_a)
    printf "SYNC_INIT defaults fields_1_2_4_6_13=%d,%d,%d,%d,%d\n", $sync_a[1], $sync_a[2], $sync_a[4], $sync_a[6], $sync_a[13]
    set $sync_b[4] = 0x13579bdf
    set $sync_b[6] = 0x2468ace0
    call ((void (*)(int *, int *))0x10026500)($sync_a, $sync_b)
    printf "SYNC_COPY source_fields_4_6=0x%08x,0x%08x target_fields_4_6=0x%08x,0x%08x\n", $sync_b[4], $sync_b[6], $sync_a[4], $sync_a[6]
  end
  if ($sync_a != 0)
    call ((void (*)(int *))0x100263f0)($sync_a)
  end
  if ($sync_b != 0)
    call ((void (*)(int *))0x100263f0)($sync_b)
  end
  printf "SYNC_FREE completed\n"

  set $initial_device = *(unsigned int *)0x1007d6e0
  set $old_device = ((unsigned int (*)(unsigned int))0x10026b20)(0x12345678)
  printf "SOUND_CARD_SET prior=0x%08x changed=0x%08x\n", $old_device, *(unsigned int *)0x1007d6e0
  set $restore_device_result = ((unsigned int (*)(unsigned int))0x10026b20)($initial_device)
  printf "SOUND_CARD_RESTORED value=0x%08x\n", *(unsigned int *)0x1007d6e0

  set $speaker_state = *(unsigned int *)(0x100a0464 + 4)
  if ($speaker_state != 0)
    set $emphasis_initial = *(int *)($speaker_state + 0x4cfc)
    call ((void (*)(int, int))0x10028310)(200, 1)
    printf "EMPHASIS input=200 stored=%d\n", *(int *)($speaker_state + 0x4cfc)
    call ((void (*)(int, int))0x10028310)(-200, 1)
    printf "EMPHASIS input=-200 stored=%d\n", *(int *)($speaker_state + 0x4cfc)
    call ((void (*)(int, int))0x10028310)($emphasis_initial, 1)
    printf "EMPHASIS restored=%d\n", *(int *)($speaker_state + 0x4cfc)
  end

  set $runtime_state = *(unsigned int *)0x100a0460
  if ($runtime_state != 0)
    set $highlight_initial = *(unsigned char *)($runtime_state + 0x20424)
    call ((void (*)(unsigned char))0x10028360)(7)
    printf "HIGHLIGHT input=7 stored=%u\n", *(unsigned char *)($runtime_state + 0x20424)
    call ((void (*)(unsigned char))0x10028360)($highlight_initial)
    set $paren_initial = *(int *)($runtime_state + 0x2041c)
    call ((void (*)(int))0x10028390)(-1)
    printf "PARENTHESIS input=-1 stored=%d\n", *(int *)($runtime_state + 0x2041c)
    call ((void (*)(int))0x10028390)($paren_initial)
    set $reading_initial = *(int *)($runtime_state + 0x20420)
    call ((void (*)(int))0x100283c0)(-1)
    printf "READING_RULE input=-1 stored=%d\n", *(int *)($runtime_state + 0x20420)
    call ((void (*)(int))0x100283c0)($reading_initial)
  printf "RUNTIME_CONFIG restored highlight=%u paren=%d reading=%d\n", *(unsigned char *)($runtime_state + 0x20424), *(int *)($runtime_state + 0x2041c), *(int *)($runtime_state + 0x20420)
  end

  set $history_initial = *(unsigned char *)0x100a045d
  call ((void (*)(unsigned char))0x100283f0)(1)
  printf "UNIT_HISTORY input=1 stored=%u load_flag=%u\n", *(unsigned char *)0x100a045d, *(unsigned char *)0x100a747c
  call ((void (*)(unsigned char))0x100283f0)(0)
  printf "UNIT_HISTORY input=0 stored=%u load_flag=%u\n", *(unsigned char *)0x100a045d, *(unsigned char *)0x100a747c
  call ((void (*)(unsigned char))0x100283f0)($history_initial)

  set $norm_text = (char *)malloc(128)
  set {char[8]}$norm_text = "example"
  call ((void (*)(unsigned char *))0x1002a550)((unsigned char *)$norm_text)
  printf "USERDICT_SOURCE_NORM text=%s completed=1\n", $norm_text
  set {char[8]}$norm_text = "example"
  set $norm_result = ((int (*)(unsigned char *))0x1002a570)((unsigned char *)$norm_text)
  printf "USERDICT_TARGET_NORM text=%s result=%d\n", $norm_text, $norm_result
  set {char[8]}$norm_text = "example"
  set $phon_result = ((int (*)(unsigned char *))0x1002a590)((unsigned char *)$norm_text)
  printf "USERDICT_TARGET_PHON text=%s result=%d\n", $norm_text, $phon_result

  set $text_buffer_ex_low = ((int (*)(int, unsigned char *, unsigned int *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, unsigned int, int *, int *))0x1001ddf0)(-1, (unsigned char *)0, (unsigned int *)0, (int *)0, 0, 0, 0, (unsigned int *)0, (unsigned int *)0, 0, 0, 0, 0, (int *)0, (int *)0)
  set $text_buffer_ex_high = ((int (*)(int, unsigned char *, unsigned int *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, unsigned int, int *, int *))0x1001ddf0)(3, (unsigned char *)0, (unsigned int *)0, (int *)0, 0, 0, 0, (unsigned int *)0, (unsigned int *)0, 0, 0, 0, 0, (int *)0, (int *)0)
  printf "TEXT_BUFFER_EX selector_low=-1 return=%d selector_high=3 return=%d\n", $text_buffer_ex_low, $text_buffer_ex_high

  set $license_info_0 = ((int (*)(int, char *, char *, unsigned int, unsigned char *, int))0x1002a110)(0, (char *)0, (char *)0, 0xffffffff, (unsigned char *)0, -1)
  set $license_info_12 = ((int (*)(int, char *, char *, unsigned int, unsigned char *, int))0x1002a110)(12, (char *)0, (char *)0, 0xffffffff, (unsigned char *)0, -1)
  set $license_info_cap_0 = ((int (*)(int, char *, char *, unsigned int, unsigned char *, int))0x1002a110)(0, (char *)0, (char *)0, 0xffffffff, (unsigned char *)0, 0)
  set $license_info_cap_3 = ((int (*)(int, char *, char *, unsigned int, unsigned char *, int))0x1002a110)(0, (char *)0, (char *)0, 0xffffffff, (unsigned char *)0, 3)
  set $license_info_cap_4 = ((int (*)(int, char *, char *, unsigned int, unsigned char *, int))0x1002a110)(0, (char *)0, (char *)0, 0xffffffff, (unsigned char *)0, 4)
  printf "LICENSE_INFO null_path selector_0_length=%d selector_12_length=%d capacity_0=%d capacity_3=%d capacity_4_null_output=%d\n", $license_info_0, $license_info_12, $license_info_cap_0, $license_info_cap_3, $license_info_cap_4

  set $license_buf = (char *)malloc(512)
  set *(char *)$license_buf = 0
  set $comment_result = ((unsigned int (*)(char *, char *, int))0x10029a00)((char *)0, $license_buf, 512)
  printf "LICENSE_COMMENT null_path result=%d output_redacted=1\n", $comment_result
  set $license_check = ((int (*)(char *, char *, unsigned int, unsigned char *))0x10029b80)((char *)0, (char *)0, 0xffffffff, (unsigned char *)0)
  printf "LICENSE_CHECK null_path result=%d\n", $license_check
  set $init_result = ((short (*)(void))0x1002aa50)()
  printf "INIT result=%d\n", $init_result

  continue
end

continue
