set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  printf "API_CONTROL entry_pitch=%d entry_speed=%d entry_volume=%d\n", *(int *)($esp + 20), *(int *)($esp + 24), *(int *)($esp + 28)
  set {int}($esp + 28) = 120
  printf "API_CONTROL pitch=%d speed=%d set_volume=%d\n", *(int *)($esp + 20), *(int *)($esp + 24), *(int *)($esp + 28)
  continue
end


set $timeline_calls = 0
set $ring_calls = 0
set $join_calls = 0
set $unit_desc_calls = 0
set $dat_reads = 0
set $upm_reads = 0
set $segment_calls = 0
set $recon_calls = 0
set $pcm_blocks = 0

break *0x1002c220
commands
  silent
  set $timeline_calls = $timeline_calls + 1
  set $timeline_rows = *(unsigned int *)($esp + 4)
  set $timeline_state = *(unsigned int *)($esp + 8)
  set $timeline_start = *(unsigned short *)($esp + 12)
  set $timeline_return = *(unsigned int *)$esp
  printf "TIMELINE_ENTRY n=%u start=%u input_count=%u output=%#x state=%#x\n", $timeline_calls, $timeline_start, *(unsigned int *)($timeline_state + 0xec624), $timeline_rows, $timeline_state
  tbreak *$timeline_return
  commands
    silent
    set $timeline_count = *(unsigned short *)($timeline_rows + 0xdf70)
    printf "TIMELINE_RETURN n=%u records=%u\n", $timeline_calls, $timeline_count
    if $timeline_count > 0 && $timeline_count < 1100
      dump binary memory /work/stage8/timeline-records.bin $timeline_rows $timeline_rows+($timeline_count*0x34)
    end
    continue
  end
  continue
end

break *0x1002c530
commands
  silent
  set $ring_calls = $ring_calls + 1
  set $ring_state = *(unsigned int *)($esp + 4)
  set $ring_rows = *(unsigned int *)($esp + 8)
  printf "RING_BUILD n=%u records=%u state=%#x rows=%#x ring_cursor=%u\n", $ring_calls, *(unsigned short *)($ring_rows + 0xdf70), $ring_state, $ring_rows, *(unsigned int *)($ring_state + 0x47780)
  continue
end

break *0x1001b200
commands
  silent
  printf "SELECTED_UNIT id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

break *0x1002c120
commands
  silent
  set $unit_desc_calls = $unit_desc_calls + 1
  if $unit_desc_calls <= 6
    set $desc = *(unsigned int *)($esp + 4)
    set $row = *(unsigned int *)($esp + 8)
    set $desc_return = *(unsigned int *)$esp
    tbreak *$desc_return
    commands
      silent
      printf "UNIT_DESC n=%u mode=%u bank=%u dat_offset=%u dat_samples=%u upm_offset=%u upm_periods=%u edge_a=%u edge_b=%u\n", $unit_desc_calls, *(unsigned char *)($desc + 0x26), *(unsigned short *)($desc + 0x12), *(unsigned int *)($desc + 0x14), *(unsigned int *)($desc + 0x0c), *(unsigned int *)($desc + 0x1c), *(unsigned char *)($desc + 0x20), *(unsigned short *)($desc + 0x22), *(unsigned short *)($desc + 0x24)
      continue
    end
    continue
  else
    continue
  end
end

break *0x1002c8b0
commands
  silent
  set $dat_reads = $dat_reads + 1
  if $dat_reads <= 6
    set $dat_descriptor = *(unsigned int *)($esp + 8)
    set $dat_row = *(unsigned int *)($dat_descriptor + 0x18)
    printf "DAT_DECODE n=%u bank=%u offset=%u coded_bytes=%u samples=%u mode=%u\n", $dat_reads, *(unsigned short *)($dat_descriptor + 0x12), *(unsigned int *)($dat_descriptor + 0x14), *(unsigned short *)($dat_row + 0x10), *(unsigned int *)($dat_descriptor + 0x0c), *(unsigned char *)($dat_descriptor + 0x26)
  end
  continue
