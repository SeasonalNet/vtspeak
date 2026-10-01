set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10025e12
commands
  silent
  printf "LIPSYNC_OPEN_HELPER_RETURN eax=%#x context=%#x\n", $eax, $ebx
  x/s *(unsigned int *)($ebp+8)
  x/s *(unsigned int *)($ebp+12)
  continue
end

break *0x1001e050
commands
  silent
  printf "LIPSYNC_CONTEXT_FIELD wrapper=%#x open_context=%#x path=", $ebx, *(unsigned int *)($ebx+0x10)
  x/s ($ebp-0x200)
  printf "LIPSYNC_CONTEXT_MODE address=%#x\n", 0x1007d534
  x/s 0x1007d534
  continue
end

break *0x1001e0c0
commands
  silent
  disable 3
  set $state = *(unsigned int *)($esp + 4)
  set $logctx = *(unsigned int *)($state + 0x2c)
  set $writer = *(unsigned int *)($logctx + 0x10)
  printf "LIPSYNC_LOG_CONTEXT state=%#x object=%#x writer=%#x\n", $state, $logctx, $writer
  x/5wx $logctx
  continue
end

break *0x10025e4d
commands
  silent
  printf "LIPSYNC_NULL_WRITER eip=%#x edx=%#x ecx=%#x eax=%#x esp=%#x\n", $eip, $edx, $ecx, $eax, $esp
  x/i $eip
  bt
  disable 4
  kill
  quit
end

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  disable 5
  set $name = (char *)malloc(128)
  set {char[128]}$name = "."
  set {int}($esp + 4) = $text
  set {int}($esp + 8) = $name
  set {int}($esp + 12) = 1
  set {int}($esp + 16) = -1
  set {int}($esp + 20) = -1
  set {int}($esp + 24) = -1
  set {int}($esp + 28) = -1
  set {int}($esp + 32) = -1
  set {int}($esp + 36) = -1
  set $eip = 0x1001de50
  printf "LIPSYNC_REDIRECT directory_target ret=%#x filename=<%s> text=<%s>\n", $ret, $name, $text
  continue
end

break *0x1006a56e
commands
  silent
  printf "LIPSYNC_CREATEFILE_RESULT handle=%#x path=", $eax
  x/s *(unsigned int *)($ebp+8)
  continue
end

break *0x1006a57a
commands
  silent
  printf "LIPSYNC_GETLASTERROR code=%u path=", $eax
  x/s *(unsigned int *)($ebp+8)
  continue
end

break *0x10064cd1
commands
  silent
  printf "LIPSYNC_STREAM_OPEN_HELPER path="
  x/s *(unsigned int *)($esp+4)
  printf "LIPSYNC_STREAM_OPEN_MODE="
  x/s *(unsigned int *)($esp+8)
  continue
end

break *0x10068005
commands
  silent
  printf "LIPSYNC_MODE_TRANSLATOR path="
  x/s *(unsigned int *)($esp+4)
  printf "LIPSYNC_MODE_TRANSLATOR_MODE="
  x/s *(unsigned int *)($esp+8)
  continue
end

break *0x1006a3b5
commands
  silent
  printf "LIPSYNC_LOWLEVEL_OPEN path="
  x/s *(unsigned int *)($esp+4)
  printf "LIPSYNC_LOWLEVEL_OPEN_MODE flags=%#x disposition=%#x\n", *(unsigned int *)($esp+8), *(unsigned int *)($esp+12)
  continue
end

break *0x10068138
commands
  silent
  printf "LIPSYNC_STREAM_OPEN_DISPATCH flags=%#x pathname=", $ecx
  x/s *(unsigned int *)($ebp+8)
  printf "LIPSYNC_STREAM_OPEN_DISPATCH_MODE="
  x/s *(unsigned int *)($ebp+12)
  continue
end

break *0x1006a532
commands
  silent
  printf "LIPSYNC_DESCRIPTOR_ALLOC result=%#x path=", $eax
  x/s *(unsigned int *)($ebp+8)
  continue
end

continue
