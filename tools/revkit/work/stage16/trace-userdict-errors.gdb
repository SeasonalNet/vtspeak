set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $path = *(unsigned int *)($esp + 12)
  disable 1

  set $limit0 = ((int (*)(int))0x1002a5b0)(0)
  set $limit1 = ((int (*)(int))0x1002a5b0)(1)
  set $limit2 = ((int (*)(int))0x1002a5b0)(2)
  set $limit3 = ((int (*)(int))0x1002a5b0)(3)
  set $limit4 = ((int (*)(int))0x1002a5b0)(4)
  set $limit5 = ((int (*)(int))0x1002a5b0)(5)
  printf "USERDICT_LIMITS selectors=0,1,2,3,4,5 values=%d,%d,%d,%d,%d,%d\n", $limit0, $limit1, $limit2, $limit3, $limit4, $limit5

  set $load_negative = ((short (*)(int, char *))0x10027960)(-1, (char *)$path)
  printf "USERDICT_LOAD_NEGATIVE_INDEX return=%d\n", $load_negative
  set $load_high = ((short (*)(int, char *))0x10027960)(1024, (char *)$path)
  printf "USERDICT_LOAD_HIGH_INDEX return=%d\n", $load_high
  set $load_bad_file = ((short (*)(int, char *))0x10027960)(0, (char *)$path)
  printf "USERDICT_LOAD_INVALID_FILE return=%d\n", $load_bad_file
  set $unload_negative = ((short (*)(int))0x10027a80)(-1)
  printf "USERDICT_UNLOAD_NEGATIVE_INDEX return=%d\n", $unload_negative
  set $unload_empty = ((short (*)(int))0x10027a80)(0)
  printf "USERDICT_UNLOAD_EMPTY_SLOT return=%d\n", $unload_empty
  continue
end

continue
