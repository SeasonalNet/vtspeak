set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $csv = ((char *(*)(unsigned int))0x1001d9c0)(11)
  set {char}($csv + 0) = 104
  set {char}($csv + 1) = 101
  set {char}($csv + 2) = 108
  set {char}($csv + 3) = 108
  set {char}($csv + 4) = 111
  set {char}($csv + 5) = 44
  set {char}($csv + 6) = 72
  set {char}($csv + 7) = 72
  set {char}($csv + 8) = 44
  set {char}($csv + 9) = 80
  set {char}($csv + 10) = 0
  set $init_before = *(int *)0x100a0458
  printf "USERDICT_EXT init_before=0x%x data_len=10\\n", $init_before
  set $load = ((short (*)(int, char *, char *, int))0x10027880)(200, 0, $csv, 10)
  printf "USERDICT_EXT memory_load index=200 path=NULL length=10 ax=%d\\n", $load
  set $duplicate = ((short (*)(int, char *, char *, int))0x10027880)(200, 0, $csv, 10)
  printf "USERDICT_EXT duplicate_load index=200 ax=%d\\n", $duplicate
  set $unload = ((short (*)(int))0x10027980)(200)
  printf "USERDICT_EXT unload index=200 ax=%d\\n", $unload
  set {int}0x100a0458 = 0
  set $uninitialized = ((short (*)(int, char *, char *, int))0x10027880)(201, 0, $csv, 10)
  printf "USERDICT_EXT forced_uninitialized index=201 ax=%d\\n", $uninitialized
  set $uninitialized_unload = ((short (*)(int))0x10027980)(201)
  printf "USERDICT_EXT forced_uninitialized_unload index=201 ax=%d\\n", $uninitialized_unload
  set {int}0x100a0458 = $init_before
  printf "USERDICT_EXT init_restored=0x%x\\n", *(int *)0x100a0458
  continue
end

continue
