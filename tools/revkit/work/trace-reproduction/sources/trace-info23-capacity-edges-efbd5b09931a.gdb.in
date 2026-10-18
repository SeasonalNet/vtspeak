set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $dst = (unsigned char *)malloc(64)
  set $sizes = (int *)malloc(8 * sizeof(int))
  set $sizes[0] = 2
  set $sizes[1] = 3
  set $sizes[2] = 4
  set $sizes[3] = 5
  set $sizes[4] = 6
  set $sizes[5] = 7
  set $sizes[6] = 8
  set $sizes[7] = 9
  set $i = 0
  while $i < 8
    set $j = 0
    while $j < 64
      set $dst[$j] = 0x5a
      set $j = $j + 1
    end
    set $result = ((int (*)(int, char *, void *, int))0x1002a690)(23, (char *)0, $dst, $sizes[$i])
    printf "INFO23_CAPACITY size=%d result=%d bytes=", $sizes[$i], $result
    x/16bx $dst
    set $i = $i + 1
  end
  kill
  quit
end

continue
