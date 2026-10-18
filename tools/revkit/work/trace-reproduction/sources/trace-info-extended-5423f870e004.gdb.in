set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $info = (char *)malloc(256)
  set $ids = (int *)malloc(7 * sizeof(int))
  set $ids[0] = -1
  set $ids[1] = 27
  set $ids[2] = 28
  set $ids[3] = 100
  set $ids[4] = 101
  set $ids[5] = 102
  set $ids[6] = 2147483647
  set $i = 0
  while $i < 7
    set *(unsigned int *)$info = 0xffffffff
    set $result = ((int (*)(int, char *, void *, int))0x1002a690)($ids[$i], (char *)0, (void *)$info, 256)
    printf "INFO_EXT request=%d result=%d value=%d\n", $ids[$i], $result, *(int *)$info
    set $i = $i + 1
  end
  continue
end

continue
