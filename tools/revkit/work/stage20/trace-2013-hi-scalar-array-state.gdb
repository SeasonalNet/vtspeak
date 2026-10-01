set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-scalar-array-state/adapted-scalar-array-state-gdb.log
set logging overwrite on
set logging enabled on

break *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if $unit == 272820 || $unit == 65331 || $unit == 272821 || $unit == 280449
    set $context = *(unsigned int *)($esp + 12)
    set $state = *(unsigned int *)($esp + 24)
    printf "SCALAR_STATE unit=%u context=%02x,%02x,%02x,%02x,%02x,%02x state=%#x param2=%#x\n", $unit, *(unsigned char *)$context, *(unsigned char *)($context + 1), *(unsigned char *)($context + 2), *(unsigned char *)($context + 3), *(unsigned char *)($context + 4), *(unsigned char *)($context + 5), $state, *(unsigned int *)($esp + 8)
    printf "SCALAR_STATE pointers="
    x/6wx ($state + 0x48)
    printf "SCALAR_STATE values="
    printf "%u,%u,%u,%u,%u,%u\n", *(unsigned char *)(*(unsigned int *)($state + 0x48) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x4c) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x50) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x54) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x58) + $unit), *(unsigned char *)(*(unsigned int *)($state + 0x5c) + $unit)
    printf "SCALAR_STATE direct_columns="
    x/3wx ($state + 0x60)
    bt 8
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "SCALAR_STATE_TRACE_READY\n"
  continue
end

continue
