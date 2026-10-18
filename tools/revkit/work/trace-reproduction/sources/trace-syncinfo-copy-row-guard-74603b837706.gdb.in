set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV stop print pass
set $copy_started = 0

break *0x1001da50
commands 1
  silent
  disable 1
  set $copy_started = 1
  set $src = ((int *(*)(void))0x10026360)()
  set $dst = ((int *(*)(void))0x10026360)()
  call ((void (*)(int *))0x10026450)($src)
  call ((void (*)(int *))0x10026450)($dst)
  set $src_rows = (unsigned char *)$src[0]
  set $src_row1 = $src_rows + 0x24
  set $src_nested0 = *(unsigned char **)($src_rows + 4)
  set $src_nested1 = *(unsigned char **)($src_row1 + 4)
  set $alloc = ((char *(*)(void *, unsigned int, unsigned int, unsigned int))0x7b68d8a0)((char *)0, 0x3000, 0x3000, 0x04)
  set $old_protect = (unsigned int *)malloc(4)
  set $protect_ok = ((int (*)(void *, unsigned int, unsigned int, unsigned int *))0x7b68df20)($alloc + 0x1000, 0x1000, 0x01, $old_protect)
  set $rowbuf = $alloc + 0xfdc
  set $guard = $alloc + 0x1000
  set $src[1] = 2
  set $src[2] = 3
  set {unsigned int}($src_rows + 8) = 0x11111111
  set {unsigned int}($src_row1 + 8) = 0x22222222
  set {unsigned int}$src_nested0 = 0x33333333
  set {unsigned short}($src_nested0 + 4) = 0x3333
  set {unsigned int}$src_nested1 = 0x44444444
  set {unsigned short}($src_nested1 + 4) = 0x4444
  set {unsigned int}($rowbuf + 4) = *(unsigned int *)($dst[0] + 4)
  set $dst[0] = (int)$rowbuf
  printf "SYNC_COPY_ROW_GUARD_SETUP alloc=%#x rowbuf=%#x row_bytes=36 guard=%#x protect_ok=%d source_dims=2/3\n", $alloc, $rowbuf, $guard, $protect_ok
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_ROW_GUARD_RETURN unexpected=1 dest_dims=%d/%d\n", $dst[1], $dst[2]
  quit
end

catch signal SIGSEGV
commands 2
  silent
  if $copy_started
    printf "SYNC_COPY_ROW_GUARD_FAULT eip=%#x guard=%#x dest_dims=%d/%d row0_count=%d row0_value=%#x\n", $eip, $guard, $dst[1], $dst[2], *(unsigned short *)$rowbuf, *(unsigned int *)($rowbuf + 8)
  end
  continue
end

continue
