set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# At position zero unit 272822 covers the full six-position input. The branch
# at 0x100232ea stores left_half + right_half + 1.
break *0x100232ea
condition $bpnum ((*(unsigned int *)($esi + 8) == 272822) && (*(int *)($ebp + 8) == 0))
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $left_sum = *(short *)($ebp - 0x14)
  set $right_sum = *(short *)($ebp - 0x18)
  set $left_count = *(int *)($ebp + 0x10)
  set $right_count = *(int *)($ebp - 4)
  set $length = *(int *)($state + 0xec624)
  set $signature = *(unsigned char *)(*(unsigned int *)($state + 0x4c) + 272822 * 7 + 6)
  printf "HELLO_WEIGHTED_PRODUCER position=0 length=%d id=272822 total_span=%d left_matches=%d right_matches=%d left_weight_sum=%d right_weight_sum=%d left_half=%d right_half=%d selected_plus10=%d sig6=%02x\n", $length, *(short *)($esi + 0xe), $left_count, $right_count, $left_sum, $right_sum, $left_sum / 2, $right_sum / 2, $eax, $signature
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_WEIGHTED_PRODUCER_TRACE_READY\n"
  continue
end

continue
