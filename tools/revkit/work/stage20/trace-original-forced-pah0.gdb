set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/original-forced-pah0-gdb.log
set logging overwrite on
set logging enabled on
set $tree_calls = 0
set $legacy_tree_calls = 0
set $parser_calls = 0
set $engine = 0

break *0x408187
commands
  silent
  printf "OLD_ENGINE_HOST_READY_BEFORE_SYNTHESIS\n"
  break *0x1000fd30
  commands
    silent
    set $engine = *(unsigned int *)($esp + 4)
    printf "OLD_PHONE_TREE_DISPATCH engine=%#x\n", $engine
    continue
  end
  break *0x100546d0
  commands
    silent
    set $parser_calls = $parser_calls + 1
    set $state = *(unsigned int *)($esp + 4)
    set $source = *(unsigned int *)($esp + 8)
    set $ret = *(unsigned int *)$esp
    printf "OLD_TEXT_PARSE_ENTRY n=%u source=", $parser_calls
    x/s $source
    tbreak *$ret
    commands
      silent
      set $record_base = *(unsigned int *)$state
      set $record_count = *(unsigned int *)($record_base + 4)
      printf "OLD_TEXT_PARSE_RETURN n=%u status=%d consumed=%u records=%u\n", $parser_calls, $eax, *(unsigned int *)$record_base, $record_count
      set $record_i = 0
      while $record_i < $record_count && $record_i < 24
        set $record_start = *(unsigned int *)($record_base + 0x14 + $record_i * 0x88)
        set $record_end = *(unsigned int *)($record_base + 0x18 + $record_i * 0x88)
        printf "OLD_TEXT_RECORD n=%u i=%u span=%u..%u\n", $parser_calls, $record_i, $record_start, $record_end
        set $record_i = $record_i + 1
      end
      continue
    end
    continue
  end
  break *0x10001570
  commands
    silent
    set $tree_calls = $tree_calls + 1
    if $tree_calls <= 128
      set $root = *(unsigned int *)($esp + 4)
      set $features = *(unsigned int *)($esp + 8)
      set $ret = *(unsigned int *)$esp
      printf "TREE2_EVAL_ENTRY n=%u root=%#x features=%#x root_node=", $tree_calls, $root, $features
      x/5wx $root
      printf "TREE2_FEATURE_VECTOR n=%u values:\n", $tree_calls
      x/32hd $features
      tbreak *$ret
      commands
        silent
        printf "TREE2_EVAL_RETURN n=%u output=%d\n", $tree_calls, (short)$eax
        continue
      end
    end
    continue
  end
  break *0x100019b0
  commands
    silent
    set $legacy_tree_calls = $legacy_tree_calls + 1
    set $tree = *(unsigned int *)($esp + 4)
    set $features = *(unsigned int *)($esp + 8)
    set $out = *(unsigned int *)($esp + 12)
    set $ret = *(unsigned int *)$esp
    if $legacy_tree_calls <= 32
      printf "TREE2_VECTOR_ENTRY n=%u tree=%#x tree_offset=%#x features=%#x output=%#x root_node=", $legacy_tree_calls, $tree, $tree - $engine, $features, $out
      printf "TREE2_VECTOR_HEADER width=%u\n", *(unsigned char *)($tree + 0x14)
      printf "TREE2_TREE_SLOTS nbf=%#x bf=%#x sbf=%#x qbf=%#x\n", *(unsigned int *)($engine + 0x228), *(unsigned int *)($engine + 0x244), *(unsigned int *)($engine + 0x260), *(unsigned int *)($engine + 0x27c)
      x/12wx $tree
      printf "TREE2_VECTOR_FEATURES n=%u values:\n", $legacy_tree_calls
      x/32hd $features
      tbreak *$ret
      commands
        silent
        printf "TREE2_VECTOR_RETURN n=%u values:\n", $legacy_tree_calls
        x/12hd $out
        continue
      end
    end
    continue
  end
  continue
end
continue
