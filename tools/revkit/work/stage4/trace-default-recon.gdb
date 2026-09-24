set pagination off
set confirm off
target remote localhost:23456
break *0x1002afb0
continue
set $synthesis_context = *(unsigned int *)($esp + 4)
set $synthesis_unit = *(unsigned short *)($esp + 8)
set $sample_cursor_before = *(unsigned int *)($synthesis_context + 0x30)
set $synthesis_return = *(unsigned int *)$esp
printf "DEFAULT RECON entry unit=%u cursor=%u return=%#x\n", $synthesis_unit, $sample_cursor_before, $synthesis_return
tbreak *$synthesis_return
continue
printf "DEFAULT RECON return eax=%u cursor=%u samples_added=%u\n", $eax, *(unsigned int *)($synthesis_context + 0x30), (*(unsigned int *)($synthesis_context + 0x30)-$sample_cursor_before)/2
delete 1
continue
