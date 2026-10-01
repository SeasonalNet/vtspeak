set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-2013-pair-transition/adapted-pair-transition-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0

# Capture the candidate IDs entering each adjacent-context pass, then read the
# cumulative costs and predecessor indices written by FUN_10018c80.
break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  if $context > 0 && $context < 4
    set $record = $state + 0xae894 + $context * 0xfc
    set $previous_record = $state + 0xae894 + ($context - 1) * 0xfc
    set $count = *(short *)($record + 0xf4)
    set $previous_count = *(short *)($previous_record + 0xf4)
    printf "HI_TRANSITION_BEGIN context=%d current_count=%d previous_count=%d current=", $context, $count, $previous_count
    set $i = 0
    while $i < $count && $i < 80
      printf "%u ", *(unsigned int *)($record + 0x7c + $i * 4)
      set $i = $i + 1
    end
    printf " previous="
    set $i = 0
    while $i < $previous_count && $i < 80
      printf "%u ", *(unsigned int *)($previous_record + 0x7c + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      set $record = $state + 0xae894 + $context * 0xfc
      set $count = *(short *)($record + 0xf4)
      printf "HI_TRANSITION_RESULT context=%d count=%d", $context, $count
      set $i = 0
      while $i < $count && $i < 30
        set $unit = *(unsigned int *)($record + 0x7c + $i * 4)
        set $slot = $context * 30 + $i
        set $cost = *(float *)($state + 0x8212c + $slot * 4)
        set $pred = *(short *)($state + 0x9f664 + $slot * 2)
        set $previous_unit = 0
        if $pred >= 0 && $pred < 80
          set $previous_unit = *(unsigned int *)($state + 0xae894 + ($context - 1) * 0xfc + 0x7c + $pred * 4)
        end
        if $unit == 272820 || $unit == 272821 || $unit == 65331 || $unit == 281932 || $previous_unit == 272820 || $previous_unit == 272821 || $previous_unit == 65331 || $previous_unit == 281932
          printf " [i=%d unit=%u cost=%g pred=%d pred_unit=%u]", $i, $unit, $cost, $pred, $previous_unit
        end
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

# Record individual edge terms for the shared source rows and their immediate
# successors; this shows whether transition additions break local-score ties.
hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if $context >= 1 && $context <= 3 && ($current == 272820 || $current == 272821 || $current == 65331 || $current == 281932 || $previous == 272820 || $previous == 272821 || $previous == 65331 || $previous == 281932)
    set $feature = *(float *)($ebp - 0x4c)
    set $divisor = *(int *)($ebp - 0x10)
    set $penalty = *(int *)($ebp - 0xc)
    set $previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
    set $previous_cost = *(float *)$previous_cost_ptr
    set $local_cost = *(float *)($ebp - 0x74)
    printf "HI_TRANSITION_EDGE context=%d current=%u previous=%u feature=%g divisor=%d penalty=%d previous_path=%g local=%g total=%g\n", $context, $current, $previous, $feature, $divisor, $penalty, $previous_cost, $local_cost, $feature / $divisor + $penalty + $previous_cost + $local_cost
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HI_TRANSITION_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_PAIR_TRANSITION_TRACE_READY\n"
  continue
end

continue
