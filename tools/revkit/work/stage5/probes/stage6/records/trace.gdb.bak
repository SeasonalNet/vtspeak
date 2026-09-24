set pagination off
set confirm off
set logging file /work/stage5/probes/stage6/records/record-capture.log
set logging overwrite on
set logging enabled on
set $record_call = 0
break *0x1000d450
commands
  silent
  set $record_call = $record_call + 1
  set $dst = *(unsigned int *)($esp + 4)
  set $token = *(unsigned int *)($esp + 8)
  set $dictionary = *(unsigned int *)($esp + 12)
  set $arg4 = *(int *)($esp + 16)
  set $arg5 = *(int *)($esp + 20)
  set $arg8 = *(int *)($esp + 32)
  set $ret = *(unsigned int *)$esp
  set $base = $dst - 2
  printf "PHONE_RECORD_ENTRY call=%d dst=%#x base=%#x token=%#x arg4=%d arg5=%d arg8=%d dict=%#x ret=%#x\n", $record_call, $dst, $base, $token, $arg4, $arg5, $arg8, $dictionary, $ret
  x/s $token
  x/32bx $dictionary
  tbreak *$ret
  commands
    silent
    printf "PHONE_RECORD_RETURN call=%d count=%u token_index=%u status=%#x boundary_marker=%#x source=", $record_call, *(unsigned short *)$dst, *(unsigned short *)($base + 4), *(unsigned char *)($base + 6), *(unsigned char *)($base + 0x25)
    x/s ($base + 7)
    printf "PHONE_RECORD_HEADER bytes:\n"
    x/40bx ($base + 2)
    set $phone_i = 0
    while $phone_i < *(unsigned short *)$dst && $phone_i < 32
      printf "PHONE_SLOT index=%u marker=%#x bytes:\n", $phone_i, *(unsigned char *)($base + 0x25 + $phone_i)
      x/65bx ($base + 0x39 + $phone_i * 0x41)
      set $phone_i = $phone_i + 1
    end
    printf "PHONE_RECORD_TAIL words:\n"
    x/8hx ($base + 0x54e)
    continue
  end
  continue
end

break *0x10060e50
commands
  silent
  set $numout = *(unsigned int *)($esp + 4)
  set $numin = *(unsigned int *)($esp + 8)
  set $numflags = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "DECIMAL_ENTRY input="
  x/s $numin
  printf " flags=%d\n", (short)$numflags
  tbreak *$ret
  commands
    silent
    printf "DECIMAL_RETURN output="
    x/s $numout
    continue
  end
  continue
end

break *0x1005fa60
commands
  silent
  set $cardout = *(unsigned int *)($esp + 4)
  set $cardnum = *(int *)($esp + 8)
  set $cardmode = *(int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "CARDINAL_ENTRY number=%d mode=%c current=", $cardnum, $cardmode
  x/s $cardout
  tbreak *$ret
  commands
    silent
    printf "CARDINAL_RETURN number=%d mode=%c output=", $cardnum, $cardmode
    x/s $cardout
    continue
  end
  continue
end

break *0x10039d50
commands
  silent
  set $yearout = *(unsigned int *)($esp + 4)
  set $yearin = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "YEAR_ENTRY input=%#x\n", $yearin
  x/s $yearin
  tbreak *$ret
  commands
    silent
    printf "YEAR_RETURN output=%#x\n", $yearout
    x/s $yearout
    continue
  end
  continue
end

continue
