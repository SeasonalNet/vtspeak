set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (unsigned char *)malloc(16)
  set $passed = 0
  set $failed = 0
  set $value = 0
  while $value < 256
    set $line[0] = 65
    set $line[1] = $value
    set $line[2] = 66
    set $line[3] = 44
    set $line[4] = 88
    set $line[5] = 0
    set $result = ((int (*)(void *, char *, int))0x10016bd0)($parser, (char *)$line, 0)
    set $count = ((int (*)(void *))0x10016c10)($parser)
    set $f0 = (unsigned char *)((char *(*)(void *, int))0x10016c30)($parser, 0)
    set $ok = 0
    if $value == 0
      if $result == 1 && $count == 1 && $f0[0] == 65 && $f0[1] == 0
        set $ok = 1
      end
      printf "CSV_BYTE value=000 result=%d count=%d f0=%02x\n", $result, $count, $f0[0]
    else
      if $value == 44
        set $f1 = (unsigned char *)((char *(*)(void *, int))0x10016c30)($parser, 1)
        set $f2 = (unsigned char *)((char *(*)(void *, int))0x10016c30)($parser, 2)
        if $result == 1 && $count == 3 && $f0[0] == 65 && $f0[1] == 0 && $f1[0] == 66 && $f1[1] == 0 && $f2[0] == 88 && $f2[1] == 0
          set $ok = 1
        end
        printf "CSV_BYTE value=%03u result=%d count=%d f0=%02x f1=%02x f2=%02x\n", $value, $result, $count, $f0[0], $f1[0], $f2[0]
      else
        set $f1 = (unsigned char *)((char *(*)(void *, int))0x10016c30)($parser, 1)
        if $result == 1 && $count == 2 && $f0[0] == 65 && $f0[1] == $value && $f0[2] == 66 && $f0[3] == 0 && $f1[0] == 88 && $f1[1] == 0
          set $ok = 1
        end
        printf "CSV_BYTE value=%03u result=%d count=%d f0=%02x%02x%02x f1=%02x\n", $value, $result, $count, $f0[0], $f0[1], $f0[2], $f1[0]
      end
    end
    if $ok
      set $passed = $passed + 1
    else
      set $failed = $failed + 1
      printf "CSV_BYTE_FAIL value=%03u result=%d count=%d\n", $value, $result, $count
    end
    set $value = $value + 1
  end
  printf "CSV_BYTE_DOMAIN passed=%d failed=%d\n", $passed, $failed
  call ((void (*)(void *))0x10016bf0)($parser)
  continue
end

continue
