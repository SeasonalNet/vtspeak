set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/kate-key5-score-path-hello/score-path-gdb.log
set logging overwrite on
set logging enabled on
set $selected_calls = 0

# Record the final per-position candidate arrays before local scoring.
break *0x100235f6
commands
  silent
  set $state = *(unsigned int *)($ebp + 0xc)
  set $position = *(int *)($ebp + 8)
  if $position < 7
    set $count = *(short *)($state + 0xae988 + $position * 0xfc)
    printf "KATE_KEY5_FINAL_CANDIDATES position=%d count=%d", $position, $count
    set $index = 0
    while $index < $count && $index < 80
      set $candidate_node = *(unsigned int *)($state + 0x477ac + $index * 4)
      printf " %u", *(unsigned int *)($candidate_node + 8)
      set $index = $index + 1
    end
    printf "\n"
  end
  continue
end

# FUN_100182e0 returns local cost in ST(0). Capture old-selected rows and
# current winners wherever they occur in the candidate chain.
break *0x1002389a
commands
  silent
  set $position = *(int *)($ebp + 8)
  set $node = *(unsigned int *)$esi
  set $unit = *(unsigned int *)($node + 8)
  if (($unit >= 272822) && ($unit <= 272825)) || (($unit == 280485) || ($unit == 21433) || ($unit >= 264072 && $unit <= 264074))
    printf "KATE_KEY5_LOCAL position=%d unit=%u rank=%g span=%d weighted=%d local=%g\n", $position, $unit, *(float *)($node + 4), *(short *)($node + 0xe), *(short *)($node + 0x10), $st0
  end
  continue
end

# Capture each candidate's best cumulative path cost and predecessor after
# adjacent-context scoring has populated the dynamic-programming row.
break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  set $record = $state + 0xae894 + $context * 0xfc
  tbreak *$return
  commands
    silent
    set $count = *(short *)($record + 0xf4)
    set $prev_record = $record - 0xfc
    printf "KATE_KEY5_PATH context=%d candidates=%d previous=%d", $context, $count, *(short *)($prev_record + 0xf4)
    set $index = 0
    while $index < $count && $index < 40
      set $unit = *(unsigned int *)($record + 0x7c + $index * 4)
      set $cost = *(float *)($state + 0x8212c + ($context * 30 + $index) * 4)
      set $pred = *(short *)($state + 0x9f664 + ($context * 30 + $index) * 2)
      set $pred_unit = *(unsigned int *)($prev_record + 0x7c + $pred * 4)
      if (($unit >= 272822) && ($unit <= 272825)) || (($unit == 280485) || ($unit == 21433) || ($unit >= 264072 && $unit <= 264074))
        printf " {unit=%u cost=%g pred-index=%d pred=%u}", $unit, $cost, $pred, $pred_unit
      end
      set $index = $index + 1
    end
    printf "\n"
    continue
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 32
    printf "KATE_KEY5_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "KATE_KEY5_SCORE_PATH_TRACE_READY\n"
  continue
end

continue
