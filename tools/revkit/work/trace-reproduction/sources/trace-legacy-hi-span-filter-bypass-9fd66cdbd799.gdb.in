set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hi-2006-span-filter-bypass/legacy-span-filter-bypass-gdb.log
set logging overwrite on
set logging enabled on

set $raw0 = -1
set $raw1 = -1
set $selected = 0

hbreak *0x1001cd60
commands
  silent
  set $position = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $raw_count = *(short *)($state + 0xbb248 + $position * 0xfc)
  if $position == 0
    set $raw0 = $raw_count
  else
    if $position == 1
      set $raw1 = $raw_count
    end
  end
  printf "HI_SPAN_FILTER_RAW position=%u candidates=%u\n", $position, $raw_count
  continue
end

hbreak *0x1001c860
commands
  silent
  set $state = *(unsigned int *)($esp + 8)
  set $target = *(unsigned int *)($esp + 12)
  set $position = ($target - ($state + 0xb99e4)) / 6
  if $position >= 0 && $position < 2
    set $count_address = $state + 0xbb154 + $position * 0xfc + 0xf4
    set $pruned_count = *(short *)$count_address
    if $position == 0
      set $raw_count = $raw0
    else
      set $raw_count = $raw1
    end
    if $raw_count > $pruned_count && $raw_count > 0
      set $selected_ids = $state + 0x528fc
      printf "HI_SPAN_FILTER_BYPASS position=%u raw=%u retained=%u candidates:", $position, $raw_count, $pruned_count
      set $i = 0
      while $i < $raw_count && $i < 100
        set $node = *(unsigned int *)($selected_ids + $i * 4)
        printf " %u/%u", *(unsigned int *)($node + 8), *(unsigned short *)($node + 0x10)
        set $i = $i + 1
      end
      printf "\n"
      set {short}$count_address = $raw_count
    end
    set $unit = *(unsigned int *)($esp + 4)
    set $return = *(unsigned int *)$esp
    set $weight = *(float *)($esp + 20)
    printf "HI_SPAN_FILTER_SCORE_ENTRY position=%u unit=%u candidates=%d weight=%g\n", $position, $unit, *(short *)$count_address, $weight
    set $score_position = $position
    set $score_unit = $unit
    set $score_weight = $weight
    thbreak *$return
    commands
      silent
      printf "HI_SPAN_FILTER_SCORE_RETURN position=%u unit=%u weight=%g cost=%g\n", $score_position, $score_unit, $score_weight, $st0
      continue
    end
  end
  continue
end

hbreak *0x1001ed7f
commands
  silent
  set $selected = $selected + 1
  if $selected <= 24
    printf "HI_SPAN_FILTER_SELECTED call=%u index=%d unit=%u\n", $selected, *(short *)$edi, $edx
  end
  continue
end

break *0x408187
commands
  silent
  printf "HI_SPAN_FILTER_TRACE_READY\n"
  continue
end

continue
