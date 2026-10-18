set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10025e4d
commands
  silent
  printf "LIPSYNC_NULL_WRITER eip=%#x edx=%#x ecx=%#x eax=%#x esp=%#x\n", $eip, $edx, $ecx, $eax, $esp
  x/i $eip
  bt
  disable 1
  kill
  quit
end

break *0x1001e0c0
commands
  silent
  disable 2
  set $state = *(unsigned int *)($esp + 4)
  set $logctx = *(unsigned int *)($state + 0x2c)
  set $writer = *(unsigned int *)($logctx + 0x10)
  printf "LIPSYNC_LOG_CONTEXT state=%#x object=%#x writer=%#x\n", $state, $logctx, $writer
  x/5wx $logctx
  continue
end

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  disable 3
  set $name = (char *)malloc(128)
  set {char[128]}$name = "no-such-lipsync-dir/report.txt"
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
  printf "LIPSYNC_REDIRECT missing_parent ret=%#x filename=<%s> text=<%s>\n", $ret, $name, $text
  continue
end

continue
