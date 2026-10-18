set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-2013-context2-candidate-boundary-2026-09-29/adapted-context2-candidate-boundary-gdb.log
set logging overwrite on
set logging enabled on
set $context = -1
set $state = 0
set $model = 0

# Save context state at FUN_10023350 entry.
hbreak *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $model = *(unsigned int *)($esp + 12)
  continue
end

# This point follows coverage ordering and full-span filtering, before
# duration/local-score evaluation and the 30-row cap.
hbreak *0x10023814
commands
  silent
  if $context == 2
    set $count = *(unsigned short *)($state + 0xae988 + $context * 0xfc)
    set $index = 0
    set $target_index = -1
    printf "APPLE_2013_CONTEXT2_AFTER_COVERAGE count=%u target188826_index=", $count
    while $index < $count && $index < 10000
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if *(unsigned int *)($node + 8) == 188826
        set $target_index = $index
      end
      set $index = $index + 1
    end
    printf "%d\n", $target_index
    set $index = 0
    while $index < $count && $index < 12
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      printf "APPLE_2013_CONTEXT2_COVERAGE_ROW rank=%u unit=%u span=%d weighted=%d right_matches=%d pre_score_key=%g\n", $index, *(unsigned int *)($node + 8), *(short *)($node + 0xe), *(short *)($node + 0x10), *(short *)($node + 0xc), *(float *)($node + 4)
      set $index = $index + 1
    end
  end
  continue
end

# FUN_100182e0 result after scoring; target row and nearby old-selected rows.
hbreak *0x1002389a
commands
  silent
  if (*(int *)($ebp + 8) == 2)
    set $node = *(unsigned int *)$esi
    set $unit = *(unsigned int *)($node + 8)
    if ($unit == 188826) || ($unit == 100180) || ($unit == 64719)
      printf "APPLE_2013_CONTEXT2_LOCAL_SCORE unit=%u score=%g coverage_key=%g span=%d weighted=%d flags=%#x\n", $unit, $st0, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), *(unsigned char *)($node + 0x12)
    end
  end
  continue
end

# Capture the suffix handed to the sort and the actual post-sort shortlist.
hbreak *0x10023948
commands
  silent
  if (*(int *)($ebp + 8) == 2)
    set $count = $ebx
    set $prefix = *(int *)($ebp - 0x20)
    set $sort_count = *(int *)$esp
    set $sort_nodes = *(unsigned int *)($esp + 4)
    printf "APPLE_2013_CONTEXT2_SORT_INPUT total=%u prefix=%u suffix=%u\n", $count, $prefix, $sort_count
    set $index = 0
    while $index < $count && $index < 10000
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      if *(unsigned int *)($node + 8) == 188826
        printf "APPLE_2013_CONTEXT2_SORT_TARGET index=%u score=%g coverage_key=%g\n", $index, *(float *)$node, *(float *)($node + 4)
      end
      set $index = $index + 1
    end
    tbreak *0x1002394d
    commands
      silent
      set $index = 0
      set $target_rank = -1
      printf "APPLE_2013_CONTEXT2_SORTED_TOP30"
      while $index < 30
        set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
        set $unit = *(unsigned int *)($node + 8)
        if $unit == 188826
          set $target_rank = $index
        end
        printf " [%u]=%u(score=%g,key=%g,span=%d,weighted=%d)", $index, $unit, *(float *)$node, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10)
        set $index = $index + 1
      end
      printf " target188826_rank=%d\n", $target_rank
      continue
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT2_CANDIDATE_TRACE_READY\n"
  continue
end

continue
