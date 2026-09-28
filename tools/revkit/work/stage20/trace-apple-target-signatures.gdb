set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-target-signatures-gdb.log
set logging overwrite on
set logging enabled on
hbreak *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $row = $state + ($slot * 3 + 0x76314) * 2
  set $bank = *(unsigned char *)$row
  set $leaf = *(unsigned char *)($row + 2)
  set $voice = *(unsigned int *)($state + 0x4c)
  set $signature = $voice + $bank * 0x3c0 + 0x6e2 + $leaf * 7
  printf "APPLE_TARGET_SIGNATURE slot=%u bank=%u leaf=%u:", $slot, $bank, $leaf
  x/7ub $signature
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_NEW_QUERY_TRACE_READY\n"
  continue
end

continue
