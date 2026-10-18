set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $allocation = (unsigned int)malloc(65536)
  set $guard = ($allocation & 0xfffff000) + 0x8000
  set $old_protect = (unsigned int *)malloc(4)
  set *$old_protect = 0
  set $virtual_protect = *(void **)0x0040b180
  set $protect_result = ((int (*)(void *, unsigned int, unsigned int, unsigned int *))$virtual_protect)((void *)$guard, 0x1000, 1, $old_protect)
  printf "SPEAKERSINFO_GUARD_SETUP allocation=%#x guard=%#x protect_result=%d old_protect=%#x\n", $allocation, $guard, $protect_result, *$old_protect
  if $protect_result != 0
    set $restore_result = ((int (*)(void *, unsigned int, unsigned int, unsigned int *))$virtual_protect)((void *)$guard, 0x1000, *$old_protect, $old_protect)
    printf "SPEAKERSINFO_GUARD_RESTORE result=%d\n", $restore_result
  end
  kill
  quit
end

continue
