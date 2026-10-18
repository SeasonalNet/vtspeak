set pagination off
set confirm off
set debuginfod enabled off
set logging file /probe/runtime-original-msi-2006-dict.log
set logging overwrite on
set logging enabled on
set $lookup_calls = 0
break *0x408187
commands
  silent
  printf "DLL_LOADED_BEFORE_FIRST_API_CALL\n"
  set $decode_calls = 0
  break *0x100040a0
  commands
    silent
    set $decode_calls = $decode_calls + 1
    if $decode_calls <= 300
      set $decoded = *(unsigned int *)($esp + 4)
      set $payload = *(unsigned int *)($esp + 8)
      set $ret = *(unsigned int *)$esp
      if $decode_calls == 1
        printf "LEGACY_PHONE_ID_TABLE_256X5_BYTES\n"
        x/1280bx 0x1009e9a0
      end
      printf "EMBEDDED_PAYLOAD_DECODE_ENTRY n=%u payload=", $decode_calls
      x/32bx $payload
      tbreak *$ret
      commands
        silent
        printf "EMBEDDED_PAYLOAD_DECODE_RETURN n=%u status=%d flags_and_counts=", $decode_calls, $eax
        x/8hx $decoded
        printf "decoded_phone_bytes="
        x/48bx ($decoded + 0x10)
        printf "decoded_phone_string="
        x/s ($decoded + 0x10)
        continue
      end
    end
    continue
  end
  set $wrapper_calls = 0
  break *0x10003eb0
  commands
    silent
    set $wrapper_calls = $wrapper_calls + 1
    if $wrapper_calls <= 300
      set $surface = *(unsigned int *)($esp + 4)
      set $result = *(unsigned int *)($esp + 8)
      set $mode = *(short *)($esp + 12)
      set $ret = *(unsigned int *)$esp
      printf "EMBEDDED_LOOKUP_WRAPPER_ENTRY n=%u mode=%d surface=", $wrapper_calls, $mode
      x/s $surface
      tbreak *$ret
      commands
        silent
        printf "EMBEDDED_LOOKUP_WRAPPER_RETURN n=%u status=%d result=", $wrapper_calls, $eax
        x/64bx $result
        continue
      end
    end
    continue
  end
  break *0x1000e440
  commands
    silent
    set $lookup_calls = $lookup_calls + 1
    if $lookup_calls <= 300
      set $query = *(unsigned int *)($esp + 4)
      set $result = *(unsigned int *)($esp + 8)
      set $ret = *(unsigned int *)$esp
      set $query_len = *(unsigned int *)($query + 200)
      printf "EMBEDDED_HASH_LOOKUP_ENTRY n=%u query_len=%u query=", $lookup_calls, $query_len
      x/s $query
      printf "query_prefix="
      x/24bx $query
      tbreak *$ret
      commands
        silent
        printf "EMBEDDED_HASH_LOOKUP_RETURN n=%u status=%d result=", $lookup_calls, $eax
        x/s $result
        printf "result_prefix="
        x/96bx $result
        continue
      end
    end
    continue
  end
  continue
end
continue
