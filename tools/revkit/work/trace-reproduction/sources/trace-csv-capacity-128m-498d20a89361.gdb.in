set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $fields = (char **)malloc(8)
  set $field0 = (char *)malloc(8)
  set $field1 = (char *)malloc(8)
  set {char[2]}$field0 = "A"
  set {char[4]}$field1 = "b,c"
  set $fields[0] = $field0
  set $fields[1] = $field1
  set $raw = (unsigned char *)malloc(134217730)
  set $out = $raw + 1
  set {unsigned char}$raw = 0x5a
  set {unsigned char}($out + 134217728) = 0xa5
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 2, $out, 134217728)
  printf "CSV_MAKE_128M low_ax=%d prefix=%02x output=", $ret, $raw[0]
  set $i = 0
  while $i < 10
    printf "%02x", $out[$i]
    set $i = $i + 1
  end
  printf " capacity_final=%02x post_guard=%02x\n", $out[134217727], $out[134217728]
  continue
end

continue
