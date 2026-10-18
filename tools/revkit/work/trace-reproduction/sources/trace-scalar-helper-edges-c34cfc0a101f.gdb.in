set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $runtime_state = *(unsigned int *)0x100a0460
  set $speaker_state = *(unsigned int *)(0x100a0464 + 4)
  set $sound_card_initial = *(int *)0x1007d6e0
  set $index = 0
  while $index < 5
    if $index == 0
      set $value = -2147483648
    else
      if $index == 1
        set $value = -1
      else
        if $index == 2
          set $value = 0
        else
          if $index == 3
            set $value = 1
          else
            set $value = 2147483647
          end
        end
      end
    end
    set $previous_sound_card = ((int (*)(int))0x10026b20)($value)
    printf "SOUND_CARD_EDGE input=%d previous=%d stored=%d\n", $value, $previous_sound_card, *(int *)0x1007d6e0
    set $index = $index + 1
  end
  set $previous_sound_card = ((int (*)(int))0x10026b20)($sound_card_initial)
  printf "SOUND_CARD_RESTORED value=%d previous=%d\n", *(int *)0x1007d6e0, $previous_sound_card
  set $emphasis_initial = *(int *)($speaker_state + 0x4cfc)
  set $paren_initial = *(int *)($runtime_state + 0x2041c)
  set $reading_initial = *(int *)($runtime_state + 0x20420)
  set $index = 0
  while $index < 9
    if $index == 0
      set $value = -2147483648
    else
      if $index == 1
        set $value = -96
      else
        if $index == 2
          set $value = -95
        else
          if $index == 3
            set $value = -94
          else
            if $index == 4
              set $value = 0
            else
              if $index == 5
                set $value = 94
              else
                if $index == 6
                  set $value = 95
                else
                  if $index == 7
                    set $value = 96
                  else
                    set $value = 2147483647
                  end
                end
              end
            end
          end
        end
      end
    end
    call ((void (*)(int, int))0x10028310)($value, 1)
    printf "EMPHASIS_EDGE input=%d stored=%d\n", $value, *(int *)($speaker_state + 0x4cfc)
    set $index = $index + 1
  end
  call ((void (*)(int, int))0x10028310)($emphasis_initial, 1)
  set $index = 0
  while $index < 5
    if $index == 0
      set $value = -2147483648
    else
      if $index == 1
        set $value = -1
      else
        if $index == 2
          set $value = 0
        else
          if $index == 3
            set $value = 1
          else
            set $value = 2147483647
          end
        end
      end
    end
    call ((void (*)(int))0x10028390)($value)
    printf "PARENTHESIS_EDGE input=%d stored=%d\n", $value, *(int *)($runtime_state + 0x2041c)
    set $index = $index + 1
  end
  call ((void (*)(int))0x10028390)($paren_initial)
  set $index = 0
  while $index < 5
    if $index == 0
      set $value = -2147483648
    else
      if $index == 1
        set $value = -1
      else
        if $index == 2
          set $value = 0
        else
          if $index == 3
            set $value = 1
          else
            set $value = 2147483647
          end
        end
      end
    end
    call ((void (*)(int))0x100283c0)($value)
    printf "READING_RULE_EDGE input=%d stored=%d\n", $value, *(int *)($runtime_state + 0x20420)
    set $index = $index + 1
  end
  call ((void (*)(int))0x100283c0)($reading_initial)
  printf "SCALAR_HELPERS_RESTORED emphasis=%d paren=%d reading=%d\n", *(int *)($speaker_state + 0x4cfc), *(int *)($runtime_state + 0x2041c), *(int *)($runtime_state + 0x20420)
  continue
end

continue
