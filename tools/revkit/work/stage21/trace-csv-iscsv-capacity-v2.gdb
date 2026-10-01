set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

define csv_probe
  set $fields = (int)$arg0
  set $bytes = 2 * $fields
  set $row = (char *)malloc($bytes + 2)
  set $group = 0
  while $group < ($fields / 2)
    set *(unsigned int *)($row + 4 * $group) = 0x2c412c41
    set $group = $group + 1
  end
  if ($fields & 1) != 0
    set *(unsigned char *)($row + $bytes - 2) = 65
  end
  set *(unsigned char *)($row + $bytes - 1) = 0

  set $expected100 = ((int (*)(char *, int, int))0x10016c70)($row, 100, $bytes)
  set $expected_fields = ((int (*)(char *, int, int))0x10016c70)($row, $fields, $bytes)
  set $expected_next = ((int (*)(char *, int, int))0x10016c70)($row, $fields + 1, $bytes)
  printf "CSV_CAPACITY fields=%d bytes=%d bound=%d expected_100=%d expected_fields=%d expected_next=%d\n", $fields, $bytes - 1, $bytes, $expected100, $expected_fields, $expected_next

  set $parser = ((void *(*)(void))0x10016bc0)()
  set $parse_result = ((int (*)(void *, char *, int))0x10016bd0)($parser, $row, 0)
  set $parsed_fields = ((int (*)(void *))0x10016c10)($parser)
  printf "CSV_PARSE_CAPACITY fields=%d low_ax=%d stored_fields=%d\n", $fields, $parse_result, $parsed_fields
  call ((void (*)(void *))0x10016bf0)($parser)
end

break *0x1001da50
commands
  silent
  disable 1

  csv_probe 99
  csv_probe 100
  csv_probe 101
  csv_probe 102
  csv_probe 103
  csv_probe 104
  csv_probe 105
  csv_probe 106
  csv_probe 107
  csv_probe 108
  csv_probe 109
  csv_probe 110
  csv_probe 111
  csv_probe 112
  csv_probe 113
  csv_probe 114
  csv_probe 115
  csv_probe 116
  csv_probe 117
  csv_probe 118
  csv_probe 119
  csv_probe 120
  csv_probe 121
  csv_probe 122
  csv_probe 123
  csv_probe 124
  csv_probe 125
  csv_probe 126
  csv_probe 127
  csv_probe 128
  csv_probe 255
  csv_probe 256
  csv_probe 512
  csv_probe 1024
  csv_probe 4096

  kill
  quit
end

continue
