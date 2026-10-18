set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /probe/candidates-gdb.log
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

continue
