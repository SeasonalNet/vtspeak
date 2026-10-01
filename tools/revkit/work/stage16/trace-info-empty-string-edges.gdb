set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $info = (char *)malloc(16)
  set $sizes = (int *)malloc(4 * sizeof(int))
  set $sizes[0] = -1
  set $sizes[1] = 0
  set $sizes[2] = 1
  set $sizes[3] = 2
  set $i = 0
  while $i < 4
    set *(unsigned int *)$info = 0x5a5a5a5a
    set *(unsigned int *)($info + 4) = 0x5a5a5a5a
    set $result = ((int (*)(int, char *, void *, int))0x1002a690)(23, (char *)0, (void *)$info, $sizes[$i])
    printf "INFO23 size=%d result=%d bytes=%#x,%#x\n", $sizes[$i], $result, *(unsigned int *)$info, *(unsigned int *)($info + 4)
    set $i = $i + 1
  end
  kill
  quit
end

continue
