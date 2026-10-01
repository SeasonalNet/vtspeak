set $range_calls = 0

# FUN_1002df50 builds the feature-range candidate scope called by
# FUN_10019570 from FUN_10023e70.
break *0x1002df50
commands
  silent
  if *(unsigned int *)$esp == 0x100195a5
    set $range_calls = $range_calls + 1
    set $range_out = *(unsigned int *)($esp + 4)
    set $range_query = *(unsigned int *)($esp + 8)
    set $range_model = *(unsigned int *)($esp + 12)
    set $range_width = *(signed char *)($range_out + 8)
    set $range_view = *(unsigned char *)$range_query
    printf "HELLO_RANGE_SCOPE_ENTRY call=%u mode=%u width=%d query=", $range_calls, $range_view, $range_width
    set $range_i = 0
    while $range_i < 10
      printf "%02x", *(unsigned char *)($range_query + $range_i)
      set $range_i = $range_i + 1
    end
    printf "\n"
  end
  continue
end

break *0x100195a5
commands
  silent
  set $range_count = *(int *)($range_out + 4)
  set $range_classes = *(unsigned int *)$range_out
  set $range_keys = *(unsigned int *)($range_model + 0x90)
  set $range_units = *(unsigned int *)($range_model + 0x8c)
  printf "HELLO_RANGE_SCOPE_RETURN call=%u count=%u classes:", $range_calls, $range_count
  set $range_i = 0
  while $range_i < $range_count && $range_i < 128
    set $range_class = *(unsigned int *)($range_classes + $range_i * 4)
    printf " {%u key=%02x%02x%02x%02x%02x members=%u}", $range_class, *(unsigned char *)($range_keys + $range_class * 5), *(unsigned char *)($range_keys + $range_class * 5 + 1), *(unsigned char *)($range_keys + $range_class * 5 + 2), *(unsigned char *)($range_keys + $range_class * 5 + 3), *(unsigned char *)($range_keys + $range_class * 5 + 4), *(unsigned short *)($range_units + $range_class * 2)
    set $range_i = $range_i + 1
  end
  printf "\n"
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_RANGE_SCOPE_TRACE_READY\n"
  continue
end

continue
