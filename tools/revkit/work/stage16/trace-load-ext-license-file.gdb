set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10027af0
commands 1
  silent
  set $slot = *(int *)($esp + 8)
  set $db_path = (char *)malloc(9)
  set {char[9]}$db_path = "Z:\\work\\"
  set $license_path = (char *)malloc(44)
  set {unsigned int}($license_path + 0) = 0x775c3a5a
  set {unsigned int}($license_path + 4) = 0x5c6b726f
  set {unsigned int}($license_path + 8) = 0x61746164
  set {unsigned int}($license_path + 12) = 0x6d6f632d
  set {unsigned int}($license_path + 16) = 0x5c6e6f6d
  set {unsigned int}($license_path + 20) = 0x69726576
  set {unsigned int}($license_path + 24) = 0x765c7966
  set {unsigned int}($license_path + 28) = 0x66697265
  set {unsigned int}($license_path + 32) = 0x74616369
  set {unsigned int}($license_path + 36) = 0x2e6e6f69
  set {unsigned int}($license_path + 40) = 0x747874
  set *(unsigned int *)($esp + 12) = $db_path
  set *(unsigned int *)($esp + 24) = $license_path
  set $arg4 = *(unsigned int *)($esp + 16)
  set $arg5 = *(unsigned int *)($esp + 20)
  set $license_data = *(unsigned int *)($esp + 28)
  set $license_size = *(unsigned int *)($esp + 32)
  set $load_return = *(unsigned int *)$esp
  printf "LOAD_EXT_ENTRY slot=%d db_path=%s arg4=%#x arg5=%#x license_path=%s license_data=%#x license_size=%d\n", $slot, $db_path, $arg4, $arg5, $license_path, $license_data, $license_size
  disable 1
  break *0x10025e70
  commands 2
    silent
    printf "LOAD_BEFORE_INIT slot=%d db_path_global=%s\n", $slot, (char *)0x1009fa4c
    disable 2
    continue
  end
  break *0x10029b80
  commands 3
    silent
    set $check_path = *(unsigned int *)($esp + 4)
    set $check_data = *(unsigned int *)($esp + 8)
    set $check_size = *(unsigned int *)($esp + 12)
    set $check_return = *(unsigned int *)$esp
    printf "CHECK_LICENSE_ENTRY file_path=%s data_pointer=%#x size=%d\n", (char *)$check_path, $check_data, $check_size
    disable 3
    tbreak *$check_return
    commands 5
      silent
      printf "CHECK_LICENSE_RETURN eax=%d\n", $eax
      continue
    end
    continue
  end
  tbreak *$load_return
  commands 4
    silent
    set $dbsize = (int *)malloc(4)
    set *$dbsize = 0
    set $dbsize_result = ((int (*)(int *, int))0x10028130)($dbsize, $slot)
    printf "LOAD_EXT_RETURN slot=%d ax=%d dbsize_result=%d dbsize=%d\n", $slot, $ax, $dbsize_result, *$dbsize
    set $slot_state = *(unsigned int *)0x100a0468
    printf "LOAD_LICENSE_STATE normalized_slot=1 gate=%d dictionary_capacity=%d\n", *(unsigned char *)0x100a7489, *(int *)($slot_state + 0x4d14)
    continue
  end
  continue
end

continue
