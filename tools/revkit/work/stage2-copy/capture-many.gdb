set pagination off
set confirm off
target remote localhost:23456
break *0x10001b30
continue
printf "CAPTURE 00\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-00.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-00.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 01\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-01.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-01.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 02\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-02.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-02.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 03\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-03.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-03.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 04\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-04.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-04.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 05\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-05.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-05.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 06\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-06.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-06.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 07\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-07.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-07.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 08\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-08.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-08.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 09\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-09.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-09.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 10\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-10.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-10.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 11\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-11.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-11.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 12\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-12.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-12.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 13\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-13.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-13.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 14\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-14.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-14.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 15\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-15.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-15.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 16\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-16.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-16.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 17\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-17.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-17.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 18\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-18.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-18.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 19\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-19.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-19.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 20\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-20.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-20.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 21\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-21.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-21.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 22\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-22.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-22.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 23\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-23.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-23.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 24\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-24.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-24.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 25\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-25.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-25.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 26\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-26.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-26.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 27\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-27.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-27.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 28\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-28.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-28.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 29\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-29.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-29.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 30\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-30.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-30.pcm $output $output + $decoded_bytes
continue
printf "CAPTURE 31\n"
set $input = *(unsigned int *)($esp + 4)
set $output = *(unsigned int *)($esp + 8)
set $coded_size = *(unsigned int *)($esp + 12)
set $output_size = *(unsigned int *)($esp + 16)
printf "input=%#x output=%#x coded_bytes=%u count_ptr=%#x\n", $input, $output, $coded_size, $output_size
dump binary memory /work/stage2-copy/candidate-31.dat $input $input + $coded_size
finish
set $decoded_bytes = *(unsigned int *)$output_size
printf "decoded=%u\n", $decoded_bytes
dump binary memory /work/stage2-copy/candidate-31.pcm $output $output + $decoded_bytes
detach
quit
