set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $mode = 1
  set $fields = (char **)malloc(4)
  set $buffer = (unsigned char *)malloc(16)
  set $field = (char *)malloc(8)
  set {char[4]}$field = "text"
  set $fields[0] = $field
  if $mode == 0
    printf "CSV_POINTER_CASE mode=output-null\n"
    set $result = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)(0, 0, 0, 8)
  else
    if $mode == 1
      printf "CSV_POINTER_CASE mode=field-array-null\n"
      set $result = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)(0, 1, $buffer, 8)
    else
      printf "CSV_POINTER_CASE mode=field-element-null\n"
      set $fields[0] = 0
      set $result = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $buffer, 8)
    end
  end
  printf "CSV_POINTER_RETURN low_ax=%d\n", $result
  continue
end

continue
