set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-repeat-2006-shortlist-sort-boundary/legacy-shortlist-sort-gdb.log
set logging overwrite on
set logging enabled on
set $legacy_current_position = -1
set $legacy_current_state = 0

# FUN_1001cff0 builds and scores each legacy position. Save the live position
# and state so the scorer and retained sort can be keyed to row 5.
hbreak *0x1001cff0
commands
  silent
  set $legacy_current_position = *(int *)($esp + 4)
  set $legacy_current_state = *(unsigned int *)($esp + 8)
  if $legacy_current_position == 5
    printf "LEGACY_SHORTLIST_ENTRY position=5 class_count=%u\n", *(unsigned short *)($legacy_current_state + 0xbb154 + 5 * 0xfc)
  end
  continue
end

# FUN_1001c860 returns the old engine's per-unit local cost in ST(0), then
# FUN_1001cff0 stores it at node +0.
hbreak *0x1001d3f4
commands
  silent
  if $legacy_current_position == 5
    set $node = *(unsigned int *)$esi
    set $unit = *(unsigned int *)($node + 8)
    if $unit == 264071 || $unit == 49874
      printf "LEGACY_UNIT_SCORE position=5 unit=%u score=%g rank_key=%g span_left=%d span_right=%d span=%d\n", $unit, $st0, *(float *)($node + 4), *(short *)($node + 0xc), *(short *)($node + 0xe), *(short *)($node + 0x10)
    end
  end
  continue
end

# Capture the pre-sort set and its sorted 30-row prefix at FUN_100140c0.
hbreak *0x1001d48e
commands
  silent
  if $legacy_current_position == 5
    set $state = $legacy_current_state
    set $record = $state + 0xbb154 + 5 * 0xfc
    set $count = *(short *)($record + 0xf4)
    set $prefix = $eax
    printf "LEGACY_SHORTLIST_SORT_INPUT position=5 total=%d presorted_prefix=%d suffix_count=%d\n", $count, $prefix, $count - $prefix
    set $i = 0
    while $i < $count && $i < 100
      set $node = *(unsigned int *)($state + 0x528fc + $i * 4)
      set $unit = *(unsigned int *)($node + 8)
      if $unit == 264071 || $unit == 49874
        printf "LEGACY_SHORTLIST_SORT_TARGET position=5 index=%d unit=%u score=%g key=%g flags=%#x\n", $i, $unit, *(float *)$node, *(float *)($node + 4), *(unsigned char *)($node + 0x12)
      end
      set $i = $i + 1
    end
    tbreak *0x1001d493
    commands
      silent
      set $i = 0
      printf "LEGACY_SHORTLIST_SORTED_TOP30 position=5"
      while $i < 30
        set $node = *(unsigned int *)($state + 0x528fc + $i * 4)
        printf " [%d]=%u(score=%g,key=%g,span=%d)", $i, *(unsigned int *)($node + 8), *(float *)$node, *(float *)($node + 4), *(short *)($node + 0x10)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_SHORTLIST_SORT_TRACE_READY\n"
  continue
end

continue
