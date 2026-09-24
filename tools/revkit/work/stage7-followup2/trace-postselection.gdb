set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/stage7-followup2/postselection-gdb.log
set logging overwrite on
set logging enabled on
set $rank_calls = 0
set $final_calls = 0
set $prosody_calls = 0

break *0x10023af0
commands
  silent
  set $rank_calls = $rank_calls + 1
  if $rank_calls <= 3
    set $limit = *(int *)($esp + 4)
    set $budget = *(int *)($esp + 8)
    set $output = *(unsigned int *)($esp + 12)
    set $signature = *(unsigned int *)($esp + 16)
    set $state = *(unsigned int *)($esp + 20)
    set $model = *(unsigned int *)($esp + 24)
    set $count = *(unsigned short *)($state + 0xec620)
    set $ids = $state + 0xec170
    set $keys = *(unsigned int *)($model + 0x90)
    set $durations = *(unsigned int *)($model + 0x8c)
    set $ret = *(unsigned int *)$esp
    printf "KEY_RANK_ENTRY call=%d limit=%d budget=%d count=%u signature=", $rank_calls, $limit, $budget, $count
    x/7bx $signature
    set $i = 0
    while $i < $count && $i < 8
      set $id = *(unsigned int *)($ids + $i * 4)
      printf "KEY_CANDIDATE call=%d index=%d class_id=%u duration=%u key=", $rank_calls, $i, $id, *(unsigned short *)($durations + $id * 2)
      x/5bx ($keys + $id * 5)
      set $i = $i + 1
    end
    tbreak *$ret
    commands
      silent
      set $selected = (unsigned short)$eax
      printf "KEY_RANK_RETURN call=%d selected=%u class_ids=", $rank_calls, $selected
      set $i = 0
      while $i < $selected && $i < 8
        printf "%u ", *(unsigned int *)($output + $i * 4)
        set $i = $i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

