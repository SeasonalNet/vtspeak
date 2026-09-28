set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-phone-tree-provenance-gdb.log
set logging overwrite on
set logging enabled on
set $tree_calls = 0
set $query_calls = 0

hbreak *0x10024680
commands
  silent
  set $tree_calls = $tree_calls + 1
  set $position = *(unsigned short *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $tree_base = *(unsigned int *)($state + 0x4c)
  set $position_record = $tree_base + $position * 0x10
  printf "HELLO_PHONE_TREE_ENTRY n=%u position=%u state=%#x tree_base=%#x record:", $tree_calls, $position, $state, $tree_base
  x/8uh $position_record
  continue
end

hbreak *0x10024060
commands
  silent
  set $query_calls = $query_calls + 1
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $row = $state + ($slot * 3 + 0x76314) * 2
  set $bank = *(unsigned char *)$row
  set $leaf = *(unsigned char *)($row + 2)
  set $voice = *(unsigned int *)($state + 0x4c)
  set $bank_base = $voice + $bank * 0x3c0
  set $signature = $bank_base + 0x6e2 + $leaf * 7
  printf "HELLO_PHONE_TREE_QUERY n=%u slot=%u bank=%u feature=%u leaf=%u flags=%u row:", $query_calls, $slot, $bank, *(unsigned char *)($row + 1), $leaf, *(unsigned char *)($row + 5)
  x/6ub $row
  printf "HELLO_PHONE_TREE_BANK count=%u categories=%u feature_table=%#x signature:", *(unsigned short *)($bank_base + 0x6de), *(unsigned char *)($bank_base + 0x6e0), *(unsigned int *)($bank_base + 0x654)
  x/7ub $signature
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_PHONE_TREE_PROVENANCE_READY\n"
  continue
end

continue
