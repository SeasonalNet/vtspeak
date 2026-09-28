set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $raw = (unsigned char *)malloc(9)
  set $out = $raw + 1
  set $i = 0
  while $i < 9
    set {unsigned char}($raw + $i) = 0xa5
    set $i = $i + 1
  end
  set {unsigned char}$raw = 0x5a
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)(0, -1, $out, 8)
  printf "CSV_MAKE_NEGATIVE count=-1 low_ax=%d bytes=", $ret
  set $i = 0
  while $i < 8
    printf "%02x", $out[$i]
    set $i = $i + 1
  end
  printf " prefix_guard=%02x\n", $raw[0]

  set $i = 0
  while $i < 9
    set {unsigned char}($raw + $i) = 0xa5
    set $i = $i + 1
  end
  set {unsigned char}$raw = 0x5a
  set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)(0, -2147483648, $out, 8)
  printf "CSV_MAKE_INT_MIN low_ax=%d bytes=", $ret
  set $i = 0
  while $i < 8
    printf "%02x", $out[$i]
    set $i = $i + 1
  end
  printf " prefix_guard=%02x\n", $raw[0]

  continue
end

continue
