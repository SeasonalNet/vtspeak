set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $name = (char *)malloc(512)
  set $path = (char *)malloc(512)
  set $k = 1
  while $k < 8
    set $slot = 0
    while $slot < 6
      set $selector = (int)((unsigned int)$slot + (unsigned int)$k * 0x20000000)
      set $key = ((char *(*)(int))0x1001c690)($selector)
      printf "PATHKEY_WRAP k=%d slot=%d selector=%d value=%s\n", $k, $slot, $selector, $key
      set $info = ((int (*)(int, char *, char *))0x1002aa60)($selector, $name, $path)
      printf "SPEAKERSINFO_WRAP k=%d slot=%d selector=%d result=%d name=%s path=%s\n", $k, $slot, $selector, $info, $name, $path
      set $slot = $slot + 1
    end
    set $k = $k + 1
  end

  kill
  quit
end

continue
