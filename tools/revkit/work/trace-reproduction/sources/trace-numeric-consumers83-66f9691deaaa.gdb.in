set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1000e0c0
commands
  silent
  set $tokens = *(unsigned int *)($esp + 4)
  set $code = *(unsigned char *)($esp + 8)
  set $first = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  if $code == 83
    printf "TPP_NUMERIC_WATCH_SETUP code=%u first=%u token=%#x\n", $code, $first, $tokens + $first * 0x94
    tbreak *$ret
    commands
      silent
      rwatch *(unsigned char *)($tokens + $first * 0x94 + 0x25)
      commands
        silent
        set $hits = $hits + 1
        printf "TPP_NUMERIC_FIELD_READ hit=%u eip=%#x value=%u\n", $hits, $eip, *(unsigned char *)($tokens + $first * 0x94 + 0x25)
        x/8i $eip-8
        bt 5
        if $hits >= 16
          disable $bpnum
          printf "TPP_NUMERIC_WATCH_DISABLED\n"
        end
        continue
      end
      continue
    end
  end
  continue
end

set $hits = 0
continue
