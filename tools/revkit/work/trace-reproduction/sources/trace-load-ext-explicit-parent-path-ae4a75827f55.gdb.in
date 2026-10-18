set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands 1
  silent
  set $slot = *(int *)($esp + 8)
  set $db_path = (char *)malloc(19)
  set {char[19]}$db_path = "Z:\\voice-data\\paul"
  set *(unsigned int *)($esp + 12) = $db_path
  set $arg4 = *(unsigned int *)($esp + 16)
  set $arg5 = *(unsigned int *)($esp + 20)
  set $arg6 = *(unsigned int *)($esp + 24)
  set $arg7 = *(unsigned int *)($esp + 28)
  set $arg8 = *(unsigned int *)($esp + 32)
  set $load_return = *(unsigned int *)$esp
  printf "LOAD_EXT_ENTRY slot=%d db_path=%s arg4=%#x arg5=%#x license_path=%#x license_data=%#x license_size=%#x\n", $slot, $db_path, $arg4, $arg5, $arg6, $arg7, $arg8
  disable 1
  break *0x10025f10
  commands 2
    silent
    printf "LOAD_CONTEXT_ENTRY slot=%d db_path_global=%s\n", $slot, (char *)0x1009fa4c
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
