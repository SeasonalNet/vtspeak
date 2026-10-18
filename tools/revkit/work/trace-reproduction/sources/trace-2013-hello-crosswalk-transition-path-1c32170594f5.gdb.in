set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-plain-crosswalk-transition-path-v3/adapted-transition-path-gdb.log
set logging overwrite on
set logging enabled on

set $last_sum = -1
break *0x10024060
commands
  silent
  set $last_sum = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    if $al == 0 && $last_sum > 2
      printf "TRANSITION_PATH_CUTOFF_OVERRIDE slot=%u sum=%d\n", $slot, $last_sum
      set $eax = 1
    end
    continue
  end
  continue
end

break *0x10023060
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $return = $caller
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      continue
    end
  end
  continue
end

# FUN_10018c80's inner edge calculation. Capture old/new competing rows and
# the observed baseline path plus the alternative reached by the feature swap.
break *0x10023350
commands
  silent
  set $rank_context = *(int *)($esp + 4)
  set $rank_state = *(unsigned int *)($esp + 8)
  set $rank_return = *(unsigned int *)$esp
  if $rank_context < 4
    tbreak *$rank_return
    commands
      silent
      set $rank_record = $rank_state + 0xae894 + $rank_context * 0xfc
      set $rank_count = *(unsigned short *)($rank_record + 0xf4)
      printf "RANK_RETURN context=%d count=%u rows=", $rank_context, $rank_count
      set $i = 0
      while $i < $rank_count && $i < 80
        set $node = *(unsigned int *)($rank_state + 0x477ac + $i * 4)
        printf "%u(score=%g,rank=%g,span=%d,weighted=%d,flags=%#x) ", *(unsigned int *)($node + 8), *(float *)$node, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), *(unsigned char *)($node + 0x12)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x10018c80
commands
  silent
  set $context_entry = *(int *)($esp + 4)
  set $state_entry = *(unsigned int *)($esp + 8)
  if $context_entry > 0 && $context_entry < 4
    set $current_record = $state_entry + 0xae894 + $context_entry * 0xfc
    set $previous_record = $state_entry + 0xae894 + ($context_entry - 1) * 0xfc
    set $current_count = *(short *)($current_record + 0xf4)
    set $previous_count = *(short *)($previous_record + 0xf4)
    printf "TRANSITION_LISTS context=%d current_count=%d current=", $context_entry, $current_count
    set $i = 0
    while $i < $current_count && $i < 80
      printf "%u ", *(unsigned int *)($current_record + 0x7c + $i * 4)
      set $i = $i + 1
    end
    printf "previous_count=%d previous=", $previous_count
    set $i = 0
    while $i < $previous_count && $i < 80
      printf "%u ", *(unsigned int *)($previous_record + 0x7c + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if ($context >= 1 && $context <= 3) && ($current == 272823 || $current == 273370 || $current == 272824 || $current == 273371 || $current == 272825 || $current == 255282) && ($previous == 272822 || $previous == 52542 || $previous == 273369 || $previous == 179388 || $previous == 272823 || $previous == 273370 || $previous == 272824 || $previous == 273371)
    set $feature = *(float *)($ebp - 0x4c)
    set $divisor = *(int *)($ebp - 0x10)
    set $penalty = *(int *)($ebp - 0xc)
    set $previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
    set $previous_cost = *(float *)$previous_cost_ptr
    set $local_cost = *(float *)($ebp - 0x74)
    printf "TRANSITION_EDGE context=%d current=%u previous=%u feature=%g divisor=%d penalty=%d previous_path=%g local=%g total=%g\n", $context, $current, $previous, $feature, $divisor, $penalty, $previous_cost, $local_cost, $feature / $divisor + $penalty + $previous_cost + $local_cost
  end
  continue
end

set $selected_calls = 0
break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 16
    printf "TRANSITION_PATH_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "TRANSITION_PATH_TRACE_READY\n"
  continue
end

continue
