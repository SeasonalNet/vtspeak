set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $fields = (char **)malloc(4)
  set $field = (unsigned char *)malloc(3)
  set $raw = (unsigned char *)malloc(10)
  set $out = $raw + 1
  set $fields[0] = (char *)$field
  set $a = 128
  set $b = 0
  set $b_values = (unsigned char *)malloc(28)
  set $b_values[0] = 0
  set $b_values[1] = 1
  set $b_values[2] = 0x22
  set $b_values[3] = 0x2c
  set $b_values[4] = 0x7f
  set $b_values[5] = 0x80
  set $b_values[6] = 0xff
  set $guard_passed = 0
  set $cases = 0
  while $a < 256
    set $n = 0
    while $n < 7
      set $b = $b_values[$n]
      set $field[0] = $a
      set $field[1] = $b
      set $field[2] = 0
      set {unsigned char}$raw = 0x5a
      set {unsigned char}($out + 8) = 0xa5
      set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
      if $raw[0] == 0x5a && $out[8] == 0xa5
        set $guard_passed = $guard_passed + 1
      else
        printf "CSV_MAKE_PAIR_GUARD_FAIL a=%02x b=%02x prefix=%02x post=%02x\n", $a, $b, $raw[0], $out[8]
      end
      printf "CSV_MAKE_PAIR a=%02x b=%02x ret=%d output=%02x%02x%02x%02x%02x%02x%02x final=%02x prefix=%02x post=%02x\n", $a, $b, $ret, $out[0], $out[1], $out[2], $out[3], $out[4], $out[5], $out[6], $out[7], $raw[0], $out[8]
      set $cases = $cases + 1
      set $n = $n + 1
    end
    set $a = $a + 1
  end
  set $lead_values = (unsigned char *)malloc(12)
  set $lead_values[0] = 0x41
  set $lead_values[1] = 0x22
  set $lead_values[2] = 0x2c
  set $n = 0
  while $n < 3
    set $a = $lead_values[$n]
    set $b = 128
    while $b < 256
      set $field[0] = $a
      set $field[1] = $b
      set $field[2] = 0
      set {unsigned char}$raw = 0x5a
      set {unsigned char}($out + 8) = 0xa5
      set $ret = ((short (*)(char **, int, unsigned char *, unsigned int))0x10016c50)($fields, 1, $out, 8)
      if $raw[0] == 0x5a && $out[8] == 0xa5
        set $guard_passed = $guard_passed + 1
      else
        printf "CSV_MAKE_PAIR_GUARD_FAIL a=%02x b=%02x prefix=%02x post=%02x\n", $a, $b, $raw[0], $out[8]
      end
      printf "CSV_MAKE_PAIR a=%02x b=%02x ret=%d output=%02x%02x%02x%02x%02x%02x%02x final=%02x prefix=%02x post=%02x\n", $a, $b, $ret, $out[0], $out[1], $out[2], $out[3], $out[4], $out[5], $out[6], $out[7], $raw[0], $out[8]
      set $cases = $cases + 1
      set $b = $b + 1
    end
    set $n = $n + 1
  end
  printf "CSV_MAKE_PAIR_DOMAIN cases=%d guards=%d\n", $cases, $guard_passed
  continue
end

continue
