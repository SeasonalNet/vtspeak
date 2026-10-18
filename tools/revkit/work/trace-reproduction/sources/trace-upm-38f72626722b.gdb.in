set pagination off
set confirm off
target remote localhost:23456

# Unit zero is selected in the reference phrase, but it may use the fast path.
# Stop on the first unit that actually reads and transforms UPM data.
break *0x1002bbd0
continue
set $upm_destination = *(unsigned int *)($esp + 4)
set $upm_descriptor = *(unsigned int *)($esp + 8)
set $upm_return = *(unsigned int *)$esp
set $upm_bank = *(unsigned short *)($upm_descriptor + 0x12)
set $upm_offset = *(unsigned int *)($upm_descriptor + 0x1c)
set $upm_count = *(unsigned char *)($upm_descriptor + 0x20)
printf "UPM entry descriptor=%#x destination=%#x bank=%u offset=%u count=%u half_a=%u half_b=%u return=%#x\\n", $upm_descriptor, $upm_destination, $upm_bank, $upm_offset, $upm_count, *(unsigned short *)($upm_descriptor + 0x22), *(unsigned short *)($upm_descriptor + 0x24), $upm_return
tbreak *$upm_return
continue
dump binary memory /work/stage4/upm-scaled.bin $upm_destination $upm_destination+($upm_count*2)
printf "UPM scaled vector:"
x/32hd $upm_destination
delete 1

# The next call converts the same vector into interpolation segment records.
break *0x1002bc60
continue
set $segment_output = *(unsigned int *)($esp + 4)
set $segment_vector = *(unsigned int *)($esp + 8)
set $segment_unit = *(unsigned int *)($esp + 12)
set $segment_periods = *(unsigned char *)($segment_unit + 8)
set $segment_return = *(unsigned int *)$esp
printf "SEGMENT entry output=%#x vector=%#x unit=%#x periods=%u first=%u second=%u return=%#x\\n", $segment_output, $segment_vector, $segment_unit, $segment_periods, *(unsigned int *)$segment_unit, *(unsigned int *)($segment_unit + 4), $segment_return
printf "SEGMENT vector:"
x/32hd $segment_vector
tbreak *$segment_return
continue
set $segments_built = $eax
printf "SEGMENT return count=%u\\n", $segments_built
dump binary memory /work/stage4/upm-segments.bin $segment_output $segment_output+($segments_built*20)
x/40wx $segment_output
delete 3

# Observe the reconstruction consumer's output cursor before and after.
break *0x1002afb0
continue
set $synthesis_context = *(unsigned int *)($esp + 4)
set $synthesis_unit = *(unsigned short *)($esp + 8)
set $sample_cursor_before = *(unsigned int *)($synthesis_context + 0x30)
set $sample_cursor_next_before = *(unsigned int *)($synthesis_context + 0x34)
set $synthesis_return = *(unsigned int *)$esp
printf "RECON entry context=%#x unit=%u cursor=%u next=%u return=%#x\\n", $synthesis_context, $synthesis_unit, $sample_cursor_before, $sample_cursor_next_before, $synthesis_return
tbreak *$synthesis_return
continue
printf "RECON return eax=%u cursor=%u next=%u\\n", $eax, *(unsigned int *)($synthesis_context + 0x30), *(unsigned int *)($synthesis_context + 0x34)
delete 5
continue