end

break *0x1002bc60
commands
  silent
  set $segment_calls = $segment_calls + 1
  set $segment_output = *(unsigned int *)($esp + 4)
  set $segment_vector = *(unsigned int *)($esp + 8)
  set $segment_unit = *(unsigned int *)($esp + 12)
  set $segment_return = *(unsigned int *)$esp
  if $segment_calls <= 8
    printf "SEGMENT_BUILD n=%u vector=%#x unit=%#x periods=%u pitch_field=%u speed_field=%u\n", $segment_calls, $segment_vector, $segment_unit, *(unsigned char *)($segment_unit + 8), *(unsigned int *)($segment_unit + 0x0c), *(unsigned int *)($segment_unit + 0x10)
    x/32hd $segment_vector
    tbreak *$segment_return
    commands
      silent
      printf "SEGMENT_RETURN n=%u records=%u\n", $segment_calls, $eax
      if $segment_calls == 1 && $eax > 0 && $eax < 256
        dump binary memory /work/stage8/map-segments.bin $segment_output $segment_output+($eax*20)
        x/40wx $segment_output
      end
      continue
    end
  end
  continue
end

break *0x1002bbd0
commands
  silent
  set $upm_reads = $upm_reads + 1
  if $upm_reads <= 8
    set $upm_descriptor = *(unsigned int *)($esp + 8)
    printf "UPM_READ n=%u bank=%u offset=%u periods=%u sample_count=%u\n", $upm_reads, *(unsigned short *)($upm_descriptor + 0x12), *(unsigned int *)($upm_descriptor + 0x1c), *(unsigned char *)($upm_descriptor + 0x20), *(unsigned int *)($upm_descriptor + 0x0c)
  end
  continue
end

break *0x1002d230
commands
  silent
  set $join_calls = $join_calls + 1
  set $join_state = *(unsigned int *)($esp + 4)
  set $join_ctx = *(unsigned int *)($esp + 8)
  set $join_index = *(unsigned short *)($esp + 12)
  set $join_return = *(unsigned int *)$esp
  tbreak *$join_return
  commands
    silent
    printf "JOIN_PREPARE n=%u state=%#x model=%#x unit_index=%u result=%#x\n", $join_calls, $join_state, $join_ctx, $join_index, $eax
    if $eax != 0 && $join_calls <= 12
      x/13wx $eax
    end
    continue
  end
  continue
end

break *0x1002afb0
commands
  silent
  set $recon_calls = $recon_calls + 1
  set $recon_context = *(unsigned int *)($esp + 4)
  set $recon_index = *(unsigned short *)($esp + 8)
  set $cursor_before = *(unsigned int *)($recon_context + 0x30)
  set $recon_return = *(unsigned int *)$esp
  tbreak *$recon_return
  commands
    silent
    printf "RECON n=%u index=%u before=%u after=%u samples=%u\n", $recon_calls, $recon_index, $cursor_before, *(unsigned int *)($recon_context + 0x30), (*(unsigned int *)($recon_context + 0x30)-$cursor_before)/2
    continue
  end
  continue
end

break *0x10026870
commands
  silent
  set $pcm_blocks = $pcm_blocks + 1
  set $pcm_ctx = *(unsigned int *)($esp + 8)
  set $pcm_out = *(unsigned int *)($esp + 12)
  set $pcm_return = *(unsigned int *)$esp
  tbreak *$pcm_return
  commands
    silent
    dump binary memory /work/stage8/map-block.pcm $pcm_out $pcm_out+$eax
    printf "PCM_BLOCK n=%u bytes=%u context_bytes=%u\n", $pcm_blocks, $eax, *(unsigned int *)($pcm_ctx + 0x30)
    continue
  end
  continue
end

continue
