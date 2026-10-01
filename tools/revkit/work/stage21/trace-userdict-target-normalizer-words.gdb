set pagination off
set confirm off
set debuginfo enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $target_input = (unsigned int)malloc(256)
  set $target_words = 1
  while $target_words <= 12
    set $target_offset = 0
    while $target_offset < 256
      set {unsigned char}($target_input + $target_offset) = 0
      set $target_offset = $target_offset + 1
    end
    set $target_offset = 0
    while $target_offset < ($target_words * 2 - 1)
      if ($target_offset % 2) == 0
        set {unsigned char}($target_input + $target_offset) = 65
      else
        set {unsigned char}($target_input + $target_offset) = 32
      end
      set $target_offset = $target_offset + 1
    end
    set $target_result = ((int (*)(char *))0x1002a570)($target_input)
    printf "TARGET_NORM_WORDS words=%d result=%d input=", $target_words, $target_result
    x/s $target_input
    set $target_words = $target_words + 1
  end
  printf "TARGET_NORM_WORDS_COMPLETE first=1 last=12 count=12\n"
  kill
end

continue
