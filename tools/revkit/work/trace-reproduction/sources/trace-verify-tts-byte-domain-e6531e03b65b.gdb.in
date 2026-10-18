set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $p = ((char *(*)(unsigned int))0x1001d9c0)(2)
  set $byte = 0
  while $byte < 256
    set {unsigned char}$p = $byte
    set {unsigned char}($p + 1) = 0
    set $r = ((int (*)(unsigned char *, int, int, int))0x1001ded0)($p, 1, -1, 0)
    printf "VERIFY_TTS_BYTE byte=%u ax=%u eax=%d\n", $byte, $r & 65535, $r
    set $byte = $byte + 1
  end
  continue
end

continue
