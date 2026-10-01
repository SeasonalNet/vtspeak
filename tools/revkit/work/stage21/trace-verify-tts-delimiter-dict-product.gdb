set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $path = ((char *(*)(unsigned int))0x1001d9c0)(41)
  set {char}($path + 0) = 90
  set {char}($path + 1) = 58
  set {char}($path + 2) = 47
  set {char}($path + 3) = 119
  set {char}($path + 4) = 111
  set {char}($path + 5) = 114
  set {char}($path + 6) = 107
  set {char}($path + 7) = 47
  set {char}($path + 8) = 115
  set {char}($path + 9) = 116
  set {char}($path + 10) = 97
  set {char}($path + 11) = 103
  set {char}($path + 12) = 101
  set {char}($path + 13) = 49
  set {char}($path + 14) = 54
  set {char}($path + 15) = 47
  set {char}($path + 16) = 117
  set {char}($path + 17) = 115
  set {char}($path + 18) = 101
  set {char}($path + 19) = 114
  set {char}($path + 20) = 100
  set {char}($path + 21) = 105
  set {char}($path + 22) = 99
  set {char}($path + 23) = 116
  set {char}($path + 24) = 45
  set {char}($path + 25) = 104
  set {char}($path + 26) = 101
  set {char}($path + 27) = 108
  set {char}($path + 28) = 108
  set {char}($path + 29) = 111
  set {char}($path + 30) = 45
  set {char}($path + 31) = 112
  set {char}($path + 32) = 108
  set {char}($path + 33) = 97
  set {char}($path + 34) = 105
  set {char}($path + 35) = 110
  set {char}($path + 36) = 46
  set {char}($path + 37) = 99
  set {char}($path + 38) = 115
  set {char}($path + 39) = 118
  set {char}($path + 40) = 0
  set $gate_original = *(unsigned char *)0x100a7489
  set $load = ((short (*)(int, char *))0x10027960)(27, $path)
  printf "VERIFY_TTS_DELIMITER_DICT_LOAD ax=%d ptr=%#x gate=%u\n", $load, *(unsigned int *)(0x100a647c + 27 * 4), $gate_original
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(16)
  set $state = 0
  while $state < 2
    if $state == 0
      set {unsigned char}0x100a7489 = 0
    else
      set {unsigned char}0x100a7489 = 1
    end
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
        set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, 27, 0)
        printf "VERIFY_TTS_DELIMITER_DICT state=%d length=%d code=%d ax=%u\n", $state, $n, $i, $r & 65535
        set $i = $i + 1
      end
      set $n = $n + 1
    end
    set $state = $state + 1
  end
  set {unsigned char}0x100a7489 = $gate_original
  set $unload = ((short (*)(int))0x10027a80)(27)
  printf "VERIFY_TTS_DELIMITER_DICT_UNLOAD ax=%d restored_gate=%u\n", $unload, *(unsigned char *)0x100a7489
  continue
end

continue
