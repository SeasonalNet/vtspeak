set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# FUN_100230a0 has chosen the left-side half because the right side reaches
# the utterance boundary. This is the direct-store branch at 0x100232d7.
break *0x100232d7
condition $bpnum ((*(unsigned int *)($esi + 8) == 266885) && ((*(int *)($ebp + 8) == 5) || (*(int *)($ebp + 8) == 6)))
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  set $unit = *(unsigned int *)($esi + 8)
  set $left_sum = *(int *)($ebp - 0x14)
  set $right_sum = *(short *)($ebp - 0x18)
  set $left_count = *(int *)($ebp + 0x10)
  set $right_count = *(int *)($ebp - 4)
  set $length = *(int *)($state + 0xec624)
  set $signature = *(unsigned char *)(*(unsigned int *)($state + 0x4c) + $unit * 7 + 6)
  printf "APPLE_WEIGHTED_PRODUCER position=%d length=%d id=%u total_span=%d left_matches=%d right_matches=%d left_weight_sum=%d right_weight_sum=%d left_half=%d right_half=%d selected_plus10=%d sig6=%02x\n", $position, $length, $unit, *(short *)($esi + 0xe), $left_count, $right_count, $left_sum, $right_sum, $left_sum / 2, $right_sum / 2, $ecx, $signature
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_WEIGHTED_PRODUCER_TRACE_READY\n"
  continue
end

continue
