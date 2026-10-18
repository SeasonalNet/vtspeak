set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $values = (int *)malloc(16)
  set $mask = 0
  while $mask < 16
    set $values[0] = 0x11111111
    set $values[1] = 0x22222222
    set $values[2] = 0x33333333
    set $values[3] = 0x44444444
    set $out0 = (int *)0
    set $out1 = (int *)0
    set $out2 = (int *)0
    set $out3 = (int *)0
    if ($mask & 1) != 0
      set $out0 = $values
    end
    if ($mask & 2) != 0
      set $out1 = $values + 1
    end
    if ($mask & 4) != 0
      set $out2 = $values + 2
    end
    if ($mask & 8) != 0
      set $out3 = $values + 3
    end
    set $getter_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($out0, $out1, $out2, $out3, 1)
    printf "PSV_NULLMASK mask=%d ret=%d fields=%08x,%08x,%08x,%08x\n", $mask, $getter_result, $values[0], $values[1], $values[2], $values[3]
    set $mask = $mask + 1
  end

  set $out0 = (int *)malloc(4)
  set $mask = 0
  while $mask < 2
    set *$out0 = 0x55667788
    if $mask == 0
      set $comma_pointer = (int *)0
    else
      set $comma_pointer = $out0
    end
    set $getter_result = ((int (*)(int *, int))0x100281f0)($comma_pointer, 1)
    printf "COMMA_NULLMASK mask=%d ret=%d value=%08x\n", $mask, $getter_result, *$out0
    set $mask = $mask + 1
  end

  set $i = 0
  while $i < 10
    if $i == 0
      set $slot = -2147483648
    else
      if $i == 1
        set $slot = -1
      else
        if $i < 8
          set $slot = $i - 2
        else
          if $i == 8
            set $slot = 6
          else
            set $slot = 2147483647
          end
        end
      end
    end
    set $values[0] = 0x11111111
    set $values[1] = 0x22222222
    set $values[2] = 0x33333333
    set $values[3] = 0x44444444
    set $getter_result = ((int (*)(int *, int *, int *, int *, int))0x100280c0)($values, $values + 1, $values + 2, $values + 3, $slot)
    printf "PSV_SLOT selector=%d ret=%d fields=%d,%d,%d,%d\n", $slot, $getter_result, $values[0], $values[1], $values[2], $values[3]
    set *$out0 = 0x55667788
    set $comma_result = ((int (*)(int *, int))0x100281f0)($out0, $slot)
    printf "COMMA_SLOT selector=%d ret=%d value=%d\n", $slot, $comma_result, *$out0
    set $i = $i + 1
  end
  detach
  quit
end

continue
