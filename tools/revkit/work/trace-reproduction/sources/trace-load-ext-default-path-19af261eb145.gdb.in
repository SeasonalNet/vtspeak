set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands 1
  silent
  set $slot = *(int *)($esp + 8)
  set $load_return = *(unsigned int *)$esp
  printf "LOAD_EXT_ENTRY slot=%d db_path_arg=%s\n", $slot, (char *)*(unsigned int *)($esp + 12)
  disable 1
  break *0x10025e70
  commands 2
    silent
    printf "LOAD_DEFAULT_PATH_BEFORE_INIT slot=%d db_path_global=%s\n", $slot, (char *)0x1009fa4c
    disable 2
    continue
  end
  tbreak *$load_return
  commands 3
    silent
    set $dbsize = (int *)malloc(4)
    set *$dbsize = 0
    set $dbsize_result = ((int (*)(int *, int))0x10028130)($dbsize, $slot)
    printf "LOAD_EXT_RETURN slot=%d ax=%#x eax=%#x dbsize_result=%d dbsize=%d\n", $slot, $ax, $eax, $dbsize_result, *$dbsize
    continue
  end
  continue
end

continue
