set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $db_size = (int *)0x100ff11c
  set $db_size_original = *$db_size
  set $db_slot = -2147483648
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = -1
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 0
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 1
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 2
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 3
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 4
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 5
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 6
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set $db_slot = 2147483647
  set $speaker_name = ((char *(*)(int))0x10028230)($db_slot)
  printf "SPEAKER_NAME_EDGE input=%d value=%s\n", $db_slot, $speaker_name
  set *$db_size = 0x13579bdf
  set $db_result = ((int (*)(int *, int))0x10028130)($db_size, $db_slot)
  printf "DB_SIZE_EDGE input=%d result=%d value=%d\n", $db_slot, $db_result, *$db_size
  set *$db_size = $db_size_original
  detach
  quit
end

continue
