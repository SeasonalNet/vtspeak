set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $fields = (char **)malloc(8)
  set $field0 = (char *)malloc(8)
  set $field1 = (char *)malloc(8)
  set {char[2]}$field0 = "A"
  set {char[4]}$field1 = "b,c"
  set $fields[0] = $field0
  set $fields[1] = $field1
  set $raw = (unsigned char *)malloc(65)
  set $out = $raw + 1

  set $cap = 0
  while $cap <= 64
    set $i = 0
    while $i < 65
      set {unsigned char}($raw + $i) = 0xa5
      set $i = $i + 1
    end
    set {unsigned char}($raw) = 0x5a
    set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 2, $out, $cap)
    printf "CSV_MAKE_CAP capacity=%u low_ax=%d bytes=", $cap, $ret
    set $i = 0
    while $i < 12
      printf "%02x", $out[$i]
      set $i = $i + 1
    end
    printf " prefix_guard=%02x guard16=%02x guard63=%02x\n", $raw[0], $out[16], $out[63]
    set $cap = $cap + 1
  end

  set {char[4]}$field0 = "a\"b"
  set {char[8]}$field1 = "plain"
  set $i = 0
  while $i < 65
    set {unsigned char}($raw + $i) = 0xa5
    set $i = $i + 1
  end
  set {unsigned char}($raw) = 0x5a
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 2, $out, 64)
  printf "CSV_MAKE_QUOTE low_ax=%d bytes=", $ret
  set $i = 0
  while $i < 18
    printf "%02x", $out[$i]
    set $i = $i + 1
  end
  printf " prefix_guard=%02x guard64=%02x\n", $raw[0], $raw[64]

  set $i = 0
  while $i < 65
    set {unsigned char}($raw + $i) = 0xa5
    set $i = $i + 1
  end
  set {unsigned char}($raw) = 0x5a
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)(0, 0, $out, 8)
  printf "CSV_MAKE_EMPTY low_ax=%d bytes=", $ret
  set $i = 0
  while $i < 8
    printf "%02x", $out[$i]
    set $i = $i + 1
  end
  printf " prefix_guard=%02x\n", $raw[0]
  continue
end

continue
