set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $runtime_state = *(unsigned int *)0x100a0460
  set $highlight_initial = *(unsigned char *)($runtime_state + 0x20424)
  set $highlight_value = 0
  set $highlight_mismatches = 0
  while $highlight_value < 256
    call ((void (*)(unsigned char))0x10028360)((unsigned char)$highlight_value)
    set $highlight_stored = *(unsigned char *)($runtime_state + 0x20424)
    if $highlight_value == 0
      set $highlight_expected = 0
    else
      set $highlight_expected = 1
    end
    if $highlight_stored != $highlight_expected
      set $highlight_mismatches = $highlight_mismatches + 1
    end
    printf "HIGHLIGHT_BYTE input=%u stored=%u expected=%u\n", $highlight_value, $highlight_stored, $highlight_expected
    set $highlight_value = $highlight_value + 1
  end
  call ((void (*)(unsigned char))0x10028360)($highlight_initial)
  printf "HIGHLIGHT_MATRIX calls=256 mismatches=%d initial=%u restored=%u\n", $highlight_mismatches, $highlight_initial, *(unsigned char *)($runtime_state + 0x20424)
  continue
end

continue
