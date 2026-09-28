set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/continuity-marker-provenance-gdb.log
set logging overwrite on
set logging enabled on
set $loader_calls = 0
set $index_calls = 0
set $loader_model = 0
set $loader_signatures = 0

# Record the model buffer and then compare it with the per-index reader's
# transposed key column before and after the versioned read.
break *0x10019f40
commands
  silent
  set $loader_calls = $loader_calls + 1
  set $loader_model = *(unsigned int *)($esp + 4)
  set $loader_signatures = *(unsigned int *)($loader_model + 0x64)
  printf "MODEL_LOADER call=%u object=%#x signature_table=%#x\n", $loader_calls, $loader_model, $loader_signatures
  continue
end

break *0x10019940
commands
  silent
  set $index_calls = $index_calls + 1
  set $index = *(unsigned int *)($esp + 4)
  set $units = *(unsigned int *)($index + 0x4c)
  if $units == 98133
    set $key_column = *(unsigned int *)($index + 0x44)
    set $return = *(unsigned int *)$esp
    printf "GEN2_INDEX_ENTRY call=%u object=%#x units=%u key_column=%#x model_table=%#x model_plus_unit_offset=%#x\n", $index_calls, $index, $units, $key_column, $loader_signatures, $loader_signatures + 272822 * 7
    printf "GEN2_BEFORE source_row92827_byte6=%#x model_unit272822_byte6=%#x\n", *(unsigned char *)($key_column + 6 * $units + 92827), *(unsigned char *)($loader_signatures + 272822 * 7 + 6)
    tbreak *$return
    commands
      silent
      printf "GEN2_INDEX_RETURN source_row92827_byte6=%#x model_unit272822_byte6=%#x\n", *(unsigned char *)($key_column + 6 * $units + 92827), *(unsigned char *)($loader_signatures + 272822 * 7 + 6)
      continue
    end
  end
  continue
end

break *0x10023350
commands
  silent
  set $model = *(unsigned int *)($esp + 12)
  set $signature_table = *(unsigned int *)($model + 0x64)
  printf "SELECTOR_MODEL object=%#x signature_table=%#x ids272822_25_byte6=%#x,%#x,%#x,%#x\n", $model, $signature_table, *(unsigned char *)($signature_table + 272822 * 7 + 6), *(unsigned char *)($signature_table + 272823 * 7 + 6), *(unsigned char *)($signature_table + 272824 * 7 + 6), *(unsigned char *)($signature_table + 272825 * 7 + 6)
  continue
end

continue
