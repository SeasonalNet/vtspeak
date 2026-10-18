set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-2013-context4-legacy-top30-2026-09-29/adapted-context4-legacy-top30-gdb.log
set logging overwrite on
set logging enabled on
set $context = -1
set $state = 0
set $selected_count = 0

hbreak *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  continue
end

# At context 4 only, retain the unique IDs in the native 2006 top-30 list.
# Deduplicate repeated 2013 entries before enforcing the 30-row bound.
hbreak *0x10023814
commands
  silent
  if $context == 4
    set $count_address = $state + 0xae988 + $context * 0xfc
    set $count = *(unsigned short *)$count_address
    set $array = $state + 0x477ac
    set $read_index = 0
    set $write_index = 0
    printf "APPLE_2013_CONTEXT4_BEFORE count=%u ids:", $count
    while $read_index < $count && $read_index < 10000
      set $node = *(unsigned int *)($array + $read_index * 4)
      set $unit = *(unsigned int *)($node + 8)
      printf " %u", $unit
      set $keep = 0
      if ($unit == 158831) || ($unit == 232463) || ($unit == 86328) || ($unit == 97885) || ($unit == 103239) || ($unit == 59559) || ($unit == 218671) || ($unit == 256582) || ($unit == 129560) || ($unit == 41058) || ($unit == 40215) || ($unit == 131119) || ($unit == 258823) || ($unit == 28632) || ($unit == 27967) || ($unit == 27795) || ($unit == 82037) || ($unit == 27441) || ($unit == 11136) || ($unit == 2555) || ($unit == 261692) || ($unit == 258101) || ($unit == 250167) || ($unit == 9835) || ($unit == 11069) || ($unit == 238346) || ($unit == 160420) || ($unit == 249481) || ($unit == 162687) || ($unit == 102356)
        set $seen = 0
        set $prior = 0
        while $prior < $write_index
          set $prior_node = *(unsigned int *)($array + $prior * 4)
          if *(unsigned int *)($prior_node + 8) == $unit
            set $seen = 1
          end
          set $prior = $prior + 1
        end
        if !$seen
          set $keep = 1
        end
      end
      if $keep
        set *(unsigned int *)($array + $write_index * 4) = $node
        set $write_index = $write_index + 1
      end
      set $read_index = $read_index + 1
    end
    printf "\nAPPLE_2013_CONTEXT4_AFTER count=%u ids:", $write_index
    set $read_index = 0
    while $read_index < $write_index
      set $node = *(unsigned int *)($array + $read_index * 4)
      printf " %u", *(unsigned int *)($node + 8)
      set $read_index = $read_index + 1
    end
    printf "\n"
    set *(unsigned short *)$count_address = $write_index
    set $ebx = $write_index
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_count = $selected_count + 1
  if $selected_count <= 64
    printf "APPLE_2013_CONTEXT4_LEGACY_TOP30_SELECTED n=%u unit=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT4_LEGACY_TOP30_TRACE_READY\n"
  continue
end

continue
