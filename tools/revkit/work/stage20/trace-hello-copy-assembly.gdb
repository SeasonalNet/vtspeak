#### EXPERIMENTAL ONLY: stopping at FUN_1002bd90 changed the synthesized WAV.
#### Do not use this trace as an output-parity control.
set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-copy-assembly-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0
set $timeline_calls = 0
set $assembly_calls = 0

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

break *0x1002c220
commands
  silent
  set $timeline_calls = $timeline_calls + 1
  printf "TIMELINE_BUILD call=%u start=%u count=%u\n", $timeline_calls, *(unsigned short *)($esp + 12), *(unsigned short *)($esp + 16)
  continue
end

break *0x1002bd90
commands
  silent
  set $assembly_calls = $assembly_calls + 1
  set $state = *(unsigned int *)($esp + 4)
  set $scratch = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  set $row_count = *(short *)($state + 0xfbd68)
  set $rows = $state + 0xeddf8
  set $manager = *(unsigned int *)($state + 0x4c)
  set $row_cursor = *(short *)($manager + 0x47706)
  set $samples_before = *(unsigned int *)($state + 0x30)
  printf "ASSEMBLY_ENTRY call=%u rows=%d row_cursor=%d samples_before=%u scratch=%#x\n", $assembly_calls, $row_count, $row_cursor, $samples_before, $scratch
  set $i = 0
  while $i < $row_count && $i < 16
    printf "ASSEMBLY_TIMELINE_ROW i=%u words:", $i
    x/13wx $rows + $i * 0x34
    set $i = $i + 1
  end
  tbreak *$ret
  commands
    silent
    printf "ASSEMBLY_RETURN call=%u samples_after_bytes=%u row_cursor=%d rows=%d\n", $assembly_calls, *(unsigned int *)($state + 0x30), *(short *)($manager + 0x47706), *(short *)($state + 0xfbd68)
    set $i = 0
    while $i < *(short *)($state + 0xfbd68) && $i < 16
      printf "ASSEMBLY_FINAL_ROW i=%u words:", $i
      x/13wx $rows + $i * 0x34
      set $i = $i + 1
    end
    continue
  end
  continue
end

continue
