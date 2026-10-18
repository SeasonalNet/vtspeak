set pagination off
set confirm off
break *0x10001670
commands
  silent
  set $tree = *(unsigned int *)($esp + 4)
  set $features = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "TREE_SCALAR tree=%#x return=%#x features:", $tree, $ret
  x/16hd $features
  tbreak *$ret
  commands
    silent
    printf "TREE_SCALAR_VALUE value=%d\n", (short)$eax
    continue
  end
  continue
end
break *0x100016a0
commands
  silent
  set $tree = *(unsigned int *)($esp + 4)
  set $features = *(unsigned int *)($esp + 8)
  set $out = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "TREE_VECTOR tree=%#x return=%#x features:", $tree, $ret
  x/16hd $features
  tbreak *$ret
  commands
    silent
    printf "TREE_VECTOR_VALUE values:"
    x/12hd $out
    continue
  end
  continue
end
continue
