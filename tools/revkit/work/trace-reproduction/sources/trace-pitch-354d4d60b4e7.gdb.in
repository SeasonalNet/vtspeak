set pagination off
set confirm off
target remote localhost:23456

# Change only the documented pitch argument at VT_TextToFile entry.
break *0x1001da50
continue
printf "API before fmt=%u pitch=%d speed=%d volume=%d\\n", *(unsigned int *)($esp + 4), *(int *)($esp + 20), *(int *)($esp + 24), *(int *)($esp + 28)
set {int}($esp + 20) = 120
printf "API after fmt=%u pitch=%d speed=%d volume=%d\\n", *(unsigned int *)($esp + 4), *(int *)($esp + 20), *(int *)($esp + 24), *(int *)($esp + 28)
delete 1

# Capture one UPM vector and the first generated interpolation segments.
break *0x1002bbd0
continue
set $upm_destination = *(unsigned int *)($esp + 4)
set $upm_descriptor = *(unsigned int *)($esp + 8)
set $upm_return = *(unsigned int *)$esp
set $upm_bank = *(unsigned short *)($upm_descriptor + 0x12)
set $upm_offset = *(unsigned int *)($upm_descriptor + 0x1c)
set $upm_count = *(unsigned char *)($upm_descriptor + 0x20)
set $sample_count = *(unsigned int *)($upm_descriptor + 0x0c)
printf "UPM entry bank=%u offset=%u count=%u samples=%u dat_offset=%u edge_a=%u edge_b=%u return=%#x\\n", $upm_bank, $upm_offset, $upm_count, $sample_count, *(unsigned int *)($upm_descriptor + 0x14), *(unsigned short *)($upm_descriptor + 0x22), *(unsigned short *)($upm_descriptor + 0x24), $upm_return
tbreak *$upm_return
continue
dump binary memory /work/stage4/pitch-upm-scaled.bin $upm_destination $upm_destination+($upm_count*2)
printf "UPM scaled vector:"
x/32hd $upm_destination
delete 2

break *0x1002bc60
continue
set $segment_output = *(unsigned int *)($esp + 4)
set $segment_vector = *(unsigned int *)($esp + 8)
set $segment_unit = *(unsigned int *)($esp + 12)
set $segment_periods = *(unsigned char *)($segment_unit + 8)
set $segment_return = *(unsigned int *)$esp
printf "SEGMENT entry output=%#x vector=%#x unit=%#x periods=%u field0=%u field1=%u return=%#x\\n", $segment_output, $segment_vector, $segment_unit, $segment_periods, *(unsigned int *)$segment_unit, *(unsigned int *)($segment_unit + 4), $segment_return
printf "SEGMENT vector:"
x/32hd $segment_vector
tbreak *$segment_return
continue
set $segments_built = $eax
printf "SEGMENT return count=%u\\n", $segments_built
dump binary memory /work/stage4/pitch-segments.bin $segment_output $segment_output+($segments_built*20)
x/40wx $segment_output
delete 4

break *0x1002afb0
continue
set $synthesis_context = *(unsigned int *)($esp + 4)
set $synthesis_unit = *(unsigned short *)($esp + 8)
set $sample_cursor_before = *(unsigned int *)($synthesis_context + 0x30)
set $synthesis_return = *(unsigned int *)$esp
printf "RECON entry unit=%u cursor=%u return=%#x\\n", $synthesis_unit, $sample_cursor_before, $synthesis_return
tbreak *$synthesis_return
continue
printf "RECON return eax=%u cursor=%u\\n", $eax, *(unsigned int *)($synthesis_context + 0x30)
delete 6

# Capture the next transformed unit through the same three boundaries.
break *0x1002bbd0
continue
set $next_upm_destination = *(unsigned int *)($esp + 4)
set $next_upm_descriptor = *(unsigned int *)($esp + 8)
set $next_upm_return = *(unsigned int *)$esp
set $next_upm_count = *(unsigned char *)($next_upm_descriptor + 0x20)
printf "NEXT UPM bank=%u offset=%u count=%u samples=%u edge_a=%u edge_b=%u return=%#x\\n", *(unsigned short *)($next_upm_descriptor + 0x12), *(unsigned int *)($next_upm_descriptor + 0x1c), $next_upm_count, *(unsigned int *)($next_upm_descriptor + 0x0c), *(unsigned short *)($next_upm_descriptor + 0x22), *(unsigned short *)($next_upm_descriptor + 0x24), $next_upm_return
tbreak *$next_upm_return
continue
dump binary memory /work/stage4/next-upm-scaled.bin $next_upm_destination $next_upm_destination+($next_upm_count*2)
printf "NEXT UPM scaled vector:"
x/32hd $next_upm_destination
delete 8

break *0x1002bc60
continue
set $next_segment_output = *(unsigned int *)($esp + 4)
set $next_segment_vector = *(unsigned int *)($esp + 8)
set $next_segment_unit = *(unsigned int *)($esp + 12)
set $next_segment_return = *(unsigned int *)$esp
printf "NEXT SEGMENT vector=%#x unit=%#x first=%u second=%u return=%#x\\n", $next_segment_vector, $next_segment_unit, *(unsigned int *)$next_segment_unit, *(unsigned int *)($next_segment_unit + 4), $next_segment_return
tbreak *$next_segment_return
continue
set $next_segments_built = $eax
printf "NEXT SEGMENT return count=%u\\n", $next_segments_built
dump binary memory /work/stage4/next-segments.bin $next_segment_output $next_segment_output+($next_segments_built*20)
delete 10

break *0x1002afb0
continue
set $next_synthesis_context = *(unsigned int *)($esp + 4)
set $next_synthesis_unit = *(unsigned short *)($esp + 8)
set $next_cursor_before = *(unsigned int *)($next_synthesis_context + 0x30)
set $next_synthesis_return = *(unsigned int *)$esp
printf "NEXT RECON entry unit=%u cursor=%u return=%#x\\n", $next_synthesis_unit, $next_cursor_before, $next_synthesis_return
tbreak *$next_synthesis_return
continue
printf "NEXT RECON return eax=%u cursor=%u samples_added=%u\\n", $eax, *(unsigned int *)($next_synthesis_context + 0x30), (*(unsigned int *)($next_synthesis_context + 0x30)-$next_cursor_before)/2
delete 12
continue
