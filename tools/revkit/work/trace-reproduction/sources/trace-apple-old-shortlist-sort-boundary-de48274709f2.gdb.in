set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-selector-byte0-old-shortlist-2026-09-29/legacy-shortlist-sort-gdb.log
set logging overwrite on
set logging enabled on
set $legacy_current_position = -1
set $legacy_current_state = 0

hbreak *0x1001cff0
commands
  silent
  set $legacy_current_position = *(int *)($esp + 4)
  set $legacy_current_state = *(unsigned int *)($esp + 8)
  if $legacy_current_position == 4
    printf "APPLE_OLD_SHORTLIST_ENTRY position=4 candidate_count=%u\n", *(unsigned short *)($legacy_current_state + 0xbb154 + 4 * 0xfc + 0xf4)
  end
  continue
end

# The local scorer's result is stored at node +0; the rank key is node +4.
hbreak *0x1001d3f4
commands
  silent
  if $legacy_current_position == 4
    set $node = *(unsigned int *)$esi
    set $unit = *(unsigned int *)($node + 8)
    if $unit == 59559 || $unit == 82037 || $unit == 2555
      printf "APPLE_OLD_POSITION4_SCORE unit=%u score=%g rank_key=%g span_left=%d span_right=%d span=%d\n", $unit, $st0, *(float *)($node + 4), *(short *)($node + 0xc), *(short *)($node + 0xe), *(short *)($node + 0x10)
    end
  end
  continue
end

# Inspect the shortlist sorter input and the post-sort top 30 at position 4.
hbreak *0x1001d48e
commands
  silent
  if $legacy_current_position == 4
    set $state = $legacy_current_state
    set $record = $state + 0xbb154 + 4 * 0xfc
    set $count = *(short *)($record + 0xf4)
    set $prefix = $eax
    printf "APPLE_OLD_POSITION4_SORT_INPUT total=%d presorted_prefix=%d suffix_count=%d\n", $count, $prefix, $count - $prefix
    set $i = 0
    while $i < $count && $i < 100
      set $node = *(unsigned int *)($state + 0x528fc + $i * 4)
      set $unit = *(unsigned int *)($node + 8)
      if $unit == 59559 || $unit == 82037 || $unit == 2555
        printf "APPLE_OLD_POSITION4_SORT_TARGET_BEFORE index=%d unit=%u score=%g key=%g span=%d flags=%#x\n", $i, $unit, *(float *)$node, *(float *)($node + 4), *(short *)($node + 0x10), *(unsigned char *)($node + 0x12)
      end
      set $i = $i + 1
    end
    tbreak *0x1001d493
    commands
      silent
      set $i = 0
      printf "APPLE_OLD_POSITION4_SORTED_TOP30"
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

hbreak *0x408187
commands
  silent
  printf "APPLE_OLD_POSITION4_SORT_TRACE_READY\n"
  continue
end

continue
