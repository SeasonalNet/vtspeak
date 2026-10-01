set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $virtual_protect = *(void **)0x0040b180
  set $old_protect = (unsigned int *)malloc(4)
  set *$old_protect = 0

  set $name_allocation = (unsigned int)malloc(65536)
  set $name_guard = ($name_allocation & 0xfffff000) + 0x8000
  set $name_protect = ((int (*)(void *, unsigned int, unsigned int, unsigned int *))$virtual_protect)((void *)$name_guard, 0x1000, 1, $old_protect)
  set $path_allocation = (unsigned int)malloc(65536)
  set $path_guard = ($path_allocation & 0xfffff000) + 0x8000
  set $path_protect = ((int (*)(void *, unsigned int, unsigned int, unsigned int *))$virtual_protect)((void *)$path_guard, 0x1000, 1, $old_protect)
  printf "SPEAKERSINFO_GUARDS name=%#x/%d path=%#x/%d old_protect=%#x\n", $name_guard, $name_protect, $path_guard, $path_protect, *$old_protect
  if $name_protect != 1 || $path_protect != 1
    printf "SPEAKERSINFO_GUARD_SETUP_FAILED\n"
    kill
    quit
  end

  set $name = (char *)($name_guard - 6)
  set $path = (char *)($path_guard - 24)
  set $result = ((int (*)(int, char *, char *))0x1002aa60)(3, $name, $path)
  printf "SPEAKERSINFO_GUARD_RESULT case=path-exact selector=3 result=%d name=%s path=%s\\n", $result, $name, $path
  kill
  quit
end

continue

