set $lookup_calls = 0

break *0x10023f90
commands
  silent
  set $lookup_calls = $lookup_calls + 1
  set $lookup_context = *(int *)($esp + 8)
  set $lookup_signature = *(unsigned int *)($esp + 4)
  set $lookup_state = *(unsigned int *)($esp + 12)
  set $lookup_model = *(unsigned int *)($esp + 16)
  set $lookup_start = *(short *)($lookup_state + 0xec620)
  set $lookup_output = $lookup_state + 0xec170 + $lookup_start * 4
  set $lookup_key0 = *(unsigned char *)(0x1007b7ec + *(unsigned char *)($lookup_signature + 1))
  set $lookup_key1 = *(unsigned char *)(0x1007b788 + *(unsigned char *)($lookup_signature + 2))
  set $lookup_key2 = *(unsigned char *)(0x1007b850 + *(unsigned char *)($lookup_signature + 3))
  set $lookup_key3 = *(unsigned char *)($lookup_signature + 5)
  set $lookup_key4 = *(unsigned char *)($lookup_signature + 6) & 0x20
  printf "HELLO_QUERY_MEMBERSHIP_ENTRY call=%u context=%u start=%d signature=%02x,%02x,%02x,%02x,%02x,%02x,%02x key=%02x,%02x,%02x,%02x,%02x\n", $lookup_calls, $lookup_context, $lookup_start, *(unsigned char *)$lookup_signature, *(unsigned char *)($lookup_signature + 1), *(unsigned char *)($lookup_signature + 2), *(unsigned char *)($lookup_signature + 3), *(unsigned char *)($lookup_signature + 4), *(unsigned char *)($lookup_signature + 5), *(unsigned char *)($lookup_signature + 6), $lookup_key0, $lookup_key1, $lookup_key2, $lookup_key3, $lookup_key4
  continue
end

# FUN_10023f90 returns from its direct FUN_10023dc0 call here, before adding
# the returned count to the state cursor.
break *0x10023fca
commands
  silent
  set $lookup_count = $eax & 0xffff
  printf "HELLO_QUERY_MEMBERSHIP_RETURN call=%u count=%u classes:", $lookup_calls, $lookup_count
  set $lookup_i = 0
  while $lookup_i < $lookup_count && $lookup_i < 300
    set $lookup_class = *(unsigned int *)($lookup_output + $lookup_i * 4)
    set $lookup_keys = *(unsigned int *)($lookup_model + 0x90)
    set $lookup_units = *(unsigned int *)($lookup_model + 0x8c)
    printf " {%u key=%02x%02x%02x%02x%02x members=%u}", $lookup_class, *(unsigned char *)($lookup_keys + $lookup_class * 5), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 1), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 2), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 3), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 4), *(unsigned short *)($lookup_units + $lookup_class * 2)
    set $lookup_i = $lookup_i + 1
  end
  printf "\n"
  continue
end

# The half-position/relaxed path returns through FUN_10023e70 to this point.
break *0x10024006
commands
  silent
  set $lookup_count = $eax & 0xffff
  printf "HELLO_QUERY_MEMBERSHIP_RELAXED_RETURN call=%u count=%u classes:", $lookup_calls, $lookup_count
  set $lookup_i = 0
  while $lookup_i < $lookup_count && $lookup_i < 300
    set $lookup_class = *(unsigned int *)($lookup_output + $lookup_i * 4)
    set $lookup_keys = *(unsigned int *)($lookup_model + 0x90)
    set $lookup_units = *(unsigned int *)($lookup_model + 0x8c)
    printf " {%u key=%02x%02x%02x%02x%02x members=%u}", $lookup_class, *(unsigned char *)($lookup_keys + $lookup_class * 5), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 1), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 2), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 3), *(unsigned char *)($lookup_keys + $lookup_class * 5 + 4), *(unsigned short *)($lookup_units + $lookup_class * 2)
    set $lookup_i = $lookup_i + 1
  end
  printf "\n"
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_QUERY_MEMBERSHIP_TRACE_READY\n"
  continue
end

continue
