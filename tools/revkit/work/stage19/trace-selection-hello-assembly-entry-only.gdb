set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-assembly-entry-only-gdb.log
set logging overwrite on
set logging enabled on
set $assembly_calls = 0

break *0x1002bd90
commands
  silent
  set $assembly_calls = $assembly_calls + 1
  printf "ASSEMBLY_ENTRY_ONLY call=%u state=%#x scratch=%#x\n", $assembly_calls, *(unsigned int *)($esp + 4), *(unsigned int *)($esp + 8)
  continue
end

continue
