set $scope_distance_calls = 0

# Save the feature-model pointer while entering FUN_10023c70. Its candidate
# distance loop calls FUN_10023a70 and returns at 0x10023d14.
break *0x10023c70
commands
  silent
  set $scope_distance_model = *(unsigned int *)($esp + 32)
  continue
end

break *0x10023a70
commands
  silent
  set $distance_return = *(unsigned int *)$esp
  set $distance_query = *(unsigned int *)($esp + 4)
  set $distance_class_key = *(unsigned int *)($esp + 8)
  set $distance_width = *(short *)($esp + 12)
  set $distance_key0 = *(unsigned char *)$distance_query
  set $distance_key1 = *(unsigned char *)($distance_query + 1)
  set $distance_key2 = *(unsigned char *)($distance_query + 2)
  set $distance_key3 = *(unsigned char *)($distance_query + 3)
  set $distance_key4 = *(unsigned char *)($distance_query + 4)
  set $distance_class_table = *(unsigned int *)($scope_distance_model + 0x90)
  if $distance_return == 0x10023d14 && (($distance_key0 == 0x5a && $distance_key1 == 0x22 && $distance_key2 == 0x17) || ($distance_key0 == 0x2b && $distance_key1 == 0x30 && $distance_key2 == 0x5a))
    set $distance_class = ($distance_class_key - $distance_class_table) / 5
    set $scope_distance_calls = $scope_distance_calls + 1
  end
  continue
end

break *0x10023d14
commands
  silent
  if $distance_return == 0x10023d14 && (($distance_key0 == 0x5a && $distance_key1 == 0x22 && $distance_key2 == 0x17) || ($distance_key0 == 0x2b && $distance_key1 == 0x30 && $distance_key2 == 0x5a))
    printf "HELLO_SCOPE_DISTANCE n=%u query=%02x%02x%02x%02x%02x width=%d class=%u key=%02x%02x%02x%02x%02x score=%u\n", $scope_distance_calls, $distance_key0, $distance_key1, $distance_key2, $distance_key3, $distance_key4, $distance_width, $distance_class, *(unsigned char *)$distance_class_key, *(unsigned char *)($distance_class_key + 1), *(unsigned char *)($distance_class_key + 2), *(unsigned char *)($distance_class_key + 3), *(unsigned char *)($distance_class_key + 4), $eax
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_SCOPE_DISTANCE_TRACE_READY\n"
  continue
end

continue
