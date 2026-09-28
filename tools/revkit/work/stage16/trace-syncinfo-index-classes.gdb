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
  set $state_count = *(short *)($state + 2)
  printf "SYNC_CLASS_SET count=%d state_count=%d state_start=%u\n", $count, $state_count, *(unsigned int *)($state + 0x64c)
  set $i = 0
  while ($i < $count) && ($i < 1100)
    set $record = $records + $i * 0x34
    set $class = *(unsigned char *)($record + 0x32)
    if $class < 0x60
      printf "SYNC_CLASS_RECORD i=%d kind=%u ordinal=%d index=%u side=%u class=%u selector=%u\n", $i, *(unsigned char *)($record + 0x27), *(short *)($record + 0x10), *(unsigned short *)($record + 0x28), *(unsigned short *)($record + 0x2a), $class, *(unsigned char *)(0x1007daa8 + $class)
    else
      printf "SYNC_CLASS_RECORD i=%d kind=%u ordinal=%d index=%u side=%u class=%u selector=OUT_OF_TABLE\n", $i, *(unsigned char *)($record + 0x27), *(short *)($record + 0x10), *(unsigned short *)($record + 0x28), *(unsigned short *)($record + 0x2a), $class
    end
    set $i = $i + 1
  end
  kill
  quit
end

continue
