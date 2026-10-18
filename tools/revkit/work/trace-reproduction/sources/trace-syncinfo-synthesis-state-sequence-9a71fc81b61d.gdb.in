set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $calls = 0

break *0x1002c530
commands
  silent
  set $calls = $calls + 1
  set $engine = *(unsigned int *)($esp + 4)
  set $records = *(unsigned int *)($esp + 8)
  set $synth_state = *(unsigned int *)($engine + 0x4c)
  set $record_total = *(short *)($records + 0xdf70)
  set $active_count = *(short *)($synth_state + 2)
  printf "SYNC_STATE_CALL call=%d records=%d active_record_count=%d source_start=%u active_endpoint=%u\n", $calls, $record_total, $active_count, *(unsigned int *)($synth_state + 0x64c), *(unsigned int *)($synth_state + $active_count * 0x3c0 + 0x290)
  set $i = 0
  while ($i < $record_total) && ($i < 1100)
    set $record = $records + $i * 0x34
    printf "SYNC_STATE_RECORD call=%d i=%d kind=%u ordinal=%d group=%u side=%u adjacent_key=%u class=%u selector=%u\n", $calls, $i, *(unsigned char *)($record + 0x27), *(short *)($record + 0x10), *(unsigned short *)($record + 0x28), *(unsigned short *)($record + 0x2a), *(unsigned short *)($record + 0x2e), *(unsigned char *)($record + 0x32), *(unsigned char *)(0x1007daa8 + *(unsigned char *)($record + 0x32))
    set $i = $i + 1
  end
  if $calls >= 20
    kill
    quit
  end
  continue
end

continue
