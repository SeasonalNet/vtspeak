set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x100232ea
condition $bpnum ((*(unsigned int *)($esi + 8) == 150652 || *(unsigned int *)($esi + 8) == 365789 || *(unsigned int *)($esi + 8) == 382197 || *(unsigned int *)($esi + 8) == 548226) && *(int *)($ebp + 8) == 5)
commands
  silent
  set $unit = *(unsigned int *)($esi + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  set $signature = *(unsigned int *)($ebp - 0x1c) + $unit * 7
  set $previous_mode = *(signed char *)($state + 0xec62c + 4 * 6)
  set $current_mode = *(signed char *)($state + 0xec62c + 5 * 6)
  set $previous_record = $state + 0xae894 + 4 * 0xfc
  set $previous_count = *(short *)($previous_record + 0xf4)
  set $previous_units = $previous_record + 0x7c
  set $contains_current = 0
  set $contains_predecessor = 0
  set $index = 0
  while $index < $previous_count && $index < 5000
    set $previous_unit = *(unsigned int *)($previous_units + $index * 4)
    if $previous_unit == $unit
      set $contains_current = 1
    end
    if $previous_unit == $unit - 1
      set $contains_predecessor = 1
    end
    set $index = $index + 1
  end
  printf "BRIDGET_POSITION5_SPAN position=5 id=%u modes_prev_current=%d,%d prev_count=%d prev_has_id=%d prev_has_id_minus_1=%d sig=%02x,%02x,%02x,%02x,%02x,%02x,%02x fields_c_e_10=%d,%d,%d left_sum=%d right_sum=%d left_count=%d right_count=%d computed_weight=%d\n", $unit, $previous_mode, $current_mode, $previous_count, $contains_current, $contains_predecessor, *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), *(short *)($esi + 0xc), *(short *)($esi + 0xe), *(short *)($esi + 0x10), *(short *)($ebp - 0x14), *(short *)($ebp - 0x18), *(int *)($ebp + 0x10), *(int *)($ebp - 4), $eax
  continue
end

hbreak *0x408187
commands
  silent
  printf "BRIDGET_POSITION5_SPAN_TRACE_READY\n"
  continue
end

continue
