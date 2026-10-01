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
  set $i = 0
  while $i < 16
    set {unsigned char}($buffer_alloc + $i) = 0xa5
    set {unsigned char}($buffer + 60000 + $i) = 0xa5
    set $i = $i + 1
  end
  set $length = (int *)malloc(4)
  set $extra_a = (unsigned int *)malloc(4)
  set $extra_b = (unsigned int *)malloc(4)
  set *$length = 60000
  set *$extra_a = 0x55555555
  set *$extra_b = 0x66666666

  printf "BUFFER_EX_RECORD_INPUT=%s\n", $text
  call ((void (*)(unsigned char))0x10028360)(1)
  printf "HIGHLIGHT_EX_SETTING flag=%u stored=%u\n", 1, *(unsigned char *)(*(unsigned int *)0x100a0460 + 0x20424)
  set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, (char *)$text, $buffer, $length, 0, 0, 1, $extra_a, $extra_b, -1, -1, -1, -1, 0, 0)
  printf "BUFFER_EX_RECORD_START selector=0 result=%d length=%d extra_a=%#x extra_b=%#x guards=%02x,%02x\n", $result, *$length, *$extra_a, *$extra_b, *(unsigned char *)($buffer - 1), *(unsigned char *)($buffer + 60000)

  set $desc = (char *)*$extra_b
  if $desc != 0
    set $record_count = *(int *)$desc
    set $records = *(unsigned int *)($desc + 4)
    printf "BUFFER_EX_RECORD_DESCRIPTOR selector=0 pointer=%#x count=%d array=%#x header=%#x,%#x,%#x,%#x\n", $desc, $record_count, $records, *(unsigned int *)$desc, *(unsigned int *)($desc + 4), *(unsigned int *)($desc + 8), *(unsigned int *)($desc + 12)
    set $i = 0
    while ($i < $record_count) && ($i < 32)
      set $row = $records + $i * 0x210
      printf "BUFFER_EX_RECORD_ROW selector=0 i=%d f0=%#x f4=%#x f8=%#x kind=%u tail=%02x,%02x,%02x,%02x payload_first32=\n", $i, *(unsigned int *)$row, *(unsigned int *)($row + 4), *(unsigned int *)($row + 8), *(unsigned char *)($row + 0x20c), *(unsigned char *)($row + 0x20c - 3), *(unsigned char *)($row + 0x20c - 2), *(unsigned char *)($row + 0x20c - 1), *(unsigned char *)($row + 0x20c)
      x/32bx $row + 0xc
      set $i = $i + 1
    end
  else
    printf "BUFFER_EX_RECORD_DESCRIPTOR selector=0 null\n"
  end

  set $sync = (char *)*$extra_a
  if $sync != 0
    set $sync_rows = *(unsigned int *)$sync
    set $sync_count = *(int *)($sync + 12)
    printf "BUFFER_EX_SYNCINFO selector=0 pointer=%#x rows=%#x count=%d header_fields_4_10=%#x,%#x,%#x,%#x,%#x,%#x,%#x\n", $sync, $sync_rows, $sync_count, *(unsigned int *)($sync + 16), *(unsigned int *)($sync + 20), *(unsigned int *)($sync + 24), *(unsigned int *)($sync + 28), *(unsigned int *)($sync + 32), *(unsigned int *)($sync + 36), *(unsigned int *)($sync + 40)
    set $i = 0
    while ($i < $sync_count) && ($i < 20)
      set $sync_row = $sync_rows + $i * 0x24
      printf "BUFFER_EX_SYNCROW selector=0 i=%d nested_count=%u nested=%#x field8=%#x text_start=%#x text_end=%#x field14=%#x field18=%#x field1c=%#x field20=%#x\n", $i, *(unsigned short *)$sync_row, *(unsigned int *)($sync_row + 4), *(unsigned int *)($sync_row + 8), *(unsigned int *)($sync_row + 12), *(unsigned int *)($sync_row + 16), *(unsigned int *)($sync_row + 20), *(unsigned int *)($sync_row + 24), *(unsigned int *)($sync_row + 28), *(unsigned int *)($sync_row + 32)
      set $i = $i + 1
    end
  end

  set *$length = 60000
  set $cancel = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(0, (char *)$text, $buffer, $length, 2, 0, 1, (unsigned int *)0, (unsigned int *)0, -1, -1, -1, -1, 0, 0)
  printf "BUFFER_EX_RECORD_CANCEL selector=0 result=%d length=%d\n", $cancel, *$length
  continue
end

continue
