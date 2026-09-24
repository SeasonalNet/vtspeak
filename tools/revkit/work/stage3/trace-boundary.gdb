set pagination off
set confirm off
target remote localhost:23456
break *0x10026870
continue
set $block_context = *(unsigned int *)($esp + 8)
set $block_output = *(unsigned int *)($esp + 12)
set $block_return = *(unsigned int *)$esp
printf "BLOCK 0 entry context=%#x output=%#x return=%#x\n", $block_context, $block_output, $block_return
tbreak *$block_return
continue
set $block_length = $eax
printf "BLOCK 0 return eax=%u context_pcm_bytes=%u\n", $block_length, *(unsigned int *)($block_context + 0x30)
dump binary memory /work/stage3/blocks/block-000.pcm $block_output $block_output+$block_length
continue
set $block_context = *(unsigned int *)($esp + 8)
set $block_output = *(unsigned int *)($esp + 12)
set $block_return = *(unsigned int *)$esp
printf "BLOCK 1 entry context=%#x output=%#x return=%#x\n", $block_context, $block_output, $block_return
tbreak *$block_return
continue
set $block_length = $eax
printf "BLOCK 1 return eax=%u context_pcm_bytes=%u\n", $block_length, *(unsigned int *)($block_context + 0x30)
dump binary memory /work/stage3/blocks/block-001.pcm $block_output $block_output+$block_length
disable 1
continue
