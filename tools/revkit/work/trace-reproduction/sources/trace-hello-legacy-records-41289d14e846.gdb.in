set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-legacy-records-gdb.log
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
  set $assembly_calls = 0
  set $legacy_unit_calls = 0
  set $legacy_audio_calls = 0
  set $legacy_append_calls = 0
  break *0x10023b40
  commands
    silent
    set $legacy_append_calls = $legacy_append_calls + 1
    if $legacy_append_calls <= 32
      set $append_out = *(unsigned int *)($esp + 4)
      set $append_state = *(unsigned int *)($esp + 8)
      set $append_index = *(short *)($esp + 12)
      set $append_return = *(unsigned int *)$esp
      set $append_count = *(short *)($append_state + 0xdf74)
      set $append_samples = *(int *)($append_state + 11000 + $append_index * 4)
      set $append_flag = *(char *)($append_state + 0x7855 + $append_index)
      set $append_overlap = *(short *)($append_state + 0x6fbc + $append_index * 2)
      set $append_used = $append_samples
      if $append_index < $append_count - 1 && $append_flag != 1
        set $append_used = $append_samples - $append_overlap
      end
      printf "LEGACY_APPEND_ENTRY call=%u index=%d count=%d raw_samples=%d flag=%d overlap=%d appended_samples=%d cursor_before=%d\n", $legacy_append_calls, $append_index, $append_count, $append_samples, $append_flag, $append_overlap, $append_used, *(int *)($append_out + 0x38)
      tbreak *$append_return
      commands
        silent
        printf "LEGACY_APPEND_RETURN call=%u cursor_after=%d\n", $legacy_append_calls, *(int *)($append_out + 0x38)
        continue
      end
    end
    continue
  end
  break *0x10023a10
  commands
    silent
    set $legacy_audio_calls = $legacy_audio_calls + 1
    if $legacy_audio_calls <= 32
      set $audio_state = *(unsigned int *)($esp + 8)
      set $audio_index = *(short *)($esp + 12)
      set $audio_return = *(unsigned int *)$esp
      set $audio_samples = *(int *)($audio_state + 11000 + $audio_index * 4)
      set $audio_record = *(unsigned int *)($audio_state + 0x3c28 + $audio_index * 4)
      printf "LEGACY_AUDIO_SEGMENT_ENTRY call=%u index=%d samples=%d record=%#x\n", $legacy_audio_calls, $audio_index, $audio_samples, $audio_record
      x/20bx $audio_record
      tbreak *$audio_return
      commands
        silent
        set $audio_data = *(unsigned int *)($audio_state + $audio_index * 4)
        printf "LEGACY_AUDIO_SEGMENT_RETURN call=%u data=%#x samples=%d first_samples=", $legacy_audio_calls, $audio_data, $audio_samples
        x/12hd $audio_data
        continue
      end
    end
    continue
  end
  break *0x10013cf0
  commands
    silent
    set $legacy_unit_calls = $legacy_unit_calls + 1
    if $legacy_unit_calls <= 32
      set $out = *(unsigned int *)($esp + 4)
      set $unit_id = *(unsigned int *)($esp + 8)
      set $context = *(unsigned int *)($esp + 12)
      set $mode = *(unsigned int *)($esp + 16)
      set $ret = *(unsigned int *)$esp
      printf "LEGACY_UNIT_RECORD_ENTRY call=%u unit_id=%u context=%#x mode=%u out=%#x\n", $legacy_unit_calls, $unit_id, $context, $mode, $out
      if $legacy_unit_calls == 1
        set $manager = *(unsigned int *)($mode + 0x88)
        set $range_count = *(short *)($mode + 0xe98) - 1
        printf "LEGACY_UNIT_BANK_RANGES manager=%#x ranges=%d first_count=%u\n", $manager, $range_count, *(unsigned int *)($manager + 0x50)
        set $range_i = 0
        while $range_i < $range_count && $range_i < 8
          set $range_start = *(int *)($manager + 0xc8 + $range_i * 0x68)
          set $range_end = *(int *)($manager + 0x130 + $range_i * 0x68)
          printf "LEGACY_UNIT_BANK_RANGE index=%u start=%u end=%u\n", $range_i, $range_start, $range_end
          set $range_i = $range_i + 1
        end
      end
      tbreak *$ret
      commands
        silent
        printf "LEGACY_UNIT_RECORD_RETURN call=%u bytes:\n", $legacy_unit_calls
        x/20bx $out
        continue
      end
    end
    continue
  end
  break *0x100232d0
  commands
    silent
    set $assembly_calls = $assembly_calls + 1
    if $assembly_calls <= 8
      set $out = *(unsigned int *)($esp + 4)
      set $state = *(unsigned int *)($esp + 8)
      set $start = *(short *)($esp + 12)
      set $count = *(short *)($esp + 16)
      set $text_state = *(unsigned int *)($state + 0x54)
      set $record_table = $text_state + 0x730 + $start * 0x2fc
      printf "LEGACY_ASSEMBLY_ENTRY call=%u out=%#x state=%#x start=%d count=%d text_state=%#x\n", $assembly_calls, $out, $state, $start, $count, $text_state
      printf "LEGACY_PHONE_RECORD_TABLE words="
      x/8wx $record_table
      set $records = *(unsigned int *)$record_table
      set $record_count = *(short *)($record_table + 8)
      printf "LEGACY_PHONE_RECORDS pointer=%#x count=%d\n", $records, $record_count
      set $record_i = 0
      while $record_i < $record_count && $record_i < 16
        printf "LEGACY_PHONE_RECORD i=%u bytes:", $record_i
        x/20bx $records + $record_i * 0x14
        set $record_i = $record_i + 1
      end
    end
    continue
  end
  continue
end
continue
