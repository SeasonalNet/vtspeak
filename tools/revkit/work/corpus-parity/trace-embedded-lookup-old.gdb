set pagination off
set confirm off
set logging file /probe/embedded-lookup-old.log
set logging overwrite on
set logging enabled on
set $calls = 0
break *0x10011820
commands
  silent
  set $query = *(unsigned int *)($esp + 4)
  set $result = *(unsigned int *)($esp + 8)
  set $mode = *(int *)($esp + 12)
  set $family = *(int *)($esp + 16)
  set $ret = *(unsigned int *)$esp
  if $mode == 0 && $family == 0
    set $calls = $calls + 1
    printf "EMBED_LOOKUP_ENTRY call=%d query=", $calls
    x/s $query
    tbreak *$ret
    commands
      silent
      printf "EMBED_LOOKUP_RETURN call=%d status=%#x length=%u payload=", $calls, $eax, *(unsigned int *)($result + 200)
      if $eax == 1
        x/32bx $result
      end
      continue
    end
  end
  continue
end
continue
