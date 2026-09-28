set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $mode = 0
  while $mode < 2
    set $parser = ((void *(*)(void))0x10016bc0)()
    set $line = (unsigned char *)malloc(32)
    set $line[0] = 99
    set $line[1] = 97
    set $line[2] = 102
    if $mode == 0
      set $line[3] = 233
      set $line[4] = 44
      set $line[5] = 110
      set $line[6] = 97
      set $line[7] = 239
      set $line[8] = 118
      set $line[9] = 101
      set $line[10] = 0
    else
      set $line[3] = 195
      set $line[4] = 169
      set $line[5] = 44
      set $line[6] = 110
      set $line[7] = 97
      set $line[8] = 195
      set $line[9] = 175
      set $line[10] = 118
      set $line[11] = 101
      set $line[12] = 0
    end
    set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, (char *)$line, 0)
    set $count = ((int (*)(void *))0x10016c10)($parser)
    printf "CSV_ENCODING mode=%d result=%d count=%d input_hex=", $mode, $parsed, $count
    set $j = 0
    while $line[$j] != 0
      printf "%02x", $line[$j]
      set $j = $j + 1
    end
    set $j = 0
    while $j < $count
      set $field = (unsigned char *)((char *(*)(void *, int))0x10016c30)($parser, $j)
      printf " field%d_hex=", $j
      set $k = 0
      while $field[$k] != 0
        printf "%02x", $field[$k]
        set $k = $k + 1
      end
      set $j = $j + 1
    end
    printf "\n"
    call ((void (*)(void *))0x10016bf0)($parser)
    set $mode = $mode + 1
  end
  continue
end

continue
