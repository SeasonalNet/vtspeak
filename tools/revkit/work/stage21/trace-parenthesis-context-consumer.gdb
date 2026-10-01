set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands 1
  silent
  printf "PAREN_CONTEXT_TEXT_ENTRY format=%d text=%s\n", *(int *)($esp + 4), *(char **)($esp + 8)
  call ((void (*)(int))0x10028390)(1)
  disable 1
  continue
end

break *0x1003e470
commands 2
  silent
  set $parenthesis_context = *(unsigned int *)($esp + 4)
  set $parenthesis_value = *(unsigned int *)($esp + 8)
  printf "PAREN_CONTEXT_SEEDED context=%#x value=%u field_plus8=%u\n", $parenthesis_context, $parenthesis_value, *(unsigned int *)($parenthesis_context + 8)
  awatch *(unsigned int *)($parenthesis_context + 8)
  commands $bpnum
    silent
    printf "PAREN_CONTEXT_ACCESS pc=%#x eax=%#x ecx=%#x edx=%#x value=%u\n", $eip, $eax, $ecx, $edx, *(unsigned int *)($parenthesis_context + 8)
    bt 8
    continue
  end
  disable 2
  continue
end

continue
