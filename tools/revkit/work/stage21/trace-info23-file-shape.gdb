set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $dst = (unsigned char *)malloc(64)
  set $j = 0
  while $j < 64
    set $dst[$j] = 0x5a
    set $j = $j + 1
  end
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(23, (char *)0, $dst, 64)
  printf "INFO23_SHAPE result=%d bytes=", $result
  x/24bx $dst
  kill
  quit
end

continue
