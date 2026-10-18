set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $init = ((int *(*)(void))0x10026360)()
  set $init_rows = (unsigned char *)$init[0]
  set $init_row0 = $init_rows
  set $init_row1 = $init_rows + 0x24
  set $init_rowlast = $init_rows + 599 * 0x24
  set $init[3] = 0x33333333
  set $init[4] = 0x44444444
  set $init[5] = 0x55555555
  set $init[6] = 0x66666666
  set $init[7] = 0x77777777
  set $init[8] = 0x88888888
  set $init[9] = 0x99999999
  set $init[10] = 0xaaaaaaaa
  set $init[11] = 0xbbbbbbbb
  set $init[12] = 0xcccccccc
  set $init[13] = 0xdddddddd
  set {unsigned short}($init_row0) = 0xa0a0
  set {unsigned short}($init_row1) = 0xa1a1
  set {unsigned short}($init_rowlast) = 0xafaf
  set $i = 0
  while $i < 7
    set {unsigned int}($init_row0 + 8 + $i * 4) = 0xa0a0a0a0
    set {unsigned int}($init_row1 + 8 + $i * 4) = 0xa1a1a1a1
    set {unsigned int}($init_rowlast + 8 + $i * 4) = 0xafafafaf
    set $i = $i + 1
  end
  set $init_nested0 = *(unsigned char **)($init_row0 + 4)
  set $init_nested1 = *(unsigned char **)($init_row1 + 4)
  set $init_nestedlast = *(unsigned char **)($init_rowlast + 4)
  set {unsigned int}($init_nested0) = 0xa0a0a0a0
  set {unsigned short}($init_nested0 + 4) = 0xa0a0
  set {unsigned int}($init_nested0 + 64 * 8) = 0xa0a0a0a0
  set {unsigned short}($init_nested0 + 64 * 8 + 4) = 0xa0a0
  set {unsigned int}($init_nested1) = 0xa1a1a1a1
  set {unsigned short}($init_nested1 + 4) = 0xa1a1
  set {unsigned int}($init_nestedlast) = 0xafafafaf
  set {unsigned short}($init_nestedlast + 4) = 0xafaf

  call ((void (*)(int *))0x10026450)($init)
  printf "SYNC_INIT_HEADER count=%d width=%d field3=%08x fields4_13=", $init[1], $init[2], $init[3]
  set $i = 4
  while $i < 14
    printf "%08x,", $init[$i]
    set $i = $i + 1
  end
  printf "\n"
  printf "SYNC_INIT_ROWS row0_word=%04x row0_first=%08x/%04x row0_last=%08x/%04x row1_word=%04x rowlast_word=%04x\n", *(unsigned short *)$init_row0, *(unsigned int *)$init_nested0, *(unsigned short *)($init_nested0 + 4), *(unsigned int *)($init_nested0 + 64 * 8), *(unsigned short *)($init_nested0 + 64 * 8 + 4), *(unsigned short *)$init_row1, *(unsigned short *)$init_rowlast
  printf "SYNC_INIT_NESTED_PTRS row0_same=%d row1_same=%d row599_same=%d\n", $init_nested0 == *(unsigned char **)($init_row0 + 4), $init_nested1 == *(unsigned char **)($init_row1 + 4), $init_nestedlast == *(unsigned char **)($init_rowlast + 4)
  printf "SYNC_INIT_ROW_WORDS row0=%04x,", *(unsigned short *)$init_row0
  set $i = 0
  while $i < 7
    printf "%08x,", *(unsigned int *)($init_row0 + 8 + $i * 4)
    set $i = $i + 1
  end
  printf " row1=%04x,", *(unsigned short *)$init_row1
  set $i = 0
  while $i < 7
    printf "%08x,", *(unsigned int *)($init_row1 + 8 + $i * 4)
    set $i = $i + 1
  end
  printf " row599=%04x,", *(unsigned short *)$init_rowlast
  set $i = 0
  while $i < 7
    printf "%08x,", *(unsigned int *)($init_rowlast + 8 + $i * 4)
    set $i = $i + 1
  end
  printf "\n"
  call ((void (*)(int *))0x100263f0)($init)

  set $src = ((int *(*)(void))0x10026360)()
  set $dst = ((int *(*)(void))0x10026360)()
  set $src_rows = (unsigned char *)$src[0]
  set $dst_rows = (unsigned char *)$dst[0]
  set $src[3] = 0x03030303
  set $dst[3] = 0xd3d3d3d3
  set $i = 4
  while $i < 14
    set $src[$i] = 0x10000000 + $i
    set $dst[$i] = 0xd0000000 + $i
    set $i = $i + 1
  end
  set $row = 0
  while $row < 600
    set $src_row = $src_rows + $row * 0x24
    set {unsigned short}($src_row) = 0x1200 + $row
    set $i = 0
    while $i < 7
      set {unsigned int}($src_row + 8 + $i * 4) = 0x11000000 + $row * 0x100 + $i
      set $i = $i + 1
    end
    set $src_nested = *(unsigned char **)($src_row + 4)
    set $item = 0
    while $item < 65
      set {unsigned int}($src_nested + $item * 8) = 0x21000000 + $row * 0x100 + $item
      set {unsigned short}($src_nested + $item * 8 + 4) = 0x2200 + $row * 0x10 + $item
      set $item = $item + 1
    end
    set $row = $row + 1
  end
  set $src_last_row = $src_rows + 599 * 0x24
  set {unsigned short}($src_last_row) = 0x17ff
  set $i = 0
  while $i < 7
    set {unsigned int}($src_last_row + 8 + $i * 4) = 0x17ff0000 + $i
    set $i = $i + 1
  end
  set $src_last_nested = *(unsigned char **)($src_last_row + 4)
  set {unsigned int}($src_last_nested + 64 * 8) = 0x27ff0000
  set {unsigned short}($src_last_nested + 64 * 8 + 4) = 0x27ff
  call ((void (*)(int *, int *))0x10026500)($dst, $src)
  printf "SYNC_COPY_HEADER rows=%d/%d width=%d/%d row_arrays_distinct=%d fields3_13=", $src[1], $dst[1], $src[2], $dst[2], $src[0] != $dst[0]
  set $i = 3
  while $i < 14
    printf "%08x/%08x,", $src[$i], $dst[$i]
    set $i = $i + 1
  end
  printf "\n"
  set $row = 0
  while $row < 2
    set $src_row = $src_rows + $row * 0x24
    set $dst_row = $dst_rows + $row * 0x24
    printf "SYNC_COPY_ROW row=%d word=%04x/%04x nested_arrays_distinct=%d fields8_32=", $row, *(unsigned short *)$src_row, *(unsigned short *)$dst_row, *(unsigned int *)($src_row + 4) != *(unsigned int *)($dst_row + 4)
    set $i = 0
    while $i < 7
      printf "%08x/%08x,", *(unsigned int *)($src_row + 8 + $i * 4), *(unsigned int *)($dst_row + 8 + $i * 4)
      set $i = $i + 1
    end
    printf " nested="
    set $item = 0
    while $item < 3
      set $src_nested = *(unsigned char **)($src_row + 4)
      set $dst_nested = *(unsigned char **)($dst_row + 4)
      printf "%08x/%08x:%04x/%04x,", *(unsigned int *)($src_nested + $item * 8), *(unsigned int *)($dst_nested + $item * 8), *(unsigned short *)($src_nested + $item * 8 + 4), *(unsigned short *)($dst_nested + $item * 8 + 4)
      set $item = $item + 1
    end
    printf "\n"
    set $row = $row + 1
  end
  set $last_row = $src_rows + 599 * 0x24
  set $last_dst_row = $dst_rows + 599 * 0x24
  set $last_nested = *(unsigned char **)($last_row + 4)
  set $last_dst_nested = *(unsigned char **)($last_dst_row + 4)
  set $row_mismatches = 0
  set $nested_mismatches = 0
  set $pointer_mismatches = 0
  set $row = 0
  while $row < 600
    set $src_row = $src_rows + $row * 0x24
    set $dst_row = $dst_rows + $row * 0x24
    if *(unsigned short *)$src_row != *(unsigned short *)$dst_row
      set $row_mismatches = $row_mismatches + 1
    end
    set $i = 0
    while $i < 7
      if *(unsigned int *)($src_row + 8 + $i * 4) != *(unsigned int *)($dst_row + 8 + $i * 4)
        set $row_mismatches = $row_mismatches + 1
      end
      set $i = $i + 1
    end
    set $src_nested = *(unsigned char **)($src_row + 4)
    set $dst_nested = *(unsigned char **)($dst_row + 4)
    if $src_nested == $dst_nested
      set $pointer_mismatches = $pointer_mismatches + 1
    end
    set $item = 0
    while $item < 65
      if *(unsigned int *)($src_nested + $item * 8) != *(unsigned int *)($dst_nested + $item * 8)
        set $nested_mismatches = $nested_mismatches + 1
      end
      if *(unsigned short *)($src_nested + $item * 8 + 4) != *(unsigned short *)($dst_nested + $item * 8 + 4)
        set $nested_mismatches = $nested_mismatches + 1
      end
      set $item = $item + 1
    end
    set $row = $row + 1
  end
  printf "SYNC_COPY_FULL rows=%d width=%d row_value_mismatches=%d nested_value_mismatches=%d shared_nested_pointers=%d last_row=%04x/%04x last_field=%08x/%08x last_entry=%08x/%08x:%04x/%04x\n", $src[1], $src[2], $row_mismatches, $nested_mismatches, $pointer_mismatches, *(unsigned short *)$last_row, *(unsigned short *)$last_dst_row, *(unsigned int *)($last_row + 32), *(unsigned int *)($last_dst_row + 32), *(unsigned int *)($last_nested + 64 * 8), *(unsigned int *)($last_dst_nested + 64 * 8), *(unsigned short *)($last_nested + 64 * 8 + 4), *(unsigned short *)($last_dst_nested + 64 * 8 + 4)

  set $src[1] = 600
  set $src[2] = 65
  set $dst[1] = 600
  set $dst[2] = 65
  call ((void (*)(int *))0x100263f0)($src)
  call ((void (*)(int *))0x100263f0)($dst)
  continue
end

continue
