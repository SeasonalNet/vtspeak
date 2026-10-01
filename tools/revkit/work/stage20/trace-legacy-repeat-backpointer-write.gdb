set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-repeat-2006-backpointer/legacy-backpointer-write-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001ec44
commands
  silent
  set $position = *(int *)($esp + 0x2c)
  set $candidate = *(unsigned int *)($esp + 0x3c)
  if ($position == 1) && ($candidate == 272823)
    set $state = *(unsigned int *)($esp + 0xa0)
    set $pred_index = $eax
    set $pred_array = $state + 0xbb154 + 0x7c
    set $pred_id = *(unsigned int *)($pred_array + $pred_index * 4)
    printf "LEGACY_BACKPOINTER_WRITE position=%d candidate=%u selected_index=%u selected_previous=%u first=%u second=%u\n", $position, $candidate, $pred_index, $pred_id, *(unsigned int *)$pred_array, *(unsigned int *)($pred_array + 4)
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_BACKPOINTER_TRACE_READY\n"
  continue
end

continue
