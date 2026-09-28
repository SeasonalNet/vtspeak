set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1002c530
commands
  silent
  set $state = *(unsigned int *)($esp + 4)
  set $records = *(unsigned int *)($esp + 8)
  set $count = *(short *)($records + 0xdf70)
  printf "SYNC_NEIGHBOR_SET count=%d state_count=%d\n", $count, *(short *)($state + 2)
  set $i = 0
  while ($i < $count) && ($i < 1100)
    set $record = $records + $i * 0x34
    printf "SYNC_NEIGHBOR_RECORD i=%d kind=%u group=%u side=%u source_guard=%u raw_class=%u prefix=%08x,%08x,%08x\n", $i, *(unsigned char *)($record + 0x27), *(unsigned short *)($record + 0x28), *(unsigned short *)($record + 0x2a), *(unsigned short *)($record + 0x2e), *(unsigned char *)($record + 0x32), *(unsigned int *)$record, *(unsigned int *)($record + 4), *(unsigned int *)($record + 8)
    set $i = $i + 1
  end
  kill
  quit
end

continue
