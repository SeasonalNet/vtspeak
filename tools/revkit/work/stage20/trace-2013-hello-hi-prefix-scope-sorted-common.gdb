# FUN_10023c70 has just returned from FUN_1001b5f0 at 0x10023d43.
# Capture the sorted score-node order for the two affected feature queries.
break *0x10023d43
commands
  silent
  set $sorted_query = (unsigned int)($ebp - 8)
  set $sorted_key0 = *(unsigned char *)$sorted_query
  set $sorted_key1 = *(unsigned char *)($sorted_query + 1)
  set $sorted_key2 = *(unsigned char *)($sorted_query + 2)
  set $sorted_key3 = *(unsigned char *)($sorted_query + 3)
  set $sorted_key4 = *(unsigned char *)($sorted_query + 4)
  set $sorted_count = $esi
  # FUN_10023c70 replaces its param7 stack slot with the scratch-array
  # pointer in the prologue, so this slot already equals param7 + 0x477ac.
  set $sorted_nodes = *(unsigned int *)($ebp + 0x20)
  set $sorted_model = *(unsigned int *)($ebp + 0x24)
  set $sorted_class_keys = *(unsigned int *)($sorted_model + 0x90)
  set $sorted_member_counts = *(unsigned int *)($sorted_model + 0x8c)
  if $sorted_count >= 30
    printf "HELLO_SORTED_SCOPE query=%02x%02x%02x%02x%02x count=%u nodes=%08x model=%08x", $sorted_key0, $sorted_key1, $sorted_key2, $sorted_key3, $sorted_key4, $sorted_count, $sorted_nodes, $sorted_model
    set $sorted_i = 0
    while $sorted_i < $sorted_count && $sorted_i < 12
      set $sorted_node = *(unsigned int *)($sorted_nodes + $sorted_i * 4)
      if $sorted_node == 0
        printf " {rank=%u node=null}", $sorted_i
      else
        set $sorted_class = *(unsigned int *)($sorted_node + 8)
        printf " {%u key=%02x%02x%02x%02x%02x score=%u members=%u}", $sorted_class, *(unsigned char *)($sorted_class_keys + $sorted_class * 5), *(unsigned char *)($sorted_class_keys + $sorted_class * 5 + 1), *(unsigned char *)($sorted_class_keys + $sorted_class * 5 + 2), *(unsigned char *)($sorted_class_keys + $sorted_class * 5 + 3), *(unsigned char *)($sorted_class_keys + $sorted_class * 5 + 4), *(unsigned int *)($sorted_node + 4), *(unsigned short *)($sorted_member_counts + $sorted_class * 2)
      end
      set $sorted_i = $sorted_i + 1
    end
    printf "\n"
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_SORTED_SCOPE_TRACE_READY\n"
  continue
end

continue
