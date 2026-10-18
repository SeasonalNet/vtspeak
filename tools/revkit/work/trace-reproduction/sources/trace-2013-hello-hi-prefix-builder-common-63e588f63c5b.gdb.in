set $helper_calls = 0
set $prefilter_calls = 0

# The class shortlist entering FUN_10023350 and the units each class expands to.
break *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  if $context == 0 || $context == 1 || $context == 6
    set $state = *(unsigned int *)($esp + 8)
    set $model = *(unsigned int *)($esp + 16)
    set $record = $state + 0xae894 + $context * 0xfc
    set $class_count = *(short *)$record
    set $unit_counts = *(unsigned int *)($model + 0x8c)
    set $unit_lists = *(unsigned int *)($model + 0x94)
    printf "HELLO_BUILDER_CLASSES context=%d count=%d", $context, $class_count
    set $i = 0
    while $i < $class_count && $i < 64
      set $class = *(unsigned int *)($record + 4 + $i * 4)
      set $member_count = *(unsigned short *)($unit_counts + $class * 2)
      set $members = *(unsigned int *)($unit_lists + $class * 4)
      printf " {class=%u members=%u units=", $class, $member_count
      set $j = 0
      while $j < $member_count && $j < 64
        printf "%s%u", $j == 0 ? "" : ",", *(unsigned int *)($members + $j * 4)
        set $j = $j + 1
      end
      printf "}"
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

# FUN_100230a0 adds coverage metadata and sorts the candidate-node array.
break *0x100230a0
commands
  silent
  set $context = *(int *)($esp + 4)
  if $context == 0 || $context == 1 || $context == 6
    set $helper_calls = $helper_calls + 1
    set $state = *(unsigned int *)($esp + 8)
    set $record = $state + 0xae894 + $context * 0xfc
    set $count = *(unsigned short *)($record + 0xf4)
    printf "HELLO_COVERAGE_ENTRY context=%d count=%d", $context, $count
    set $i = 0
    while $i < $count && $i < 12
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      printf " [%d]=%u", $i, *(unsigned int *)($node + 8)
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

# This is the observed pre-local-score checkpoint used by the existing list trace.
break *0x100235f6
commands
  silent
  set $context = *(int *)($ebp + 8)
  if $context == 0 || $context == 1 || $context == 6
    set $prefilter_calls = $prefilter_calls + 1
    set $state = *(unsigned int *)($ebp + 0xc)
    set $count = *(short *)($state + 0xae988 + $context * 0xfc)
    printf "HELLO_PRE_SCORE context=%d count=%d", $context, $count
    set $i = 0
    while $i < $count && $i < 12
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      printf " [%d]=%u(span=%d,weighted=%d)", $i, *(unsigned int *)($node + 8), *(short *)($node + 0xe), *(short *)($node + 0x10)
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HELLO_HI_PREFIX_BUILDER_TRACE_READY\n"
  continue
end

continue
