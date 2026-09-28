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
  set $raw = (unsigned char *)malloc(65540)
  set $out = $raw + 1
  set $i = 0
  while $i < 5
    if $i == 0
      set $cap = 65
    end
    if $i == 1
      set $cap = 128
    end
    if $i == 2
      set $cap = 1024
    end
    if $i == 3
      set $cap = 4096
    end
    if $i == 4
      set $cap = 65536
    end
    set {unsigned char}$raw = 0x5a
    set {unsigned char}($out + $cap) = 0xa5
    set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 2, $out, $cap)
    printf "CSV_MAKE_LARGE capacity=%u low_ax=%d prefix=%02x first=", $cap, $ret, $raw[0]
    set $j = 0
    while $j < 12
      printf "%02x", $out[$j]
      set $j = $j + 1
    end
    printf " final=%02x post_guard=%02x\n", $out[$cap - 1], $out[$cap]
    set $i = $i + 1
  end
  continue
end

continue
