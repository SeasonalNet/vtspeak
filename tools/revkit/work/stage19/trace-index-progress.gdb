set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x100199f0
commands
  silent
  set $index_state = *(unsigned int *)($esp + 4)
  printf "index dispatch old_flag=%u units=%u\n", *(unsigned char *)($index_state + 0x62), *(unsigned int *)($index_state + 0x4c)
  continue
end

break *0x100197d0
commands
  silent
  set $index_state = *(unsigned int *)($esp + 4)
  printf "legacy index reader entry units=%u\n", *(unsigned int *)($index_state + 0x4c)
  continue
end

break *0x10019940
commands
  silent
  printf "versioned reader entry\n"
  continue
end

break *0x100197f0
commands
  silent
  if $ebx % 50000 == 0
    printf "legacy-index unit=%u total=%u stream-cursor=%u\n", $ebx, *(unsigned int *)($edi + 0x4c), *(unsigned int *)($esi + 0x14)
  end
  continue
end

continue
