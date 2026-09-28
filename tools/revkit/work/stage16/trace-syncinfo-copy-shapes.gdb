set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $src = ((int *(*)(void))0x10026360)()
  set $dst = ((int *(*)(void))0x10026360)()
  call ((void (*)(int *))0x10026450)($src)
  call ((void (*)(int *))0x10026450)($dst)
  set $src_rows = (unsigned char *)$src[0]
  set $dst_rows = (unsigned char *)$dst[0]
  set $src_row0 = $src_rows
  set $src_row1 = $src_rows + 0x24
  set $dst_row0 = $dst_rows
  set $dst_row1 = $dst_rows + 0x24
  set $src_nested0 = *(unsigned char **)($src_row0 + 4)
  set $src_nested1 = *(unsigned char **)($src_row1 + 4)
  set $dst_nested0 = *(unsigned char **)($dst_row0 + 4)
  set $dst_nested1 = *(unsigned char **)($dst_row1 + 4)

  set $src[1] = 2
  set $src[2] = 3
  set $dst[1] = 1
  set $dst[2] = 1
  set {unsigned short}$src_row0 = 0x1111
  set {unsigned int}($src_row0 + 8) = 0x11111111
  set {unsigned short}$src_row1 = 0x2222
  set {unsigned int}($src_row1 + 8) = 0x22222222
  set {unsigned int}($src_nested0 + 2 * 8) = 0x33333333
  set {unsigned short}($src_nested0 + 2 * 8 + 4) = 0x3333
  set {unsigned int}($src_nested1 + 2 * 8) = 0x44444444
  set {unsigned short}($src_nested1 + 2 * 8 + 4) = 0x4444
  set {unsigned short}$dst_row1 = 0xaaaa
  set {unsigned int}($dst_row1 + 8) = 0xaaaaaaaa
  set {unsigned int}($dst_nested1 + 2 * 8) = 0xbbbbbbbb
  set {unsigned short}($dst_nested1 + 2 * 8 + 4) = 0xbbbb
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_DST_SMALLER dest_before=1/1 source=2/3 dest_after=%d/%d row0=%04x/%08x row1=%04x/%08x entry2_row0=%08x/%04x entry2_row1=%08x/%04x pointers_distinct=%d/%d\n", $dst[1], $dst[2], *(unsigned short *)$dst_row0, *(unsigned int *)($dst_row0 + 8), *(unsigned short *)$dst_row1, *(unsigned int *)($dst_row1 + 8), *(unsigned int *)($dst_nested0 + 2 * 8), *(unsigned short *)($dst_nested0 + 2 * 8 + 4), *(unsigned int *)($dst_nested1 + 2 * 8), *(unsigned short *)($dst_nested1 + 2 * 8 + 4), $dst_nested0 != $src_nested0, $dst_nested1 != $src_nested1

  call ((void (*)(int *))0x10026450)($dst)
  set $src[1] = -1
  set $src[2] = 3
  set {unsigned short}$src_row0 = 0x5151
  set {unsigned short}$dst_row0 = 0x6161
  set {unsigned int}($dst_row0 + 8) = 0x61616161
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_NEGATIVE_ROWS source=-1/3 dest=%d/%d row0_untouched=%04x/%08x\n", $dst[1], $dst[2], *(unsigned short *)$dst_row0, *(unsigned int *)($dst_row0 + 8)

  call ((void (*)(int *))0x10026450)($dst)
  set $src[1] = 2
  set $src[2] = -1
  set {unsigned short}$src_row0 = 0x7171
  set {unsigned int}($src_row0 + 8) = 0x71717171
  set {unsigned short}$src_row1 = 0x7272
  set {unsigned int}($src_row1 + 8) = 0x72727272
  set {unsigned short}$dst_row0 = 0x8181
  set {unsigned int}($dst_row0 + 8) = 0x81818181
  set {unsigned short}$dst_row1 = 0x8282
  set {unsigned int}($dst_row1 + 8) = 0x82828282
  set {unsigned int}($dst_nested0) = 0x91919191
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_NEGATIVE_WIDTH source=2/-1 dest=%d/%d row0=%04x/%08x row1=%04x/%08x nested_untouched=%08x\n", $dst[1], $dst[2], *(unsigned short *)$dst_row0, *(unsigned int *)($dst_row0 + 8), *(unsigned short *)$dst_row1, *(unsigned int *)($dst_row1 + 8), *(unsigned int *)$dst_nested0

  set $src[1] = 600
  set $src[2] = 65
  set $dst[1] = 600
  set $dst[2] = 65
  call ((void (*)(int *))0x100263f0)($src)
  call ((void (*)(int *))0x100263f0)($dst)
  printf "SYNC_COPY_SHAPES_DONE\n"
  kill
  quit
end

continue
