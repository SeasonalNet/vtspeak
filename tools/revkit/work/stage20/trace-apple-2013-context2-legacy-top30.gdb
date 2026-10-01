set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-2013-context2-legacy-top30-2026-09-29/adapted-context2-legacy-top30-gdb.log
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

# At context 2 only, retain the unique IDs in the native 2006 top-30 list.
hbreak *0x10023814
commands
  silent
  if $context == 2
    set $count_address = $state + 0xae988 + $context * 0xfc
    set $count = *(unsigned short *)$count_address
    set $array = $state + 0x477ac
    set $read_index = 0
    set $write_index = 0
    printf "APPLE_2013_CONTEXT2_BEFORE count=%u ids:", $count
    while $read_index < $count && $read_index < 10000
      set $node = *(unsigned int *)($array + $read_index * 4)
      set $unit = *(unsigned int *)($node + 8)
      printf " %u", $unit
      set $keep = 0
      if ($unit == 189455) || ($unit == 100180) || ($unit == 2645) || ($unit == 188302) || ($unit == 199880) || ($unit == 38294) || ($unit == 137215) || ($unit == 45509) || ($unit == 61725) || ($unit == 40810) || ($unit == 151210) || ($unit == 41008) || ($unit == 96145) || ($unit == 231014) || ($unit == 118661) || ($unit == 56167) || ($unit == 108878) || ($unit == 244717) || ($unit == 213256) || ($unit == 64719) || ($unit == 64728) || ($unit == 82377) || ($unit == 82709) || ($unit == 97601) || ($unit == 128863) || ($unit == 40503) || ($unit == 86012) || ($unit == 78432) || ($unit == 150075) || ($unit == 32957)
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
    printf "\nAPPLE_2013_CONTEXT2_AFTER count=%u ids:", $write_index
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
    printf "APPLE_2013_CONTEXT2_LEGACY_TOP30_SELECTED n=%u unit=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT2_LEGACY_TOP30_TRACE_READY\n"
  continue
end

continue
