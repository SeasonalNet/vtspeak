set pagination off
set confirm off
set logging file /probe/cnas-path-gdb.log
set logging overwrite on
set logging enabled on

break *0x10001670
commands
  silent
  set $tree = *(unsigned int *)($esp + 4)
  set $features = *(unsigned int *)($esp + 8)
  set $return_address = *(unsigned int *)$esp
  set $nodes = *(short *)($tree + 24)
  if $nodes == 281
    printf "CNAS_CALL return=%#x tree=%#x features=%#x\n", $return_address, $tree, $features
    x/12hd $features
    tbreak *$return_address
    commands
      silent
      printf "CNAS_RETURN value=%d\n", (int)$eax
      continue
    end
  end
  continue
end

break *0x1000160b
commands
  silent
  if *(short *)($ebx + 24) == 281
    set $root = *(unsigned int *)$ebx
    set $node_index = ($esi - $root) / 16
    set $caller = *(unsigned int *)(*(unsigned int *)$ebp + 4)
    set $operation = *(unsigned char *)($esi + 11)
    set $feature = (signed char)*(unsigned char *)($esi + 10)
    if $operation == 67
      set $feature_value = *(short *)($edi + $feature * 2)
      set $threshold = *(short *)($esi + 4)
      set $true_child = *(short *)($esi + 6)
      set $false_child = *(short *)($esi + 8)
      printf "CNAS_NODE caller=%#x node=%d op=C feature=%d value=%d threshold=%d true=%d false=%d\n", $caller, $node_index, $feature, $feature_value, $threshold, $true_child, $false_child
    else
      printf "CNAS_NODE caller=%#x node=%d op=%c feature=%d left=%d right=%d\n", $caller, $node_index, $operation, $feature, *(short *)($esi + 6), *(short *)($esi + 8)
    end
  end
  continue
end

continue
