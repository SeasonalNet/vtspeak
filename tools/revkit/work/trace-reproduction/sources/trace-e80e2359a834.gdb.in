set pagination off
set confirm off
set logging file /work/stage5/probes/combo/tree-lookups-gdb.log
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
  x/12hd $features
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
  x/12hd $features
  tbreak *$return_address
  commands
    silent
    printf "VECTOR_RETURN values:"
    x/12hd $output
    continue
  end
  continue
end

continue
