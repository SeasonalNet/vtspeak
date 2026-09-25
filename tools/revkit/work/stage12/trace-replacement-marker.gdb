set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Observe whether marker bytes are copied into the later TPP component records.
break *0x10034180
commands
  silent
  set $context = *(unsigned int *)($esp + 4)
  set $components = *(unsigned int *)($context + 0x39e4)
  set $count = *(unsigned int *)$components
  printf "TPP_COMPONENTS count=%u", $count
  set $i = 0
  while $i < $count && $i < 5
    set $entry = $components + $i * 0x140
    printf " component=%u marker=%#x surface=", $i, *(unsigned char *)($entry + 0x1a)
    x/s ($entry + 0x46)
    set $i = $i + 1
  end
  continue
end

# Verify that the special tokenizer fallback is reached for the input.
break *0x10043810
commands
  silent
  printf "REPLACEMENT_ROUTINE_ENTER p1=%#x p2=%#x p3=%#x p4=%#x\n", *(unsigned int *)($esp + 4), *(unsigned int *)($esp + 8), *(unsigned int *)($esp + 12), *(unsigned int *)($esp + 16)
  continue
end

# FUN_10043810's indexed replacement-table lookup at 0x10043a8a.
break *0x10043a8f
commands
  silent
  printf "REPLACEMENT_LOOKUP result=%d key=", $eax
  x/s ($ebp - 0x60)
  continue
end

# The successful indexed table path writes lowercase 'd' to local_a4.
break *0x10043b10
commands
  silent
  set $marker_index = *(unsigned int *)($ebp - 4)
  printf "REPLACEMENT_MARKER value=%#x index=%u key=", *(unsigned char *)($ebp - 0xa0 + $marker_index), $marker_index
  x/s ($ebp - 0x60)
  printf "REPLACEMENT_VALUE "
  x/s ($ebp - 0x1ac)
  continue
end

# Apartment/basement lookup is a separate branch and can emit uppercase 'A'.
break *0x10043d75
commands
  silent
  printf "ADDRESS_LOOKUP result=%d key=", $eax
  x/s ($ebp - 0x60)
  continue
end

break *0x10043e49
commands
  silent
  set $marker_index = *(unsigned int *)($ebp - 4)
  printf "ADDRESS_MARKER value=%#x index=%u key=", *(unsigned char *)($ebp - 0xa0 + $marker_index), $marker_index
  x/s ($ebp - 0x60)
  printf "ADDRESS_VALUE "
  x/s ($ebp - 0x1ac)
  continue
end

break *0x10043e87
commands
  silent
  set $marker_index = *(unsigned int *)($ebp - 4)
  printf "ADDRESS_MARKER value=%#x index=%u key=", *(unsigned char *)($ebp - 0xa0 + $marker_index), $marker_index
  x/s ($ebp - 0x60)
  printf "ADDRESS_VALUE "
  x/s ($ebp - 0x1ac)
  continue
end

break *0x10043ed0
commands
  silent
  set $marker_index = *(unsigned int *)($ebp - 4)
  printf "ADDRESS_MARKER value=%#x index=%u key=", *(unsigned char *)($ebp - 0xa0 + $marker_index), $marker_index
  x/s ($ebp - 0x60)
  printf "ADDRESS_VALUE "
  x/s ($ebp - 0x1ac)
  continue
end

continue
