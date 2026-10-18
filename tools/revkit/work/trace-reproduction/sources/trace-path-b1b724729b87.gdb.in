set pagination off
set confirm off
set $n_reader = 0
set $n_segment = 0
set $n_reconstruct = 0
set $n_rebuild = 0
target remote localhost:23456

break *0x1002bbd0
commands
  silent
  set $n_reader = $n_reader + 1
  if $n_reader <= 4
    set $trace_desc = *(unsigned int *)($esp + 8)
    printf "UPM#%u bank=%u offset=%u count=%u half_a=%u half_b=%u caller=%#x\\n", $n_reader, *(unsigned short *)($trace_desc + 0x12), *(unsigned int *)($trace_desc + 0x1c), *(unsigned char *)($trace_desc + 0x20), *(unsigned short *)($trace_desc + 0x22), *(unsigned short *)($trace_desc + 0x24), *(unsigned int *)$esp
  end
  continue
end

break *0x1002bc60
commands
  silent
  set $n_segment = $n_segment + 1
  if $n_segment <= 4
    set $trace_unit = *(unsigned int *)($esp + 12)
    printf "SEGMENT#%u periods=%u field0=%u field1=%u caller=%#x\\n", $n_segment, *(unsigned char *)($trace_unit + 8), *(unsigned int *)$trace_unit, *(unsigned int *)($trace_unit + 4), *(unsigned int *)$esp
  end
  continue
end

break *0x1002afb0
commands
  silent
  set $n_reconstruct = $n_reconstruct + 1
  if $n_reconstruct <= 4
    printf "RECON#%u unit=%u caller=%#x\\n", $n_reconstruct, *(unsigned short *)($esp + 8), *(unsigned int *)$esp
  end
  continue
end

break *0x1002bd90
commands
  silent
  set $n_rebuild = $n_rebuild + 1
  if $n_rebuild <= 4
    printf "REBUILD#%u caller=%#x\\n", $n_rebuild, *(unsigned int *)$esp
  end
  continue
end

continue
