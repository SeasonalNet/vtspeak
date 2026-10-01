set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hi-2006-volume-setter-gdb.log
set logging overwrite on
set logging enabled on
set $setter_calls = 0

hbreak *0x100219a0
commands
  silent
  set $setter_calls = $setter_calls + 1
  printf "LEGACY_HI_VOLUME_SETTER call=%u pitch=%d speed=%d volume=%d pause=%d speaker=%d\n", $setter_calls, *(int *)($esp + 4), *(int *)($esp + 8), *(int *)($esp + 12), *(int *)($esp + 16), *(int *)($esp + 20)
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_HI_VOLUME_SETTER_TRACE_READY\n"
  continue
end

continue
