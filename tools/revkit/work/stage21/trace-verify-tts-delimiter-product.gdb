set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(16)
  set $n = 1
  while $n <= 4
    set $limit = 1
    set $j = 0
    while $j < $n
      set $limit = $limit * 4
      set $j = $j + 1
    end
    set $i = 0
    while $i < $limit
      set $v = $i
      set $j = 0
      while $j < $n
        set $d = $v % 4
        if $d == 0
          set {char}($p + $j) = 60
        else
          if $d == 1
            set {char}($p + $j) = 62
          else
            if $d == 2
              set {char}($p + $j) = 47
            else
              set {char}($p + $j) = 65
            end
          end
        end
        set $v = $v / 4
        set $j = $j + 1
      end
      set {char}($p + $n) = 0
      set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
      printf "VERIFY_TTS_DELIMITER_PRODUCT length=%d code=%d ax=%u\n", $n, $i, $r & 65535
      set $i = $i + 1
    end
    set $n = $n + 1
  end
  continue
end

continue
