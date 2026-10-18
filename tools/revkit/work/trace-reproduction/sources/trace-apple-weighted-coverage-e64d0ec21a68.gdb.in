set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  set $count = *(short *)($state + 0xae988 + $position * 0xfc)
  set $model = *(unsigned int *)($state + 0x4c)
  set $index = 0
  set $weighted_count = 0
  set $weighted_max = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $weighted = *(short *)($node + 0x10)
    if $weighted > 0
      set $weighted_count = $weighted_count + 1
      if $weighted > $weighted_max
        set $weighted_max = $weighted
      end
    end
    set $index = $index + 1
  end
  printf "APPLE_WEIGHTED_COVERAGE position=%d count=%d nonzero=%d max=%d\n", $position, $count, $weighted_count, $weighted_max
  set $index = 0
  while $index < $count
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $weighted = *(short *)($node + 0x10)
    if $weighted > 0
      set $unit = *(unsigned int *)($node + 8)
      set $span = *(short *)($node + 0xe)
      set $right = *(short *)($node + 0xc)
      set $rank_term = -($span + 100 * $weighted)
      set $signature = *(unsigned char *)($model + 0x64 + $unit * 7 + 6)
      printf "APPLE_WEIGHTED_ROW position=%d index=%d id=%u rank_term=%d fields_c_e_10=%d,%d,%d sig6=%02x\n", $position, $index, $unit, $rank_term, $right, $span, $weighted, $signature
    end
    set $index = $index + 1
  end
  set $index = 0
  while $index < $count && $index < 12
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $unit = *(unsigned int *)($node + 8)
    set $span = *(short *)($node + 0xe)
    set $weighted = *(short *)($node + 0x10)
    set $rank_term = -($span + 100 * $weighted)
    printf "APPLE_WEIGHTED_TOP position=%d index=%d id=%u rank_term=%d fields_c_e_10=%d,%d,%d\n", $position, $index, $unit, $rank_term, *(short *)($node + 0xc), $span, $weighted
    set $index = $index + 1
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_WEIGHTED_COVERAGE_TRACE_READY\n"
  continue
end

continue
