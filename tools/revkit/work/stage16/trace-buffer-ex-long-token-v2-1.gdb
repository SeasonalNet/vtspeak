set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $size = 16384
  set $text = (char *)malloc($size + 1)
  set $i = 0
  while $i < $size
    set {char}($text + $i) = 97
    set $i = $i + 1
  end
  set {char}($text + $size) = 0
  set $buffer = (char *)malloc(60000)
  set $length = (int *)malloc(4)
  set *$length = 60000
  set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(1, $text, $buffer, $length, 0, 0, 1, 0, 0, -1, -1, -1, -1, 0, 0)
  printf "LONG_TOKEN_16K_START selector=1 result=%d length=%d\n", $result, *$length
  set $poll = 1
  while ($result == 0) && ($poll < 256)
    set *$length = 60000
    set $result = ((int (*)(int, char *, char *, int *, int, int, unsigned int, unsigned int *, unsigned int *, int, int, int, int, int, int))0x1001ddf0)(1, $text, $buffer, $length, 1, 0, 1, 0, 0, -1, -1, -1, -1, 0, 0)
    printf "LONG_TOKEN_16K_POLL selector=1 poll=%d result=%d length=%d\n", $poll, $result, *$length
    set $poll = $poll + 1
  end
  printf "LONG_TOKEN_16K_DONE selector=1 result=%d polls=%d\n", $result, $poll
  continue
end

continue
