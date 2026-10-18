set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-plain-features-v3-gdb.log
set logging overwrite on
set logging enabled on

break *0x408187
commands
  silent
  printf "APPLE_OLD_FEATURE_TRACE_READY\n"

break *0x1001da80
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x1001dfb0 && $caller < 0x1001e110
    set $key = *(unsigned int *)($esp + 4)
    set $slot = *(unsigned int *)($esp + 8)
    set $state = *(unsigned int *)($esp + 12)
    set $model = *(unsigned int *)($esp + 16)
    set $mode = *(unsigned short *)($esp + 20)
    set $return = $caller
    printf "APPLE_OLD_CLASS_LOOKUP slot=%u mode=%u key:", $slot, $mode
    x/5ub $key
    tbreak *$return
    commands
      silent
      set $count = *(unsigned short *)($state + 0xf8e64)
      set $features = *(unsigned int *)($model + 0x98)
      printf "APPLE_OLD_CLASS_RETURN slot=%u appended=%u total=%u ids:", $slot, $eax, $count
      set $i = 0
      while $i < $count && $i < 32
        set $id = *(unsigned int *)($state + 0xf89b4 + $i * 4)
        printf " %u", $id
        printf "\nAPPLE_OLD_CLASS_FEATURE id=%u:", $id
        x/7ub ($features + $id * 7)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

continue
end

continue
