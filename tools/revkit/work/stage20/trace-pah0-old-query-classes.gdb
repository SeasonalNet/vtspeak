set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/pah0-old-query-classes-v2-gdb.log
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
  printf "OLD_PAH0_SPLIT_ENTER n=%u slot=%u state=%#x context=%#x\n", $old_split_calls, $slot, $state, $context
  thbreak *$return
  commands
    silent
    printf "OLD_PAH0_SPLIT_RETURN n=%u slot=%u accepted=%u\n", $old_split_calls, $slot, $al
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
    set $mode = *(unsigned short *)($esp + 20)
    set $return = $caller
    printf "OLD_PAH0_CLASS_LOOKUP slot=%u mode=%u key:", $slot, $mode
    x/5ub $key
    thbreak *$return
    commands
      silent
      set $count = *(unsigned short *)($state + 0xf8e64)
      printf "OLD_PAH0_CLASS_RETURN slot=%u mode=%u appended=%u total=%u ids:", $slot, $mode, $eax, $count
      set $i = 0
      while $i < $count && $i < 32
        printf " %u", *(unsigned int *)($state + 0xf89b4 + $i * 4)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

hbreak *0x1001d5c0
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x1001dfb0 && $caller < 0x1001e110
    set $ids = *(unsigned int *)($esp + 4)
    set $count = *(unsigned short *)($esp + 8)
    set $context = *(unsigned int *)($esp + 12)
    set $weights = *(unsigned int *)($context + 0x8c)
    set $return = $caller
    printf "OLD_PAH0_METRIC_ENTRY count=%u ids_weights:", $count
    set $i = 0
    while $i < $count && $i < 64
      set $id = *(unsigned int *)($ids + $i * 4)
      printf " %u:%u", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    thbreak *$return
    commands
      silent
      printf "OLD_PAH0_METRIC_RETURN sum=%u\n", $eax
      continue
    end
  end
  continue
end

hbreak *0x1001e110
commands
  silent
  printf "OLD_PAH0_FALLBACK_ENTER slot=%u\n", *(unsigned int *)($esp + 4)
  continue
end

hbreak *0x408187
commands
  silent
  printf "OLD_PAH0_QUERY_CLASS_TRACE_READY\n"
  continue
end

continue
