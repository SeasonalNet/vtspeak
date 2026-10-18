set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  if $context == 4 || $context == 5
    set $state = *(unsigned int *)($esp + 8)
    set $model = *(unsigned int *)($esp + 12)
    set $record = $state + 0xae894 + $context * 0xfc
    set $class_count = *(short *)$record
    set $unit_counts = *(unsigned int *)($model + 0x8c)
    set $unit_lists = *(unsigned int *)($model + 0x94)
    printf "BRIDGET_CONTEXT_CLASSES context=%d count=%d", $context, $class_count
    set $i = 0
    while $i < $class_count && $i < 64
    set $class = *(unsigned int *)($record + 4 + $i * 4)
    set $member_count = *(unsigned short *)($unit_counts + $class * 2)
    set $members = *(unsigned int *)($unit_lists + $class * 4)
      set $class_keys = *(unsigned int *)($model + 0x90)
      printf " {class=%u key=%02x%02x%02x%02x%02x members=%u targets=", $class, *(unsigned char *)($class_keys + $class * 5), *(unsigned char *)($class_keys + $class * 5 + 1), *(unsigned char *)($class_keys + $class * 5 + 2), *(unsigned char *)($class_keys + $class * 5 + 3), *(unsigned char *)($class_keys + $class * 5 + 4), $member_count
      set $j = 0
      set $target_count = 0
      while $j < $member_count && $j < 10000
        set $unit = *(unsigned int *)($members + $j * 4)
        if $unit == 150652 || $unit == 365789 || $unit == 382197 || $unit == 548226
          printf "%u ", $unit
          set $target_count = $target_count + 1
        end
        set $j = $j + 1
      end
      if $target_count == 0
        printf "-"
      end
      printf "}"
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "BRIDGET_CONTEXT_CLASSES_TRACE_READY\n"
  continue
end

continue
