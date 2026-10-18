set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $row = (char *)malloc(64)
  set {char[64]}$row = "a,b,c"
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, -2147483648, 64)
  printf "CSV_ISCSV case=three_fields expected=INT_MIN limit=64 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, -1, 64)
  printf "CSV_ISCSV case=three_fields expected=-1 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 0, 64)
  printf "CSV_ISCSV case=three_fields expected=0 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 1, 64)
  printf "CSV_ISCSV case=three_fields expected=1 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 2, 64)
  printf "CSV_ISCSV case=three_fields expected=2 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 64)
  printf "CSV_ISCSV case=three_fields expected=3 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 4, 64)
  printf "CSV_ISCSV case=three_fields expected=4 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 2147483647, 64)
  printf "CSV_ISCSV case=three_fields expected=INT_MAX limit=64 low_ax=%d\n", $result

  set $row = (char *)malloc(64)
  set {char[64]}$row = "a,b,c"
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, -1)
  printf "CSV_ISCSV case=short_bound expected=3 limit=-1 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 0)
  printf "CSV_ISCSV case=short_bound expected=3 limit=0 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 2)
  printf "CSV_ISCSV case=short_bound expected=3 limit=2 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 4)
  printf "CSV_ISCSV case=short_bound expected=3 limit=4 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 5)
  printf "CSV_ISCSV case=short_bound expected=3 limit=5 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 6)
  printf "CSV_ISCSV case=short_bound expected=3 limit=6 low_ax=%d\n", $result

  set $row = (char *)malloc(64)
  set {char[64]}$row = "a,b\r\nc,d"
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 2, 64)
  printf "CSV_ISCSV case=crlf_first_record expected=2 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $row = (char *)malloc(64)
  set {char[64]}$row = "a,b\r\nc,d"
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 64)
  printf "CSV_ISCSV case=crlf_first_record expected=3 limit=64 low_ax=%d source=<%s>\n", $result, $row

  set $row = (char *)malloc(64)
  set {char[64]}$row = "a,b\nc,d"
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 2, 64)
  printf "CSV_ISCSV case=lf_first_record expected=2 limit=64 low_ax=%d source=<%s>\n", $result, $row
  set $row = (char *)malloc(64)
  set {char[64]}$row = "a,b\nc,d"
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, 64)
  printf "CSV_ISCSV case=lf_first_record expected=3 limit=64 low_ax=%d source=<%s>\n", $result, $row

  set $row = (char *)malloc(64)
  set {char[64]}$row = ""
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 0, 64)
  printf "CSV_ISCSV case=empty expected=0 limit=64 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 1, 64)
  printf "CSV_ISCSV case=empty expected=1 limit=64 low_ax=%d\n", $result
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 2, 64)
  printf "CSV_ISCSV case=empty expected=2 limit=64 low_ax=%d\n", $result

  set $field_count = 98
  while $field_count <= 101
    set $row = (char *)malloc(512)
    set $field = 0
    while $field < $field_count
      if $field > 0
        set *(unsigned char *)($row + 2 * $field - 1) = 44
      end
      set *(unsigned char *)($row + 2 * $field) = 65
      set $field = $field + 1
    end
    set *(unsigned char *)($row + 2 * $field_count - 1) = 0
    set $expected = $field_count - 1
    set $result = ((int (*)(char *, int, int))0x10016c70)($row, $expected, 512)
    printf "CSV_ISCSV case=field_capacity fields=%d expected=%d limit=512 low_ax=%d\n", $field_count, $expected, $result
    set $result = ((int (*)(char *, int, int))0x10016c70)($row, $field_count, 512)
    printf "CSV_ISCSV case=field_capacity fields=%d expected=%d limit=512 low_ax=%d\n", $field_count, $field_count, $result
    set $expected = $field_count + 1
    set $result = ((int (*)(char *, int, int))0x10016c70)($row, $expected, 512)
    printf "CSV_ISCSV case=field_capacity fields=%d expected=%d limit=512 low_ax=%d\n", $field_count, $expected, $result
    set $parser = ((void *(*)(void))0x10016bc0)()
    set $parse_result = ((int (*)(void *, char *, int))0x10016bd0)($parser, $row, 0)
    set $parsed_count = ((int (*)(void *))0x10016c10)($parser)
    printf "CSV_PARSE_CAPACITY input_fields=%d low_ax=%d stored_fields=%d\n", $field_count, $parse_result, $parsed_count
    call ((void (*)(void *))0x10016bf0)($parser)
    set $field_count = $field_count + 1
  end

  kill
  quit
end

continue
