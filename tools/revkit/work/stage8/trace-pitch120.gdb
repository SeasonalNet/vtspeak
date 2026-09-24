set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
continue
printf "API_ARGUMENTS before fmt=%u pitch=%d speed=%d volume=%d\n", *(unsigned int *)($esp + 4), *(int *)($esp + 20), *(int *)($esp + 24), *(int *)($esp + 28)
set {int}($esp + 20) = 120
printf "API_ARGUMENTS after fmt=%u pitch=%d speed=%d volume=%d\n", *(unsigned int *)($esp + 4), *(int *)($esp + 20), *(int *)($esp + 24), *(int *)($esp + 28)
delete 1

set $handoffs = 0
set $reconstructions = 0
set $blend_a = 0
set $blend_b = 0
set $upm_reads = 0
set $segment_builds = 0
set $blocks = 0
set $wav_updates = 0
set $stream_writes = 0

break *0x1001b200
commands
  silent
  set $handoffs = $handoffs + 1
  printf "SELECTED_UNIT n=%u id=%u\n", $handoffs, *(unsigned int *)($esp + 8)
  continue
end

break *0x1002bbd0
commands
  silent
  set $upm_reads = $upm_reads + 1
  set $unit = *(unsigned int *)($esp + 8)
  printf "UPM_READ n=%u bank=%u offset=%u periods=%u dat_samples=%u edges=%u/%u\n", $upm_reads, *(unsigned short *)($unit + 0x12), *(unsigned int *)($unit + 0x1c), *(unsigned char *)($unit + 0x20), *(unsigned int *)($unit + 0x0c), *(unsigned short *)($unit + 0x22), *(unsigned short *)($unit + 0x24)
  continue
end

break *0x1002bc60
commands
  silent
  set $segment_builds = $segment_builds + 1
  set $return = *(unsigned int *)$esp
  set $segment_unit = *(unsigned int *)($esp + 12)
  tbreak *$return
  commands
    silent
    printf "SEGMENTS n=%u unit_first=%u unit_target=%u periods=%u segments=%u\n", $segment_builds, *(unsigned int *)$segment_unit, *(unsigned int *)($segment_unit + 4), *(unsigned char *)($segment_unit + 8), $eax
    continue
  end
  continue
end

break *0x1002afb0
commands
  silent
  set $reconstructions = $reconstructions + 1
  set $context = *(unsigned int *)($esp + 4)
  set $unit_index = *(unsigned short *)($esp + 8)
  set $cursor_before = *(unsigned int *)($context + 0x30)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    printf "RECON n=%u unit_index=%u cursor_before=%u cursor_after=%u samples_added=%u\n", $reconstructions, $unit_index, $cursor_before, *(unsigned int *)($context + 0x30), (*(unsigned int *)($context + 0x30)-$cursor_before)/2
    continue
  end
  continue
end

break *0x1002d010
commands
  silent
  set $blend_a = $blend_a + 1
  if $blend_a <= 8
    printf "JOIN_RESAMPLE_A n=%u source=%#x destination=%#x source_start=%d output_count=%d\n", $blend_a, *(unsigned int *)($esp + 4), *(unsigned int *)($esp + 8), *(short *)($esp + 12), *(short *)($esp + 16)
  end
  continue
end

break *0x1002cdf0
commands
  silent
  set $blend_b = $blend_b + 1
  if $blend_b <= 8
    printf "JOIN_RESAMPLE_B n=%u source=%#x destination=%#x source_start=%d output_count=%d\n", $blend_b, *(unsigned int *)($esp + 4), *(unsigned int *)($esp + 8), *(short *)($esp + 12), *(short *)($esp + 16)
  end
  continue
end

break *0x10026870
commands
  silent
  set $blocks = $blocks + 1
  set $ctx = *(unsigned int *)($esp + 8)
  set $out = *(unsigned int *)($esp + 12)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    if $blocks == 1
      dump binary memory /work/stage8/pitch120-block-000.pcm $out $out+$eax
    else
      dump binary memory /work/stage8/pitch120-block-001.pcm $out $out+$eax
    end
    printf "PCM_BLOCK n=%u returned_bytes=%u context_bytes=%u pcm_buffer=%#x join_resample_a_total=%u join_resample_b_total=%u\n", $blocks, $eax, *(unsigned int *)($ctx + 0x30), $out, $blend_a, $blend_b
    continue
  end
  continue
end

break *0x1001f600
commands
  silent
  set $wav_updates = $wav_updates + 1
  set $ctx = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  printf "WAVE_WRITE_ENTRY n=%u format=%u pcm_bytes_ready=%u stream_bytes_before=%u frames_before=%u\n", $wav_updates, *(unsigned short *)($ctx + 0x11b244), *(unsigned int *)($ctx + 0x30), *(unsigned int *)($ctx + 0xeddf0), *(unsigned int *)($ctx + 0xedde8)
  tbreak *$return
  commands
    silent
    printf "WAVE_WRITE_RETURN n=%u stream_bytes_after=%u pcm_bytes_ready=%u\n", $wav_updates, *(unsigned int *)($ctx + 0xeddf0), *(unsigned int *)($ctx + 0x30)
    continue
  end
  continue
end

break *0x10025500
commands
  silent
  set $stream_writes = $stream_writes + 1
  printf "STREAM_WRITE n=%u offset=%u kind=%u width=%u count=%u source=%#x\n", $stream_writes, *(unsigned int *)($esp + 8), *(unsigned int *)($esp + 12), *(unsigned int *)($esp + 20), *(unsigned int *)($esp + 24), *(unsigned int *)($esp + 16)
  continue
end

continue
