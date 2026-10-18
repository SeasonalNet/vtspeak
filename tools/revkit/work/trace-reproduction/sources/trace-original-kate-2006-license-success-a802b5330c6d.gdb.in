set pagination off
set confirm off
set debuginfod enabled off
set logging file /probe/runtime-license-success-2006-parser.log
set logging overwrite on
set logging enabled on
set $parser_calls = 0
break *0x408187
commands
  silent
  printf "DLL_LOADED_BEFORE_FIRST_API_CALL\n"
  break *0x100546d0
  commands
    silent
    set $parser_calls = $parser_calls + 1
    set $state = *(unsigned int *)($esp + 4)
    set $source = *(unsigned int *)($esp + 8)
    set $result = *(unsigned int *)$esp
    if $parser_calls <= 300
      printf "LEGACY_TEXT_PARSE_ENTRY n=%u source=", $parser_calls
      x/s $source
      printf "state=%#x result_buffer=%#x\n", $state, *(unsigned int *)$state
      tbreak *$result
      commands
        silent
        set $record_base = *(unsigned int *)$state
        set $record_count = *(unsigned int *)($record_base + 4)
        printf "LEGACY_TEXT_PARSE_RETURN n=%u status=%d consumed=%u records=%u\n", $parser_calls, $eax, *(unsigned int *)$record_base, $record_count
        set $record_i = 0
        while $record_i < $record_count && $record_i < 32
          set $record_start = *(unsigned int *)($record_base + 0x14 + $record_i * 0x88)
          set $record_end = *(unsigned int *)($record_base + 0x18 + $record_i * 0x88)
          printf "LEGACY_RECORD n=%u i=%u span=%u..%u\n", $parser_calls, $record_i, $record_start, $record_end
          set $record_i = $record_i + 1
        end
        continue
      end
    end
    continue
  end
  continue
end
continue
