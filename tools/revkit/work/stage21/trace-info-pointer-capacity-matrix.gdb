set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $out = ((char *(*)(unsigned int))0x1001d9c0)(32)
  set $i = 0
  while $i < 32
    set {unsigned char}($out + $i) = 165
    set $i = $i + 1
  end
  set $ids = (int *)((char *(*)(unsigned int))0x1001d9c0)(7 * sizeof(int))
  set $ids[0] = -1
  set $ids[1] = 1
  set $ids[2] = 2
  set $ids[3] = 3
  set $ids[4] = 23
  set $ids[5] = 27
  set $ids[6] = 101
  set $i = 0
  while $i < 7
    set $result = ((int (*)(int, char *, void *, int))0x1002a690)($ids[$i], (char *)0, (void *)0, 0)
    printf "INFO_PTR null_output request=%d size=0 result=%d\n", $ids[$i], $result
    set $i = $i + 1
  end

  set $i = 0
  while $i < 32
    set {unsigned char}($out + $i) = 165
    set $i = $i + 1
  end
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(-1, (char *)0, (void *)$out, 4)
  printf "INFO_PTR invalid request=-1 size=4 result=%d first=%#x\n", $result, *(unsigned int *)$out
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(27, (char *)0, (void *)$out, 4)
  printf "INFO_PTR invalid request=27 size=4 result=%d first=%#x\n", $result, *(unsigned int *)$out

  set $i = 0
  while $i < 32
    set {unsigned char}($out + $i) = 165
    set $i = $i + 1
  end
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(0, (char *)0, (void *)$out, 11)
  printf "INFO_PTR build_date size=11 result=%d bytes=%#x,%#x,%#x,%#x\n", $result, *(unsigned char *)$out, *(unsigned char *)($out + 1), *(unsigned char *)($out + 2), *(unsigned char *)($out + 3)
  set $i = 0
  while $i < 32
    set {unsigned char}($out + $i) = 165
    set $i = $i + 1
  end
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(0, (char *)0, (void *)$out, 12)
  printf "INFO_PTR build_date size=12 result=%d prefix=%c%c%c\n", $result, *(char *)$out, *(char *)($out + 1), *(char *)($out + 2)

  set $sizes = (int *)((char *(*)(unsigned int))0x1001d9c0)(5 * sizeof(int))
  set $sizes[0] = -1
  set $sizes[1] = 0
  set $sizes[2] = 1
  set $sizes[3] = 3
  set $sizes[4] = 4
  set $i = 0
  while $i < 5
    set $j = 0
    while $j < 32
      set {unsigned char}($out + $j) = 165
      set $j = $j + 1
    end
    set $result = ((int (*)(int, char *, void *, int))0x1002a690)(3, (char *)0, (void *)$out, $sizes[$i])
    printf "INFO_PTR db_directory size=%d result=%d bytes=%#x,%#x,%#x,%#x\n", $sizes[$i], $result, *(unsigned char *)$out, *(unsigned char *)($out + 1), *(unsigned char *)($out + 2), *(unsigned char *)($out + 3)
    set $i = $i + 1
  end

  set $i = 0
  while $i < 32
    set {unsigned char}($out + $i) = 165
    set $i = $i + 1
  end
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(23, (char *)0, (void *)$out, 1)
  printf "INFO_PTR db_build_date_missing size=1 result=%d first=%#x\n", $result, *(unsigned int *)$out
  set $result = ((int (*)(int, char *, void *, int))0x1002a690)(101, (char *)0, (void *)$out, -1)
  printf "INFO_PTR playback_state size=-1 result=%d first=%#x\n", $result, *(unsigned int *)$out
  continue
end

continue
