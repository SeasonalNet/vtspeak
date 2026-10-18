set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hi-2006-scalar-array-state/legacy-scalar-array-state-gdb.log
set logging overwrite on
set logging enabled on

break *0x1001c860
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if $unit == 272820 || $unit == 65331 || $unit == 272821 || $unit == 280449
    set $state = *(unsigned int *)($esp + 24)
    printf "LEGACY_SCALAR_STATE unit=%u state=%#x\n", $unit, $state
    printf "LEGACY_SCALAR_STATE pointers="
    x/6wx ($state + 0x48)
    printf "LEGACY_SCALAR_STATE values="
    printf "%u,%u,%u,%u,%u,%u\n", *(unsigned char *)(*(unsigned int *)($state + 0x48) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x4c) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x50) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x54) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x58) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x5c) + $unit)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_SCALAR_STATE_TRACE_READY\n"
  continue
end

continue
