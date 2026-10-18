set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/kate-key5-class-pools-hello/class-pools-gdb.log
set logging overwrite on
set logging enabled on
set $context = -1
set $state = 0
set $model = 0

# Capture the class shortlist for every position before unit expansion.
break *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $model = *(unsigned int *)($esp + 12)
  if $context >= 0 && $context < 12
    set $record = $state + 0xae894 + $context * 0xfc
    set $class_count = *(short *)$record
    set $unit_counts = *(unsigned int *)($model + 0x8c)
    set $unit_lists = *(unsigned int *)($model + 0x94)
    set $class_keys = *(unsigned int *)($model + 0x90)
    printf "KATE_OLDROW_CLASSES context=%d class_count=%d", $context, $class_count
    set $i = 0
    while $i < $class_count && $i < 64
      set $class = *(unsigned int *)($record + 4 + $i * 4)
      set $member_count = *(unsigned short *)($unit_counts + $class * 2)
      set $members = *(unsigned int *)($unit_lists + $class * 4)
      printf " {class=%u key=%02x%02x%02x%02x%02x members=%u old_rows=", $class, *(unsigned char *)($class_keys + $class * 5), *(unsigned char *)($class_keys + $class * 5 + 1), *(unsigned char *)($class_keys + $class * 5 + 2), *(unsigned char *)($class_keys + $class * 5 + 3), *(unsigned char *)($class_keys + $class * 5 + 4), $member_count
      set $j = 0
      set $found = 0
      while $j < $member_count && $j < 50000
        set $unit = *(unsigned int *)($members + $j * 4)
        if $unit == 272822 || $unit == 272823 || $unit == 272824 || $unit == 272825
          printf "%s%u", $found == 0 ? "" : ",", $unit
          set $found = $found + 1
        end
        set $j = $j + 1
      end
      if $found == 0
        printf "-"
      end
      if $j < $member_count
        printf "+truncated"
      end
      printf "}"
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

# Record whether those rows survive unit expansion and continuity filtering.
break *0x100235f6
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $state = *(unsigned int *)($ebp + 0xc)
  if $context >= 0 && $context < 12
    set $count = *(short *)($state + 0xae988 + $context * 0xfc)
    set $i = 0
    set $found = 0
    while $i < $count && $i < 10000
      set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
      set $unit = *(unsigned int *)($node + 8)
      if $unit == 272822 || $unit == 272823 || $unit == 272824 || $unit == 272825
        printf "KATE_OLDROW_CANDIDATE context=%d rank=%u id=%u span=%d\n", $context, $i, $unit, *(short *)($node + 0xe)
        set $found = $found + 1
      end
      set $i = $i + 1
    end
    if $found == 0
      printf "KATE_OLDROW_CANDIDATES_NONE context=%d count=%d\n", $context, $count
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "KATE_OLDROW_CLASS_TRACE_READY\n"
  continue
end

continue
