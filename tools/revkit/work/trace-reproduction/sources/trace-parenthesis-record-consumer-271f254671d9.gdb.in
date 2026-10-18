set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  printf "PAREN_RECORD_TEXT_ENTRY format=%d text=%s\n", *(int *)($esp + 4), *(char **)($esp + 8)
  call ((void (*)(int))0x10028390)(1)
  disable 1
  continue
end

break *0x1005567f
commands 2
  silent
  set $paren_record_address = $edi + $edx * 4 - 0x54
  printf "PAREN_RECORD_MARKED address=%#x index_state=%d value=%u\n", $paren_record_address, *(int *)($ebp - 4), *(unsigned int *)$paren_record_address
  awatch *(unsigned int *)$paren_record_address
  commands $bpnum
    silent
    printf "PAREN_RECORD_ACCESS pc=%#x eax=%#x ecx=%#x edx=%#x value=%u\n", $eip, $eax, $ecx, $edx, *(unsigned int *)$paren_record_address
    bt 8
    continue
  end
  disable 2
  continue
end

continue
