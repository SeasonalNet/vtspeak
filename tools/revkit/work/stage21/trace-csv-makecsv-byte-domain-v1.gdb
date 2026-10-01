set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $fields = (char **)malloc(4)
  set $field = (unsigned char *)malloc(2)
  set $raw = (unsigned char *)malloc(10)
  set $out = $raw + 1
  set $fields[0] = (char *)$field
  set $passed = 0
  set $failed = 0
  set $value = 0
  while $value < 256
    set $field[0] = $value
    set $field[1] = 0
    set {unsigned char}$raw = 0x5a
    set {unsigned char}($out + 8) = 0xa5
    set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
    set $ok = 0
    if $value == 0
      if $ret == -1 && $out[0] == 0x22 && $out[1] == 0x31 && $out[2] == 0x31 && $out[3] == 0x31 && $out[4] == 0x31 && $out[5] == 0x31 && $out[6] == 0x31 && $out[7] == 0 && $raw[0] == 0x5a && $out[8] == 0xa5
        set $ok = 1
      end
    else
      if $value == 0x22
        if $ret == 1 && $out[0] == 0x22 && $out[1] == 0x22 && $out[2] == 0x22 && $out[3] == 0x22 && $out[4] == 0 && $out[5] == 0x31 && $out[6] == 0x31 && $out[7] == 0 && $raw[0] == 0x5a && $out[8] == 0xa5
          set $ok = 1
        end
      else
        if $ret == 1 && $out[0] == 0x22 && $out[1] == $value && $out[2] == 0x22 && $out[3] == 0 && $out[4] == 0x31 && $out[5] == 0x31 && $out[6] == 0x31 && $out[7] == 0 && $raw[0] == 0x5a && $out[8] == 0xa5
          set $ok = 1
        end
      end
    end
    if $ok
      set $passed = $passed + 1
    else
      set $failed = $failed + 1
      printf "CSV_MAKE_BYTE_FAIL value=%02x ret=%d out=%02x%02x%02x%02x%02x%02x%02x final=%02x prefix=%02x post=%02x\n", $value, $ret, $out[0], $out[1], $out[2], $out[3], $out[4], $out[5], $out[6], $out[7], $raw[0], $out[8]
    end
    printf "CSV_MAKE_BYTE value=%02x ret=%d output=%02x%02x%02x%02x%02x%02x%02x final=%02x prefix=%02x post=%02x\n", $value, $ret, $out[0], $out[1], $out[2], $out[3], $out[4], $out[5], $out[6], $out[7], $raw[0], $out[8]
    set $value = $value + 1
  end
  printf "CSV_MAKE_BYTE_DOMAIN passed=%d failed=%d\n", $passed, $failed
  continue
end

continue
