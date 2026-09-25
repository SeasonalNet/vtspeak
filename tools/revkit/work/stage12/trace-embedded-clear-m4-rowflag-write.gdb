set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $rowflag_hits = 0
break *0x10007520
commands
  silent
  set $ctx = *(unsigned int *)($esp + 4)
  set $rows = $ctx + 0x429a2
  set $entry_row = $rows + 4
  set {unsigned short}($entry_row + 0x6e) = 0
  set $return_address = *(unsigned int *)$esp
  set $rowflag_address = $entry_row
  printf "ROWFLAG_WATCH_SETUP address=%#x initial=%#x metadata4=%u\n", $rowflag_address, *(unsigned char *)$rowflag_address, *(unsigned short *)($entry_row + 0x6e)
  watch *(unsigned char *)$rowflag_address
  set $rowflag_watch = $bpnum
  commands
    silent
    set $rowflag_hits = $rowflag_hits + 1
    printf "ROWFLAG_WRITE hit=%u eip=%#x value=%#x metadata4=%u\n", $rowflag_hits, $eip, *(unsigned char *)$rowflag_address, *(unsigned short *)($entry_row + 0x6e)
    x/6i $eip-8
    bt 8
    if $rowflag_hits >= 24
      disable $rowflag_watch
      printf "ROWFLAG_WATCH_DISABLED\n"
    end
    continue
  end
  tbreak *$return_address
  commands
    silent
    printf "ROWFLAG_RETURN value=%#x metadata4=%u\n", *(unsigned char *)$rowflag_address, *(unsigned short *)($entry_row + 0x6e)
    continue
  end
  continue
end
continue
