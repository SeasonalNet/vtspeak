set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $dst = (unsigned char *)malloc(1100)
  set $sizes = (int *)malloc(3 * sizeof(int))
  set $sizes[0] = 1023
  set $sizes[1] = 1024
  set $sizes[2] = 1100
  set $i = 0
  while $i < 3
    set $j = 0
    while $j < 1100
      set $dst[$j] = 0x5a
      set $j = $j + 1
    end
    set $result = ((int (*)(int, char *, void *, int))0x1002a690)(23, (char *)0, $dst, $sizes[$i])
    printf "INFO23_LONG size=%d result=%d first=", $sizes[$i], $result
    x/8bx $dst
    printf "INFO23_LONG_TAIL size=%d bytes=", $sizes[$i]
    x/4bx ($dst + 1021)
    set $i = $i + 1
  end
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(23, (char *)0, (void *)0, 1100)
  printf "INFO23_LONG_NULL_DEST result=%d\n", $result
  kill
  quit
end

continue
