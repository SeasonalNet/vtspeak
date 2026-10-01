set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/apple-selector-byte0-context3-pool-2026-09-29/adapted-context3-pool-gdb.log
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

# Observe the complete unique-row pool before transition scoring and the cap.
hbreak *0x10023814
commands
  silent
  if $context == 3
    set $count_address = $state + 0xae988 + $context * 0xfc
    set $count = *(unsigned short *)$count_address
    set $array = $state + 0x477ac
    set $read_index = 0
    printf "APPLE_2013_CONTEXT3_POOL count=%u ids:", $count
    while $read_index < $count && $read_index < 10000
      set $node = *(unsigned int *)($array + $read_index * 4)
      printf " %u", *(unsigned int *)($node + 8)
      set $read_index = $read_index + 1
    end
    printf "\n"
  end
  continue
end

continue
