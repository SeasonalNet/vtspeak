set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV stop print nopass
break *0x1001da50
commands 1
  silent
  disable 1
  set $fields = (char **)malloc(8)
  set $a = (unsigned char *)malloc(2)
  set $a[0] = 65
  set $a[1] = 0
  set $empty = (unsigned char *)malloc(1)
  set $empty[0] = 0
  set $fields[0] = (char *)$a
  set $fields[1] = 0
  set $raw = (unsigned char *)malloc(16)
  set $ignored = ((void *(*)(void *, int, unsigned int))memset)($raw, 165, 16)
  set $out = $raw + 1
  set $stack = (unsigned char *)malloc(16384)
  set $fields[0] = (char *)0xffffffff
  printf "CSV_PTR_BEGIN name=bad_first_cap2 count=1 cap=2\n"
  break *0x10016c6d
  commands 2
    silent
    printf "CSV_PTR_RETURN name=bad_first_cap2 low_ax=%d raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x\n", (short)$eax, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15]
    kill
    quit
  end
  catch signal SIGSEGV
  commands 3
    silent
    printf "CSV_PTR_FAULT name=bad_first_cap2 eip=%08x raw=%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x\n", $eip, $raw[0], $raw[1], $raw[2], $raw[3], $raw[4], $raw[5], $raw[6], $raw[7], $raw[8], $raw[9], $raw[10], $raw[11], $raw[12], $raw[13], $raw[14], $raw[15]
    kill
    quit
  end
  set $esp = (unsigned int)$stack + 16128
  set {unsigned int}$esp = 0x10016c6d
  set {unsigned int}($esp + 4) = (unsigned int)$fields
  set {int}($esp + 8) = 1
  set {unsigned int}($esp + 12) = (unsigned int)$out
  set {unsigned int}($esp + 16) = 2
  set $eip = 0x10016c50
  continue
end
continue
