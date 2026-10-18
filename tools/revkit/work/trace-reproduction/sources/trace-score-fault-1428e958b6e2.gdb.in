set pagination off
set confirm off
set debuginfod enabled off

break *0x10018f17
commands
  silent
  printf "candidate-score lookup current_metric=%u candidate_metric=%u row=%#x distance=%g\n", *(unsigned int *)($ebp-0x10), *(unsigned int *)($ebp-0x4c), $ecx, *(float *)($ecx+$eax*4)
  printf "score context pair-table=%#x feature-context=%#x\n", *(unsigned int *)($ebx+0xc), $ebx
  x/24wx $ebp-0x70
  x/8wx $ecx-16
  bt
  kill
end

continue
