set pagination off
set confirm off
target remote 127.0.0.1:23456
break *0x10001b30
continue
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "decoder input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/engine-input-first.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoder returned %u PCM bytes (%u samples)\n", $decoded_bytes, $decoded_bytes / 2
dump binary memory /work/stage2-copy/engine-output-first.pcm $output $output + $decoded_bytes
detach
quit
