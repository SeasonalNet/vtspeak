set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands
  silent
  disable 1
  set $history_initial = *(unsigned char *)0x100a045d
  set $history_load_flag = *(unsigned char *)0x100a747c
  set $history_value = 0
  set $history_mismatches = 0
  while $history_value < 256
    call ((void (*)(unsigned char))0x100283f0)((unsigned char)$history_value)
    set $history_stored = *(unsigned char *)0x100a045d
    if $history_value == 1
      set $history_expected = 1
    else
      set $history_expected = 0
    end
    if $history_stored != $history_expected
      set $history_mismatches = $history_mismatches + 1
    end
    printf "UNIT_HISTORY_BYTE input=%u stored=%u expected=%u load_flag=%u\n", $history_value, $history_stored, $history_expected, *(unsigned char *)0x100a747c
    set $history_value = $history_value + 1
  end
  call ((void (*)(unsigned char))0x100283f0)($history_initial)
  printf "UNIT_HISTORY_MATRIX calls=256 mismatches=%d initial=%u restored=%u load_flag=%u\n", $history_mismatches, $history_initial, *(unsigned char *)0x100a045d, $history_load_flag
  continue
end

break *0x1001da50
commands
  silent
  disable 2
  printf "UNIT_HISTORY_AFTER_LOAD stored=%u load_flag=%u\n", *(unsigned char *)0x100a045d, *(unsigned char *)0x100a747c
  continue
end

continue
