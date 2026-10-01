set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/repeat-hello-old-query-pools/legacy-query-pools-gdb.log
set logging overwrite on
set logging enabled on
set $old_split_calls = 0

hbreak *0x1001dfb0
commands
  silent
  set $old_split_calls = $old_split_calls + 1
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $context = *(unsigned int *)($esp + 12)
  set $return = *(unsigned int *)$esp
  printf "OLD_REPEAT_SPLIT_ENTER n=%u slot=%u state=%#x context=%#x\n", $old_split_calls, $slot, $state, $context
  thbreak *$return
  commands
    silent
    printf "OLD_REPEAT_SPLIT_RETURN n=%u slot=%u accepted=%u\n", $old_split_calls, $slot, $al
    continue
  end
  continue
end

hbreak *0x1001da80
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
    set $before = *(unsigned short *)($state + 0xf8e64)
    printf "OLD_REPEAT_CLASS_LOOKUP slot=%u mode=%u key:", $slot, $mode
    x/5ub $key
    thbreak *$return
    commands
      silent
      set $count = *(unsigned short *)($state + 0xf8e64)
      set $feature_base = *(unsigned int *)($model + 0x98)
      set $weight_base = *(unsigned int *)($model + 0x8c)
      set $added = $count - $before
      printf "OLD_REPEAT_CLASS_RETURN slot=%u mode=%u appended_metric=%u total=%u model=%#x feature_base=%#x ids:", $slot, $mode, $eax, $count, $model, $feature_base
      set $i = $before
      while $i < $count && $i < 256
        set $id = *(unsigned int *)($state + 0xf89b4 + $i * 4)
        printf " %u", $id
        printf "\nOLD_REPEAT_CLASS_MEMBER slot=%u id=%u weight=%u feature:", $slot, $id, *(unsigned short *)($weight_base + $id * 2)
        x/7ub ($feature_base + $id * 7)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

hbreak *0x1001e110
commands
  silent
  printf "OLD_REPEAT_FALLBACK_ENTER slot=%u\n", *(unsigned int *)($esp + 4)
  continue
end

hbreak *0x408187
commands
  silent
  printf "OLD_REPEAT_QUERY_TRACE_READY\n"
  continue
end

continue
