set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/stage7-complete/trace.log
set logging overwrite on
set logging enabled on
set $rank_calls = 0
set $prune_calls = 0
set $pair_calls = 0
set $handoff_calls = 0

break *0x10023af0
commands
  silent
  set $rank_calls = $rank_calls + 1
  set $count = *(unsigned short *)(*(unsigned int *)($esp + 20) + 0xec620)
  set $limit = *(int *)($esp + 4)
  set $budget = *(int *)($esp + 8)
  set $state = *(unsigned int *)($esp + 20)
  set $model = *(unsigned int *)($esp + 24)
  set $ids = $state + 0xec170
  set $sizes = *(unsigned int *)($model + 0x8c)
  set $output = *(unsigned int *)($esp + 12)
  set $population = 0
  set $i = 0
  while $i < $count
    set $id = *(unsigned int *)($ids + $i * 4)
    set $population = $population + *(unsigned short *)($sizes + $id * 2)
    set $i = $i + 1
  end
  set $ret = *(unsigned int *)$esp
  printf "CLASS_RANK_ENTRY call=%d input_classes=%u population=%d max_classes=%d population_budget=%d\n", $rank_calls, $count, $population, $limit, $budget
  tbreak *$ret
  commands
    silent
    printf "CLASS_RANK_RETURN call=%d selected_classes=%u output_ids=", $rank_calls, (unsigned short)$eax
    set $i = 0
    while $i < (unsigned short)$eax && $i < 10
      printf "%u ", *(unsigned int *)($output + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    continue
  end
  continue
end

break *0x10024510
commands
  silent
  set $prune_calls = $prune_calls + 1
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $multiplier = *(int *)($esp + 12)
  set $record = $state + 0xae894 + $context * 0xfc
  set $count = *(unsigned short *)($record + 0xf4)
  set $scores = $state + 0x8212c + $context * 0x78
  set $sum = 0.0
  set $minimum = *(float *)$scores
  set $flagged = 0
  set $flagged_above = 0
  set $i = 0
  while $i < $count
    set $score = *(float *)($scores + $i * 4)
    set $sum = $sum + $score
    if $score < $minimum
      set $minimum = $score
    end
    set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
    if *(short *)($node + 0xc) != 0
      set $flagged = $flagged + 1
      if $score > $cutoff
        set $flagged_above = $flagged_above + 1
      end
    end
    set $i = $i + 1
  end
  if $count > 0
    set $average = $sum / $count
    set $cutoff = $average + $multiplier * ($average - $minimum) / 4.0
  else
    set $average = 0.0
    set $cutoff = 0.0
  end
  set $token_row = *(unsigned char *)($state + 0xec628 + $context * 6)
  set $gate = *(signed char *)(*(unsigned int *)($state + 0x4c) + $token_row * 0x3c0 + 0x92b)
  set $ret = *(unsigned int *)$esp
  printf "PRUNE_ENTRY context=%d count=%u multiplier=%d average=%g minimum=%g cutoff=%g flagged=%d flagged_above_cutoff=%d gate_value=%d\n", $context, $count, $multiplier, $average, $minimum, $cutoff, $flagged, $flagged_above, $gate
  tbreak *$ret
  commands
    silent
    printf "PRUNE_RETURN context=%d kept=%u first_ids=", $context, *(unsigned short *)($record + 0xf4)
    set $i = 0
    while $i < *(unsigned short *)($record + 0xf4) && $i < 5
      printf "%u ", *(unsigned int *)($record + 0x7c + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    continue
  end
  continue
end

break *0x10018c80
commands
  silent
  set $pair_calls = $pair_calls + 1
  if $pair_calls <= 2
    set $context = *(int *)($esp + 4)
    set $state = *(unsigned int *)($esp + 8)
    set $record = $state + 0xae894 + $context * 0xfc
    set $previous = $state + 0xae894 + ($context - 1) * 0xfc
    printf "TRANSITION context=%d current_count=%u previous_count=%u mode=%u\n", $context, *(unsigned short *)($record + 0xf4), *(unsigned short *)($previous + 0xf4), *(unsigned char *)($state + 0xec628 + $context * 6 + 4)
  end
  continue
end

break *0x10024900
commands
  silent
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  printf "BACKTRACK_SPAN first=%d last=%d contexts=%d\n", $first, $last, $last - $first + 1
  continue
end

break *0x1001b200
commands
  silent
  set $handoff_calls = $handoff_calls + 1
  continue
end

continue
