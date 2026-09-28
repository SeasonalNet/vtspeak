set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/legacy-continuity-candidates-hardware-gdb.log
set logging overwrite on
set logging enabled on
set $legacy_continuity_calls = 0

hbreak *0x1001cd60
commands
  silent
  set $legacy_continuity_calls = $legacy_continuity_calls + 1
  set $position = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  thbreak *$return
  commands
    silent
    if $legacy_continuity_calls <= 4
      set $context = $state + 0xbb154 + $position * 0xfc
      set $count = *(unsigned short *)($context + 0xf4)
      set $candidate_array = $state + 0x528fc
      set $total_positions = *(unsigned int *)($state + 0xf8e68)
      printf "LEGACY_CONTINUITY_RETURN position=%u candidates=%u total_positions=%u\n", $position, $count, $total_positions
      set $i = 0
      while $i < $count && $i < 100
        set $candidate = *(unsigned int *)($candidate_array + $i * 4)
        printf "LEGACY_CONTINUITY_CANDIDATE position=%u index=%u unit=%u coverage=%u left=%u right=%u preliminary_cost=%g\n", $position, $i, *(unsigned int *)($candidate + 8), *(unsigned short *)($candidate + 0x10), *(unsigned short *)($candidate + 0xc), *(unsigned short *)($candidate + 0xe), *(float *)($candidate + 4)
        set $i = $i + 1
      end
    end
    continue
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_CONTINUITY_TRACE_READY\n"
  continue
end

continue