break *0x10023350
commands
  silent
  set $final_calls = $final_calls + 1
  if $final_calls <= 3
    set $context = *(int *)($esp + 4)
    set $state = *(unsigned int *)($esp + 8)
    set $record = $state + 0xae894 + $context * 0xfc
    set $input_count = *(unsigned short *)$record
    set $ret = *(unsigned int *)$esp
    printf "FINAL_RANK_ENTRY call=%d context=%d input_count=%u input_ids=", $final_calls, $context, $input_count
    set $i = 0
    while $i < $input_count && $i < 8
      printf "%u ", *(unsigned int *)($record + 4 + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$ret
    commands
      silent
      set $selected = *(unsigned short *)($record + 0xf4)
      printf "FINAL_RANK_RETURN call=%d selected=%u unit_ids=", $final_calls, $selected
      set $i = 0
      while $i < $selected && $i < 10
        printf "%u ", *(unsigned int *)($record + 0x7c + $i * 4)
        set $i = $i + 1
      end
      printf "\nFINAL_RANK_COSTS call=%d\n", $final_calls
      set $i = 0
      while $i < $selected && $i < 10
        set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
        printf "FINAL_COST position=%d unit=%u score=%g duration_penalty=%d context_penalty=%d\n", $i, *(unsigned int *)($node + 8), *(float *)$node, *(short *)($node + 0xe), *(short *)($node + 0x10)
        set $i = $i + 1
      end
      continue
    end
  end
  continue
end

break *0x100182e0
commands
  silent
  set $prosody_calls = $prosody_calls + 1
  if $prosody_calls <= 8
    set $unit = *(int *)($esp + 4)
    set $target = *(unsigned int *)($esp + 12)
    set $feature_view = *(unsigned int *)($esp + 16)
    set $scale = *(float *)($esp + 20)
    set $ret = *(unsigned int *)$esp
    printf "CEP_PROSODY_ENTRY call=%d unit=%d scale=%g target=", $prosody_calls, $unit, $scale
    x/7bx $target
    printf "CEP_PROSODY_FEATURE_VIEW="
    x/10bx $feature_view
    tbreak *$ret
    commands
      silent
      printf "CEP_PROSODY_RETURN call=%d cost=%g\n", $prosody_calls, $st0
      continue
    end
  end
  continue
end

set $pair_calls = 0
break *0x10018c80
commands
  silent
  set $pair_calls = $pair_calls + 1
  if $pair_calls <= 4
    set $context = *(int *)($esp + 4)
    set $state = *(unsigned int *)($esp + 8)
    set $model = *(unsigned int *)($esp + 12)
    set $record = $state + 0xae894 + $context * 0xfc
    set $previous = $state + 0xae894 + ($context - 1) * 0xfc
    set $count = *(unsigned short *)($record + 0xf4)
    set $previous_count = *(unsigned short *)($previous + 0xf4)
    set $ret = *(unsigned int *)$esp
    printf "PAIR_SCORE_ENTRY call=%d context=%d current_count=%u previous_count=%u raw_count=%u derived_count=%u raw_rows=%p derived_rows=%p\n", $pair_calls, $context, $count, $previous_count, *(unsigned short *)($model + 8), *(unsigned short *)($model + 0x14), *(unsigned int *)($model + 0xc), *(unsigned int *)($model + 0x18)
    set $i = 0
    while $i < $count && $i < 5
      printf "PAIR_CURRENT call=%d index=%d unit=%u\n", $pair_calls, $i, *(unsigned int *)($record + 0x7c + $i * 4)
      set $i = $i + 1
    end
    set $i = 0
    while $i < $previous_count && $i < 5
      printf "PAIR_PREVIOUS call=%d index=%d unit=%u\n", $pair_calls, $i, *(unsigned int *)($previous + 0x7c + $i * 4)
      set $i = $i + 1
    end
    set $phone_mode = *(unsigned char *)($state + 0xec628 + $context * 6 + 4)
    set $current_mode = 0
    set $previous_mode = 2
    if $phone_mode == 2
      set $current_mode = 1
      set $previous_mode = 1
    end
    set $current_unit = *(unsigned int *)($record + 0x7c)
    set $previous_unit = *(unsigned int *)($previous + 0x7c)
    set $current_metric_base = *(unsigned int *)($model + 0x3c + $current_mode * 4)
    set $previous_metric_base = *(unsigned int *)($model + 0x3c + $previous_mode * 4)
    set $current_metric = *(unsigned short *)($current_metric_base + $current_unit * 2) & 0x3fff
    set $previous_metric = *(unsigned short *)($previous_metric_base + $previous_unit * 2) & 0x3fff
    set $row = $current_metric
    set $column = $previous_metric
    if $column > $row
      set $row = $previous_metric
      set $column = $current_metric
    end
    set $raw_rows = *(unsigned int *)($model + 0xc)
    set $raw_distance = *(float *)(*(unsigned int *)($raw_rows + $row * 4) + $column * 4)
    printf "CEP_TABLE_PAIR call=%d phone_mode=%u current_mode=%d previous_mode=%d current_metric=%u previous_metric=%u raw_distance=%g\n", $pair_calls, $phone_mode, $current_mode, $previous_mode, $current_metric, $previous_metric, $raw_distance
    tbreak *$ret
    commands
      silent
      set $count = *(unsigned short *)($record + 0xf4)
      printf "PAIR_SCORE_RETURN call=%d current_count=%u\n", $pair_calls, $count
      set $i = 0
      while $i < $count && $i < 5
        set $slot = $context * 30 + $i
        printf "PAIR_RESULT call=%d index=%d unit=%u previous_index=%d cumulative_cost=%g\n", $pair_calls, $i, *(unsigned int *)($record + 0x7c + $i * 4), *(short *)($state + 0x9f664 + $slot * 2), *(float *)($state + 0x8212c + $slot * 4)
        set $i = $i + 1
      end
      continue
    end
  end
  continue
end

continue
