set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $out = (unsigned int *)((char *(*)(unsigned int))0x1001d9c0)(4)
  set $id = -128
  while $id <= 256
    set *$out = 0xa5a5a5a5
    set $result = ((int (*)(int, char *, void *, int))0x1002a690)($id, (char *)0, (void *)$out, 4)
    printf "INFO_ID request=%d result=%d output=%#x\n", $id, $result, *$out
    set $id = $id + 1
  end
  continue
end

continue
