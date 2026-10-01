set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x100232ea
condition $bpnum (*(unsigned int *)($esi + 8) == 16223 && *(int *)($ebp + 8) == 4)
commands
  silent
  set $signature = *(unsigned int *)($ebp - 0x1c) + 16223 * 7
  printf "BRIDGET_SPAN_PRODUCER position=4 id=16223 sig=%02x,%02x,%02x,%02x,%02x,%02x,%02x fields_c_e_10=%d,%d,%d left_sum=%d right_sum=%d left_count=%d right_count=%d computed_weight=%d\n", *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), *(short *)($esi + 0xc), *(short *)($esi + 0xe), *(short *)($esi + 0x10), *(short *)($ebp - 0x14), *(short *)($ebp - 0x18), *(int *)($ebp + 0x10), *(int *)($ebp - 4), $eax
  continue
end

hbreak *0x408187
commands
  silent
  printf "BRIDGET_SPAN_TRACE_READY\n"
  continue
end

continue
