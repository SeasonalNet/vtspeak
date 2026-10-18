set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-plain-crosswalk-old-span-filter-bypass-v3/legacy-span-filter-bypass-gdb.log
set logging overwrite on
set logging enabled on

set $calls = 0
set $raw0 = -1
set $raw1 = -1
set $raw2 = -1
set $raw3 = -1
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
    else
      if $position == 2
        set $raw2 = $raw_count
      else
        if $position == 3
          set $raw3 = $raw_count
        end
      end
    end
  end
  set $calls = $calls + 1
  if $calls <= 16
    printf "SPAN_FILTER_RAW position=%u call=%u candidates=%u\n", $position, $calls, $raw_count
  end
  continue
end

hbreak *0x1001c860
commands
  silent
  set $state = *(unsigned int *)($esp + 8)
  set $target = *(unsigned int *)($esp + 12)
  set $position = ($target - ($state + 0xb99e4)) / 6
  if $position >= 0 && $position < 4
    set $count_address = $state + 0xbb154 + $position * 0xfc + 0xf4
    set $pruned_count = *(short *)$count_address
    if $position == 0
      set $raw_count = $raw0
    else
      if $position == 1
        set $raw_count = $raw1
      else
        if $position == 2
          set $raw_count = $raw2
        else
          set $raw_count = $raw3
        end
      end
    end
    if $raw_count > $pruned_count && $raw_count > 0
      set $selected_ids = $state + 0x528fc
      printf "SPAN_FILTER_BYPASS position=%u raw=%u retained=%u candidate_units:", $position, $raw_count, $pruned_count
      set $i = 0
      while $i < $raw_count && $i < 100
        set $node = *(unsigned int *)($selected_ids + $i * 4)
        printf " %u/%u", *(unsigned int *)($node + 8), *(unsigned short *)($node + 0x10)
        set $i = $i + 1
      end
      printf "\n"
      set {short}($count_address) = $raw_count
    end
    set $scored_unit = *(unsigned int *)($esp + 4)
    if ($scored_unit == 272822) || ($scored_unit == 273369) || ($scored_unit == 52542) || ($scored_unit == 272823) || ($scored_unit == 273370) || ($scored_unit == 264072) || ($scored_unit == 272824) || ($scored_unit == 273371) || ($scored_unit == 264073) || ($scored_unit == 272825) || ($scored_unit == 255282) || ($scored_unit == 264074)
      set $score_return = *(unsigned int *)$esp
      set $score_position = $position
      set $score_weight = *(float *)($esp + 20)
      printf "SPAN_FILTER_SCORER_ENTER position=%u unit=%u count=%d weight=%g\n", $position, $scored_unit, *(short *)$count_address, $score_weight
      thbreak *$score_return
      commands
        silent
        printf "SPAN_FILTER_SCORER_RETURN position=%u unit=%u weight=%g score=%g\n", $score_position, $scored_unit, $score_weight, $st0
        continue
      end
    end
  end
  continue
end

hbreak *0x1001ed7f
commands
  silent
  set $selected = $selected + 1
  if $selected <= 48
    printf "SPAN_FILTER_SELECTED call=%u index=%d unit=%u\n", $selected, *(short *)$edi, $edx
  end
  continue
end

break *0x408187
commands
  silent
  printf "SPAN_FILTER_BYPASS_TRACE_READY\n"
  continue
end

continue
