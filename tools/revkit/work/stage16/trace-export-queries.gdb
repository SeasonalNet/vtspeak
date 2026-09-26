set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $name = (char *)malloc(256)
  set $path = (char *)malloc(256)
  set $version = (char *)0x0
  set $slot = 0
  while $slot < 6
    printf "SPEAKER_NAME slot=%d value=%s\n", $slot, ((char *(*)(int))0x10028230)($slot)
    set *(char *)$name = 0
    set *(char *)$path = 0
    set $speakers_result = ((int (*)(int, char *, char *))0x1002aa60)($slot, $name, $path)
    printf "SPEAKERS_INFO slot=%d result=%d name=%s path=%s\n", $slot, $speakers_result, $name, $path
    set $slot = $slot + 1
  end
  printf "SPEAKER_NAME invalid=-1 value=%s\n", ((char *(*)(int))0x10028230)(-1)
  set $version = (char *)((char *(*)(void))0x10028260)()
  printf "DEFAULT_VERSION value=%s\n", $version
  set $slot = 0
  set $dbsize_ptr = (int *)malloc(4)
  while $slot < 6
    set *$dbsize_ptr = 0
    set $dbsize_result = ((int (*)(int *, int))0x10028130)($dbsize_ptr, $slot)
    printf "DB_SIZE slot=%d result=%d value=%d\n", $slot, $dbsize_result, *$dbsize_ptr
    set $slot = $slot + 1
  end
  continue
end

continue
