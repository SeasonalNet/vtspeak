set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  call ((void (*)(int *, int *))0x10026500)(0, 0)
  set $src = ((int *(*)(void))0x10026360)()
  call ((void (*)(int *, int *))0x10026500)($src, 0)
  call ((void (*)(int *, int *))0x10026500)(0, $src)
  printf "SYNC_COPY_NULL_CALLS completed=1\n"

  set $self = ((int *(*)(void))0x10026360)()
  call ((void (*)(int *))0x10026450)($self)
  set $self_rows = (unsigned char *)$self[0]
  set $self_nested = *(unsigned char **)($self_rows + 4)
  set $self[3] = 0x12345678
  set $self[13] = 0xabcdef01
  set {unsigned short}$self_rows = 0x1234
  set {unsigned int}($self_rows + 8) = 0x23456789
  set {unsigned int}$self_nested = 0x3456789a
  set {unsigned short}($self_nested + 4) = 0x4567
  call ((void (*)(int *, int *))0x10026500)($self, $self)
  printf "SYNC_COPY_SELF rows=%d width=%d field3=%08x field13=%08x row0=%04x/%08x nested0=%08x/%04x row_array=%#x nested_array=%#x\n", $self[1], $self[2], $self[3], $self[13], *(unsigned short *)$self_rows, *(unsigned int *)($self_rows + 8), *(unsigned int *)$self_nested, *(unsigned short *)($self_nested + 4), $self[0], *(unsigned int *)($self_rows + 4)

  set $dst = ((int *(*)(void))0x10026360)()
  call ((void (*)(int *))0x10026450)($src)
  call ((void (*)(int *))0x10026450)($dst)
  set $src_rows = (unsigned char *)$src[0]
  set $dst_rows = (unsigned char *)$dst[0]
  set $dst_nested0 = *(unsigned char **)($dst_rows + 4)
  set $src[1] = 2
  set $src[2] = 3
  set $src[3] = 0x33333333
  set $src[4] = 0x44444444
  set {unsigned short}$src_rows = 0x1010
  set {unsigned int}($src_rows + 8) = 0x11111111
  set $src_row1 = $src_rows + 0x24
  set {unsigned short}$src_row1 = 0x2020
  set {unsigned int}($src_row1 + 8) = 0x22222222
  set $src_nested0 = *(unsigned char **)($src_rows + 4)
  set $src_nested1 = *(unsigned char **)($src_row1 + 4)
  set {unsigned int}$src_nested0 = 0x12121212
  set {unsigned short}($src_nested0 + 4) = 0x1212
  set {unsigned int}($src_nested0 + 2 * 8) = 0x13131313
  set {unsigned short}($src_nested0 + 2 * 8 + 4) = 0x1313
  set {unsigned int}$src_nested1 = 0x21212121
  set {unsigned short}($src_nested1 + 4) = 0x2121
  set {unsigned int}$dst_nested0 = 0x61616161
  set {unsigned short}($dst_nested0 + 4) = 0x6161
  set {unsigned int}($dst_nested0 + 3 * 8) = 0x63636363
  set {unsigned short}($dst_nested0 + 3 * 8 + 4) = 0x6363
  set {unsigned short}($dst_rows + 2 * 0x24) = 0x7272
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_SMALL_SOURCE dst_dims=%d/%d fields3_4=%08x/%08x row0=%04x/%08x row1=%04x/%08x row2_untouched=%04x entry0=%08x/%04x entry2=%08x/%04x entry3_untouched=%08x/%04x nested_distinct=%d\n", $dst[1], $dst[2], $dst[3], $dst[4], *(unsigned short *)$dst_rows, *(unsigned int *)($dst_rows + 8), *(unsigned short *)($dst_rows + 0x24), *(unsigned int *)($dst_rows + 0x24 + 8), *(unsigned short *)($dst_rows + 2 * 0x24), *(unsigned int *)$dst_nested0, *(unsigned short *)($dst_nested0 + 4), *(unsigned int *)($dst_nested0 + 2 * 8), *(unsigned short *)($dst_nested0 + 2 * 8 + 4), *(unsigned int *)($dst_nested0 + 3 * 8), *(unsigned short *)($dst_nested0 + 3 * 8 + 4), $dst_nested0 != $src_nested0

  set $dst[1] = 600
  set $dst[2] = 65
  call ((void (*)(int *))0x10026450)($dst)
  set $src[1] = 0
  set $src[2] = 3
  set $src[3] = 0x51515151
  set {unsigned short}$dst_rows = 0x7474
  set {unsigned short}$src_rows = 0x5151
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_ZERO_ROWS dst_dims=%d/%d field3=%08x row0_untouched=%04x\n", $dst[1], $dst[2], $dst[3], *(unsigned short *)$dst_rows

  set $dst[1] = 600
  set $dst[2] = 65
  call ((void (*)(int *))0x10026450)($dst)
  set $src[1] = 2
  set $src[2] = 0
  set {unsigned short}$src_rows = 0x5252
  set {unsigned int}($src_rows + 8) = 0x52525252
  set {unsigned short}$dst_rows = 0x7575
  set {unsigned int}($dst_rows + 8) = 0x75757575
  set {unsigned int}$src_nested0 = 0x53535353
  set {unsigned int}$dst_nested0 = 0x76767676
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_ZERO_WIDTH dst_dims=%d/%d row0=%04x/%08x nested0_untouched=%08x\n", $dst[1], $dst[2], *(unsigned short *)$dst_rows, *(unsigned int *)($dst_rows + 8), *(unsigned int *)$dst_nested0

  set $src[1] = 600
  set $src[2] = 65
  set $dst[1] = 600
  set $dst[2] = 65
  call ((void (*)(int *))0x100263f0)($src)
  call ((void (*)(int *))0x100263f0)($dst)
  call ((void (*)(int *))0x100263f0)($self)
  continue
end

continue
