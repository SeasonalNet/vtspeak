set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  disable 1
  set $parser = ((void *(*)(void))0x10016bc0)()
  set $line = (char *)malloc(16)
  set $len = 0
  set $total_cases = 0
  set $low_ax_failures = 0
  while $len <= 7
    set $limit = 1
    set $k = 0
    while $k < $len
      set $limit = $limit * 3
      set $k = $k + 1
    end
    set $code = 0
    while $code < $limit
      set $remain = $code
      set $pos = 0
      while $pos < $len
        set $digit = $remain % 3
        set $remain = $remain / 3
        if $digit == 0
          set {unsigned char}($line + $pos) = 65
        else
          if $digit == 1
            set {unsigned char}($line + $pos) = 34
          else
            set {unsigned char}($line + $pos) = 44
          end
        end
        set $pos = $pos + 1
      end
      set {unsigned char}($line + $len) = 0
      set $parsed = ((int (*)(void *, char *, int))0x10016bd0)($parser, $line, 0)
      set $count = ((int (*)(void *))0x10016c10)($parser)
      printf "CSV_QDOMAIN7 case=%d:%d input=", $len, $code
      set $pos = 0
      while $pos < $len
        printf "%02x", *(unsigned char *)($line + $pos)
        set $pos = $pos + 1
      end
      printf " low_ax=%d count=%d", $parsed, $count
      if $parsed != 1
        set $low_ax_failures = $low_ax_failures + 1
      end
      set $field_index = 0
      while $field_index < $count
        set $field = (unsigned char *)((char *(*)(void *, int))0x10016c30)($parser, $field_index)
        printf " f%d=", $field_index
        if $field == 0
          printf "NULL"
        else
          set $byte_index = 0
          while $byte_index < 32 && *(unsigned char *)($field + $byte_index) != 0
            printf "%02x", *(unsigned char *)($field + $byte_index)
            set $byte_index = $byte_index + 1
          end
        end
        set $field_index = $field_index + 1
      end
      printf "\n"
      set $total_cases = $total_cases + 1
      set $code = $code + 1
    end
    set $len = $len + 1
  end
  printf "CSV_QUOTE_DOMAIN7 cases=%d low_ax_failures=%d\n", $total_cases, $low_ax_failures
  call ((void (*)(void *))0x10016bf0)($parser)
  continue
end

continue
