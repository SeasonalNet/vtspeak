set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-selector-byte0-old-context2-candidate-boundary-2026-09-29/legacy-context2-candidate-boundary-gdb.log
set logging overwrite on
set logging enabled on
set $legacy_position = -1
set $legacy_state = 0

# Trace one context through class expansion, continuity coverage, full-span
# filtering, local scoring, and final unit-list handoff.
hbreak *0x1001cff0
commands
  silent
  set $legacy_position = *(int *)($esp + 4)
  set $legacy_state = *(unsigned int *)($esp + 8)
  if $legacy_position == 2
    set $record = $legacy_state + 0xbb154 + $legacy_position * 0xfc
    set $class_count = *(unsigned short *)$record
    printf "APPLE_OLD_CONTEXT2_CLASS_INPUT count=%u classes:", $class_count
    set $i = 0
    while $i < $class_count && $i < 40
      printf " %u", *(unsigned int *)($record + 4 + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    set $return = *(unsigned int *)$esp
    thbreak *$return
    commands
      silent
      set $record = $legacy_state + 0xbb154 + $legacy_position * 0xfc
      set $count = *(unsigned short *)($record + 0xf4)
      printf "APPLE_OLD_CONTEXT2_FINAL count=%u ids:", $count
      set $i = 0
      while $i < $count && $i < 40
        printf " %u", *(unsigned int *)($record + 0x7c + $i * 4)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

# Immediately after FUN_1001cd60 attaches coverage, before FUN_1001cff0's
# full-span gate or its local-score cutoff.
hbreak *0x1001cd60
commands
  silent
  if $legacy_position == 2
    set $return = *(unsigned int *)$esp
    thbreak *$return
    commands
      silent
      set $record = $legacy_state + 0xbb154 + $legacy_position * 0xfc
      set $count = *(unsigned short *)($record + 0xf4)
      set $array = $legacy_state + 0x528fc
      printf "APPLE_OLD_CONTEXT2_AFTER_COVERAGE count=%u target188826_index=", $count
      set $i = 0
      set $target_index = -1
      while $i < $count && $i < 10000
        set $node = *(unsigned int *)($array + $i * 4)
        if *(unsigned int *)($node + 8) == 188826
          set $target_index = $i
        end
        set $i = $i + 1
      end
      printf "%d\n", $target_index
      set $i = 0
      while $i < $count && $i < 60
        set $node = *(unsigned int *)($array + $i * 4)
        printf "APPLE_OLD_CONTEXT2_COVERAGE_ROW index=%u id=%u span=%u left=%u right=%u local=%g\n", $i, *(unsigned int *)($node + 8), *(unsigned short *)($node + 0x10), *(unsigned short *)($node + 0xc), *(unsigned short *)($node + 0xe), *(float *)$node
        set $i = $i + 1
      end
      continue
    end
  end
  continue
end

# Capture the post-filter candidate count and ranked top-30 output when the
# legacy unit ranker enters its tail sorter.
hbreak *0x1001d48e
commands
  silent
  if $legacy_position == 2
    set $record = $legacy_state + 0xbb154 + $legacy_position * 0xfc
    set $count = *(unsigned short *)($record + 0xf4)
    set $prefix = $eax
    printf "APPLE_OLD_CONTEXT2_SORT_INPUT total=%u presorted_prefix=%u\n", $count, $prefix
    set $i = 0
    while $i < $count && $i < 10000
      set $node = *(unsigned int *)($legacy_state + 0x528fc + $i * 4)
      if *(unsigned int *)($node + 8) == 188826
        printf "APPLE_OLD_CONTEXT2_SORT_TARGET index=%u score=%g rank_key=%g span=%u\n", $i, *(float *)$node, *(float *)($node + 4), *(unsigned short *)($node + 0x10)
      end
      set $i = $i + 1
    end
  end
  continue
end

hbreak *0x1001d493
commands
  silent
  if $legacy_position == 2
    set $i = 0
    printf "APPLE_OLD_CONTEXT2_SORTED_TOP30"
    while $i < 30
      set $node = *(unsigned int *)($legacy_state + 0x528fc + $i * 4)
      printf " [%u]=%u(score=%g,span=%u)", $i, *(unsigned int *)($node + 8), *(float *)$node, *(unsigned short *)($node + 0x10)
      set $i = $i + 1
    end
    printf "\n"
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_CONTEXT2_CANDIDATE_TRACE_READY\n"
  continue
end

continue
