set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/stage7-complete/trace-forced-rank.log
set logging overwrite on
set logging enabled on
set $rank_calls = 0
set $score_calls = 0

break *0x10023a70
commands
  silent
  set $score_calls = $score_calls + 1
  set $target = *(unsigned int *)($esp + 4)
  set $candidate = *(unsigned int *)($esp + 8)
  set $mode = *(short *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "CLASS_MISMATCH_ENTRY call=%d mode=%d target=", $score_calls, $mode
  x/5bx $target
  printf " candidate="
  x/5bx $candidate
  tbreak *$ret
  commands
    silent
    printf "CLASS_MISMATCH_RETURN call=%d score=%d\n", $score_calls, $eax
    continue
  end
  continue
end

break *0x10023af0
commands
  silent
  set $rank_calls = $rank_calls + 1
  set $limit = *(int *)($esp + 4)
  set $budget = *(int *)($esp + 8)
  set $output = *(unsigned int *)($esp + 12)
  set $state = *(unsigned int *)($esp + 20)
  set $model = *(unsigned int *)($esp + 24)
  set $count = *(unsigned short *)($state + 0xec620)
  set $ids = $state + 0xec170
  set $keys = *(unsigned int *)($model + 0x90)
  set $sizes = *(unsigned int *)($model + 0x8c)
  printf "CLASS_RANK_ENTRY call=%d classes=%u max=%d budget=%d\n", $rank_calls, $count, $limit, $budget
  set $i = 0
  while $i < $count && $i < 8
    set $id = *(unsigned int *)($ids + $i * 4)
    printf "CLASS_INPUT index=%d id=%u population=%u key=", $i, $id, *(unsigned short *)($sizes + $id * 2)
    x/5bx ($keys + $id * 5)
    set $i = $i + 1
  end
  if $rank_calls == 2 && $count > 1
    set *(int *)($esp + 8) = 1
    printf "CONTROLLED_OVERRIDE call=%d population_budget=1\n", $rank_calls
  end
  set $ret = *(unsigned int *)$esp
  tbreak *$ret
  commands
    silent
    printf "CLASS_RANK_RETURN call=%d selected=%u output_ids=", $rank_calls, (unsigned short)$eax
    set $i = 0
    while $i < (unsigned short)$eax && $i < 8
      printf "%u ", *(unsigned int *)($output + $i * 4)
      set $i = $i + 1
    end
    printf "\n"
    continue
  end
  continue
end

continue
