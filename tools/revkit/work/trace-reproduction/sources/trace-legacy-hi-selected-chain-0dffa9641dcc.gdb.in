set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hi-2006-selected-chain/legacy-selected-chain-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0

hbreak *0x1001ed7f
commands
  silent
  set $selected_calls = $selected_calls + 1
  set $selected_index = *(short *)$edi
  printf "LEGACY_HI_FINAL_CHAIN call=%u selected_index=%d unit_id=%u selected_index_array=%#x\n", $selected_calls, $selected_index, $edx, $edi
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_HI_CHAIN_TRACE_READY\n"
  continue
end

continue
