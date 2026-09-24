set pagination off
set confirm off
set logging file /work/stage5/probes/stage6/abbreviations/tree-lookups-gdb.log
set logging overwrite on
set logging enabled on
break *0x10001670
commands
  silent
  set $tree = *(unsigned int *)($esp + 4)
  set $features = *(unsigned int *)($esp + 8)
  set $return_address = *(unsigned int *)$esp
  set $nodes = *(short *)($tree + 24)
  set $output_width = *(unsigned char *)($tree + 26)
  printf "SCALAR_ENTRY tree=%#x nodes=%d output_width=%u features=%#x return=%#x\n", $tree, $nodes, $output_width, $features, $return_address
  x/16hd $features
  tbreak *$return_address
  commands
    silent
    printf "SCALAR_RETURN value=%d\n", (short)$eax
    continue
  end
  continue
end

break *0x100016a0
commands
  silent
  set $tree = *(unsigned int *)($esp + 4)
  set $features = *(unsigned int *)($esp + 8)
  set $output = *(unsigned int *)($esp + 12)
  set $return_address = *(unsigned int *)$esp
  set $nodes = *(short *)($tree + 24)
  set $output_width = *(unsigned char *)($tree + 26)
  printf "VECTOR_ENTRY tree=%#x nodes=%d output_width=%u features=%#x output=%#x return=%#x\n", $tree, $nodes, $output_width, $features, $output, $return_address
  x/16hd $features
  tbreak *$return_address
  commands
    silent
    printf "VECTOR_RETURN values:"
    x/12hd $output
    continue
  end
  continue
end

break *0x1001c990
commands
  silent
  set $ctx = *(unsigned int *)($esp + 4)
  set $input = *(unsigned int *)($esp + 8)
  set $return_address = *(unsigned int *)$esp
  printf "TEXT_PREPROCESS_ENTRY input=%#x\n", $input
  x/s $input
  tbreak *$return_address
  commands
    silent
    printf "TEXT_PREPROCESS_RETURN value=%#x\n", $eax
    x/s $eax
    set $records = *(unsigned int *)($ctx + 0x122448)
    set $record_count = *(int *)$records
    set $record_data = *(unsigned int *)($records + 4)
    printf "TEXT_RECORDS count=%d base=%#x stride=0x210\n", $record_count, $record_data
    set $ri = 0
    while $ri < $record_count && $ri < 8
      printf "TEXT_RECORD_RAW index=%d bytes:\n", $ri
      x/24bx ($record_data + $ri * 0x210)
      set $ri = $ri + 1
    end
    continue
  end
  continue
end

break *0x1002e990
commands
  silent
  set $text_base = *(unsigned int *)($esp + 8)
  set $index_ptr = *(unsigned int *)($esp + 12)
  set $return_address = *(unsigned int *)$esp
  printf "RULE_SCAN offset=%d text=\n", *(int *)$index_ptr
  x/s ($text_base + *(int *)$index_ptr)
  tbreak *$return_address
  commands
    silent
    printf "RULE_SCAN_RETURN handled=%d next_offset=%d\n", $eax, *(int *)$index_ptr
    continue
  end
  continue
end

continue
