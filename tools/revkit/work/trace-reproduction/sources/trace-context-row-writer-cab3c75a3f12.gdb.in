set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $row_watch_installed = 0
set $row_watch_state = 0
set $row_watch_addr = 0
set $row_watch_value = 0

# FUN_100267d0 receives the engine state as its second argument. Watch
# context 1's row byte +0x04 before the selection call.
break *0x100267d0
commands
  silent
  if $row_watch_installed == 0
    set $row_watch_state = *(unsigned int *)($esp + 8)
    set $row_watch_addr = $row_watch_state + 0xec628 + 6 + 4
    set $row_watch_value = *(unsigned char *)$row_watch_addr
    printf "CONTEXT_ROW4_WATCH_SETUP state=%#x address=%#x initial=%u\n", $row_watch_state, $row_watch_addr, $row_watch_value
    watch *(unsigned char *)$row_watch_addr
    set $row_watch_installed = 1
    commands
      silent
      printf "CONTEXT_ROW4_WRITE pc=%#x old=%u new=%u\n", $pc, $row_watch_value, *(unsigned char *)$row_watch_addr
      x/i $pc
      bt 8
      set $row_watch_value = *(unsigned char *)$row_watch_addr
      continue
    end
  end
  continue
end

continue
