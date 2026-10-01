set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $buffer_alloc = (char *)malloc(60032)
  set $buffer = $buffer_alloc + 16
  set $i = 0
  while $i < 60000
    set {unsigned char}($buffer + $i) = 0xa5
    set $i = $i + 1
  end
  set $length = (int *)malloc(4)
  set $sync_out = (unsigned int *)malloc(4)
  set $records_out = (unsigned int *)malloc(4)
  set *$length = 60000
  set *$sync_out = 0x55555555
  set *$records_out = 0x66666666
  call ((void (*)(unsigned char))0x10028360)(0)
  printf "HIGHLIGHT_EX_SETTING flag=0 stored=%u\n", *(unsigned char *)(*(unsigned int *)0x100a0460 + 0x20424)
  set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, (char *)$text, $buffer, $length, 0, 0, 1, $sync_out, $records_out, -1, -1, -1, -1, 0, 0)
  printf "BUFFER_EX_RAW_START selector=0 result=%d length=%d\n", $result, *$length

  set $desc = (char *)*$records_out
  if $desc != 0
    set $record_count = *(int *)$desc
    set $records = *(unsigned int *)($desc + 4)
    printf "BUFFER_EX_RAW_DESCRIPTOR count=%d array=%#x header=%#x,%#x,%#x,%#x\n", $record_count, $records, *(unsigned int *)$desc, *(unsigned int *)($desc + 4), *(unsigned int *)($desc + 8), *(unsigned int *)($desc + 12)
    if $record_count > 0
      dump binary memory /work/stage21/highlight-ex-two-byte-classes-replay-0-0-descriptor.bin $records $records + $record_count * 0x210
    end
    set $i = 0
    while $i < $record_count
      set $row = $records + $i * 0x210
      printf "BUFFER_EX_RAW_ROW i=%d f0=%#x f4=%#x f8=%#x kind=%u\n", $i, *(unsigned int *)$row, *(unsigned int *)($row + 4), *(unsigned int *)($row + 8), *(unsigned char *)($row + 0x20c)
      set $i = $i + 1
    end
  else
    printf "BUFFER_EX_RAW_DESCRIPTOR null\n"
  end

  set $sync = (char *)*$sync_out
  if $sync != 0
    set $sync_rows = *(unsigned int *)$sync
    set $sync_count = *(int *)($sync + 12)
    printf "BUFFER_EX_RAW_SYNCINFO count=%d header=%#x,%#x,%#x,%#x,%#x,%#x,%#x\n", $sync_count, *(unsigned int *)($sync + 16), *(unsigned int *)($sync + 20), *(unsigned int *)($sync + 24), *(unsigned int *)($sync + 28), *(unsigned int *)($sync + 32), *(unsigned int *)($sync + 36), *(unsigned int *)($sync + 40)
    set $i = 0
    while $i < $sync_count
      set $sync_row = $sync_rows + $i * 0x24
      printf "BUFFER_EX_RAW_SYNCROW i=%d nested_count=%u field8=%#x text_start=%#x text_end=%#x field14=%#x field18=%#x field1c=%#x field20=%#x\n", $i, *(unsigned short *)$sync_row, *(unsigned int *)($sync_row + 8), *(unsigned int *)($sync_row + 12), *(unsigned int *)($sync_row + 16), *(unsigned int *)($sync_row + 20), *(unsigned int *)($sync_row + 24), *(unsigned int *)($sync_row + 28), *(unsigned int *)($sync_row + 32)
      set $i = $i + 1
    end
  end
  set *$length = 60000
  set $cancel = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, (char *)$text, $buffer, $length, 2, 0, 1, (unsigned int *)0, (unsigned int *)0, -1, -1, -1, -1, 0, 0)
  printf "BUFFER_EX_RAW_CANCEL result=%d length=%d\n", $cancel, *$length
  continue
end

continue
