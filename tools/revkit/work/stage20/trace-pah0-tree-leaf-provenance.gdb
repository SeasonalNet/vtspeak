set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/pah0-tree-leaf-provenance/gdb.log
set logging overwrite on
set logging enabled on
set $split_calls = 0
set $signature_calls = 0

break *0x10024680
commands
  silent
  set $tree_pos = *(unsigned short *)($esp + 4)
  set $tree_state = *(unsigned int *)($esp + 8)
  set $tree_arg = *(unsigned int *)($esp + 12)
  set $tree_base = *(unsigned int *)($tree_state + 0x4c)
  printf "PAH0_PHONE_TREE_ENTRY pos=%u state=%#x arg=%#x tree_base=%#x\n", $tree_pos, $tree_state, $tree_arg, $tree_base
  continue
end

break *0x10024060
commands
  silent
  set $split_calls = $split_calls + 1
  set $query_slot = *(unsigned int *)($esp + 4)
  set $query_state = *(unsigned int *)($esp + 8)
  set $query_voice = *(unsigned int *)($query_state + 0x4c)
  set $query_row = $query_state + ($query_slot * 3 + 0x76314) * 2
  set $query_bank = *(unsigned char *)$query_row
  set $query_leaf = *(unsigned char *)($query_row + 2)
  set $query_bank_base = $query_voice + $query_bank * 0x3c0
  set $query_sig = $query_bank_base + 0x6e2 + $query_leaf * 7
  if $split_calls <= 32
    printf "PAH0_QUERY_ROW n=%u slot=%u return=%#x state=%#x voice=%#x row=%#x fields:", $split_calls, $query_slot, *(unsigned int *)$esp, $query_state, $query_voice, $query_row
    x/6ub $query_row
    printf "PAH0_QUERY_SIGNATURE n=%u bank=%u leaf=%u signature:", $split_calls, $query_bank, $query_leaf
    x/7ub $query_sig
    printf "PAH0_BANK_HEADER n=%u base=%#x signature_count=%u context_count=%u node_table=%#x feature_table=%#x\n", $split_calls, $query_bank_base, *(unsigned short *)($query_bank_base + 0x6de), *(unsigned char *)($query_bank_base + 0x6e0), *(unsigned int *)($query_bank_base + 0x658), *(unsigned int *)($query_bank_base + 0x654)
    printf "PAH0_BANK_SIGNATURES n=%u first_leaf_through_leaf:\n", $split_calls
    x/14ub ($query_bank_base + 0x6e2 + (($query_leaf / 7) * 7))
    printf "PAH0_BANK_CONTEXT_RECORDS n=%u records:\n", $split_calls
    x/48ub (*(unsigned int *)($query_bank_base + 0x658) + 1)
  end
  continue
end

break *0x10018770
commands
  silent
  set $signature_calls = $signature_calls + 1
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $signature = *(unsigned int *)($esp + 4)
    set $signature_slot = *(unsigned int *)($esp + 8)
    printf "PAH0_SIGNATURE_CALL n=%u caller=%#x slot=%u raw:", $signature_calls, $caller, $signature_slot
    x/7ub $signature
  end
  continue
end

break *0x408187
commands
  silent
  printf "PAH0_TREE_LEAF_PROVENANCE_READY\n"
  continue
end

continue
