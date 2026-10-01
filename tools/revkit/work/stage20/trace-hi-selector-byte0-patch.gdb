set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-selector-byte0-patch/adapted-suffix-map-gdb.log
set logging overwrite on
set logging enabled on
set $lookup_calls = 0

break *0x10023dc0
commands
  silent
  set $lookup_calls = $lookup_calls + 1
  set $signature = *(unsigned int *)($ebp + 8)
  set $return = *(unsigned int *)($ebp + 4)
  set $key0 = *(unsigned char *)(0x1007b7ec + *(unsigned char *)($signature + 1))
  set $key1 = *(unsigned char *)(0x1007b788 + *(unsigned char *)($signature + 2))
  set $key2 = *(unsigned char *)(0x1007b850 + *(unsigned char *)($signature + 3))
  set $key3 = *(unsigned char *)($signature + 5)
  set $key4 = *(unsigned char *)($signature + 6) & 0x20
  tbreak *$return
  commands
    silent
    set $count = $eax & 0xffff
    printf "HI_SUFFIX_LOOKUP call=%u signature=%02x,%02x,%02x,%02x,%02x,%02x,%02x key=%02x,%02x,%02x,%02x,%02x count=%u\n", $lookup_calls, *(unsigned char *)$signature, *(unsigned char *)($signature + 1), *(unsigned char *)($signature + 2), *(unsigned char *)($signature + 3), *(unsigned char *)($signature + 4), *(unsigned char *)($signature + 5), *(unsigned char *)($signature + 6), $key0, $key1, $key2, $key3, $key4, $count
    continue
  end
  continue
end

break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  set $found_272820 = 0
  set $found_272821 = 0
  set $found_65331 = 0
  set $index = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $unit = *(unsigned int *)($node + 8)
    if $unit == 272820
      set $found_272820 = 1
    end
    if $unit == 272821
      set $found_272821 = 1
    end
    if $unit == 65331
      set $found_65331 = 1
    end
    set $index = $index + 1
  end
  printf "HI_SUFFIX_LIST position=%d count=%d has272820=%d has272821=%d has65331=%d\n", $position, $count, $found_272820, $found_272821, $found_65331
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "HI_SUFFIX_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_SUFFIX_TRACE_READY\n"
  continue
end

continue
