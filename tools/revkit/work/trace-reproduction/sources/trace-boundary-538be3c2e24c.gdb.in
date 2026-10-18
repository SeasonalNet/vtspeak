set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $timeline_calls = 0
set $pcm_blocks = 0

break *0x1002c220
commands
  silent
  set $timeline_calls = $timeline_calls + 1
  set $rows = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    printf "TIMELINE n=%u rows=%u\n", $timeline_calls, *(unsigned short *)($rows + 0xdf70)
    continue
  end
  continue
end

break *0x10026870
commands
  silent
  set $pcm_out = *(unsigned int *)($esp + 12)
  set $pcm_return = *(unsigned int *)$esp
  tbreak *$pcm_return
  commands
    silent
    set $pcm_blocks = $pcm_blocks + 1
    printf "PCM_BLOCK n=%u bytes=%u\n", $pcm_blocks, $eax
    append binary memory /work/corpus-parity/current-blocks.pcm $pcm_out $pcm_out+$eax
    continue
  end
  continue
end

continue
