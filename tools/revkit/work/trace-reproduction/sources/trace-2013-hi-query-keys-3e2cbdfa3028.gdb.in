set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-2013-query-keys/adapted-query-keys-gdb.log
set logging overwrite on
set logging enabled on
set $query_calls = 0

break *0x10023af0
commands
  silent
  set $query_calls = $query_calls + 1
  set $min_count = *(unsigned int *)($esp + 4)
  set $max_weight = *(int *)($esp + 8)
  set $result_array = *(unsigned int *)($esp + 12)
  set $signature = *(unsigned int *)($esp + 16)
  set $state = *(unsigned int *)($esp + 20)
  set $model = *(unsigned int *)($esp + 24)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    set $packed = $eax
    set $count = $packed & 0xffff
    printf "HI_2013_QUERY call=%u min_count=%u max_weight=%d signature=%02x,%02x,%02x,%02x,%02x,%02x,%02x key=%02x,%02x,%02x,%02x,%02x result_count=%u ids:", $query_calls, $min_count, $max_weight, *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), *(unsigned char *)($ebp - 0xc), *(unsigned char *)($ebp - 0xb), *(unsigned char *)($ebp - 0xa), *(unsigned char *)($ebp - 9), *(unsigned char *)($ebp - 8), $count
    set $i = 0
    while $i < $count && $i < 80
      printf " %u", *(unsigned int *)($result_array + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    continue
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_2013_QUERY_KEY_TRACE_READY\n"
  continue
end

continue
