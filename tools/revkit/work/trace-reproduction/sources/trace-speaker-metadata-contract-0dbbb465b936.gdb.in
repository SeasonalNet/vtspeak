set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1

  set $path_copy = (char *)malloc(64)
  set $path_slot = 0
  while $path_slot < 6
    set $path_key = ((char *(*)(int))0x1001c690)($path_slot)
    if $path_slot == 0
      set $first_key_pointer = $path_key
      set $path_i = 0
      while $path_i < 64
        set {unsigned char}($path_copy + $path_i) = *(unsigned char *)($path_key + $path_i)
        if *(unsigned char *)($path_key + $path_i) == 0
          set $path_i = 64
        else
          set $path_i = $path_i + 1
        end
      end
    end
    set $last_key_pointer = $path_key
    printf "PATH_KEY slot=%d pointer=%p value=%s\n", $path_slot, $path_key, $path_key
    set $path_slot = $path_slot + 1
  end
  printf "PATH_KEY_BUFFER same_pointer=%d saved_first=%s current_after_sweep=%s\n", $first_key_pointer == $last_key_pointer, $path_copy, $last_key_pointer
  set $path_key_extreme = ((char *(*)(int))0x1001c690)(-2147483648)
  printf "PATH_KEY_EDGE input=INT_MIN value=%s\n", $path_key_extreme

  set $name_buffer = (char *)malloc(512)
  set $path_buffer = (char *)malloc(512)
  set $slot = 0
  while $slot < 6
    set $i = 0
    while $i < 512
      set {unsigned char}($name_buffer + $i) = 0xa5
      set {unsigned char}($path_buffer + $i) = 0xa5
      set $i = $i + 1
    end
    set $info_result = ((int (*)(int, char *, char *))0x1002aa60)($slot, $name_buffer, $path_buffer)
    set $name_length = (int)strlen($name_buffer)
    set $path_length = (int)strlen($path_buffer)
    set $name_tail_bad = 0
    set $path_tail_bad = 0
    set $i = $name_length + 1
    while $i < 512
      if *(unsigned char *)($name_buffer + $i) != 0xa5
        set $name_tail_bad = $name_tail_bad + 1
      end
      set $i = $i + 1
    end
    set $i = $path_length + 1
    while $i < 512
      if *(unsigned char *)($path_buffer + $i) != 0xa5
        set $path_tail_bad = $path_tail_bad + 1
      end
      set $i = $i + 1
    end
    printf "SPEAKERS_INFO_CONTRACT slot=%d result=%d name_len=%d name_nul=%d name_next=0x%02x name_tail_bad=%d path_len=%d path_nul=%d path_next=0x%02x path_tail_bad=%d name=%s path=%s\n", $slot, $info_result, $name_length, *(unsigned char *)($name_buffer + $name_length), *(unsigned char *)($name_buffer + $name_length + 1), $name_tail_bad, $path_length, *(unsigned char *)($path_buffer + $path_length), *(unsigned char *)($path_buffer + $path_length + 1), $path_tail_bad, $name_buffer, $path_buffer
    set $slot = $slot + 1
  end
  set $info_extreme = ((int (*)(int, char *, char *))0x1002aa60)(-2147483648, $name_buffer, $path_buffer)
  printf "SPEAKERS_INFO_EDGE input=INT_MIN result=%d name=%s path=%s\n", $info_extreme, $name_buffer, $path_buffer
  detach
  quit
end

continue
