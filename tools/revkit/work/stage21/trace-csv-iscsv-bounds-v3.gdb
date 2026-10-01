set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

define csv_bound
  set $limit = (int)$arg0
  set $result = ((int (*)(char *, int, int))0x10016c70)($row, 3, $limit)
  printf "CSV_BOUND pointer=%#x limit=%d end=%#x low_ax=%d source=<%s>\n", $row, $limit, (unsigned int)$row + (unsigned int)$limit, $result, $row
end

break *0x1001da50
commands
  silent
  disable 1
  set $row = (char *)malloc(16)
  set {char[16]}$row = "a,b,c"

  csv_bound -2147483648
  csv_bound -1024
  csv_bound -1
  csv_bound 0
  csv_bound 1
  csv_bound 4
  csv_bound 5
  csv_bound 6
  csv_bound 2147483647

  set $base = (unsigned int)$row
  set $limit_end0 = (int)(0 - $base)
  set $limit_end1 = $limit_end0 + 1
  set $limit_end4 = $limit_end0 + 4
  set $limit_end5 = $limit_end0 + 5
  set $limit_wrap = $limit_end0 - 1
  csv_bound $limit_end0
  csv_bound $limit_end1
  csv_bound $limit_end4
  csv_bound $limit_end5
  csv_bound $limit_wrap

  kill
  quit
end

continue
