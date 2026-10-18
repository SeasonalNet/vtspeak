set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-2013-context2-prune-new-only-2026-09-29/adapted-context2-prune-new-only-gdb.log
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

# Remove only the seven IDs outside the 2006 context-2 set. Keep the shared
# candidate records and all coverage, scoring, and transition logic intact.
hbreak *0x10023814
commands
  silent
  if $context == 2
    set $count_address = $state + 0xae988 + $context * 0xfc
    set $count = *(unsigned short *)$count_address
    set $array = $state + 0x477ac
    set $read_index = 0
    set $write_index = 0
    printf "APPLE_2013_CONTEXT2_PRUNE_BEFORE count=%u ids:", $count
    while $read_index < $count && $read_index < 10000
      set $node = *(unsigned int *)($array + $read_index * 4)
      set $unit = *(unsigned int *)($node + 8)
      printf " %u", $unit
      if ($unit != 42412) && ($unit != 92769) && ($unit != 127405) && ($unit != 145945) && ($unit != 258032) && ($unit != 265201) && ($unit != 271727)
        set *(unsigned int *)($array + $write_index * 4) = $node
        set $write_index = $write_index + 1
      end
      set $read_index = $read_index + 1
    end
    printf "\nAPPLE_2013_CONTEXT2_PRUNE_AFTER count=%u ids:", $write_index
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
    printf "APPLE_2013_PRUNED_SELECTED n=%u unit=%u\n", $selected_count, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_CONTEXT2_PRUNE_TRACE_READY\n"
  continue
end

continue
