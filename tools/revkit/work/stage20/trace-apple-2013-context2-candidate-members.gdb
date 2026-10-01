set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-2013-context2-candidate-members-2026-09-29/adapted-context2-candidate-members-gdb.log
set logging overwrite on
set logging enabled on
set $context = -1
set $state = 0

hbreak *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  continue
end

# Dump the complete post-coverage set before scoring or the shortlist cap.
hbreak *0x10023814
commands
  silent
  if $context == 2
    set $count = *(unsigned short *)($state + 0xae988 + $context * 0xfc)
    set $index = 0
    printf "APPLE_2013_CONTEXT2_ALL_CANDIDATES count=%u ids:", $count
    while $index < $count && $index < 10000
      set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
      printf " %u", *(unsigned int *)($node + 8)
      set $index = $index + 1
    end
    printf "\n"
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
