set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/continuity-gate-flags/2013-continuity-gate-flags-override-gdb.log
set logging overwrite on
set logging enabled on
set $postselection_calls = 0
set $current_context = 0
set $current_state = 0
set $current_model = 0
set $last_sum = -1
set $unit_score_calls = 0
set $selected_calls = 0

# Reproduce the controlled old-cutoff intervention used for the matched
# four-position comparison: accept a nonempty metric sum above two.
break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  set $last_sum = -1
  tbreak *$return
  commands
    silent
    set $native = $al
    if $native == 0 && $last_sum > 2
      printf "OLD_CUTOFF_OVERRIDE slot=%u sum=%u native=%u overridden=1\n", $slot, $last_sum, $native
      set $eax = 1
    end
    continue
  end
  continue
end

break *0x10023060
commands
  silent
  set $caller = *(unsigned int *)$esp
  if $caller >= 0x10024060 && $caller < 0x100242a0
    set $ids = *(unsigned int *)($esp + 4)
    set $count = *(unsigned short *)($esp + 8)
    set $context = *(unsigned int *)($esp + 12)
    set $weights = *(unsigned int *)($context + 0x8c)
    set $return = $caller
    set $i = 0
    set $last_sum = 0
    while $i < $count && $i < 64
      set $id = *(unsigned int *)($ids + $i * 4)
      set $last_sum = $last_sum + *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      continue
    end
  end
  continue
end

# Save arguments at the function entry. The shared branch address below is
# reached after the compiler has adjusted ESP, so it cannot read entry args.
break *0x10023350
commands
  silent
  set $current_context = *(int *)($esp + 4)
  set $current_state = *(unsigned int *)($esp + 8)
  set $current_model = *(unsigned int *)($esp + 12)
  set $unit_signature = *(unsigned int *)($current_model + 0x64)
  set *(unsigned char *)($unit_signature + 272822 * 7 + 6) = 0x80
  set *(unsigned char *)($unit_signature + 272823 * 7 + 6) = 0x80
  set *(unsigned char *)($unit_signature + 272824 * 7 + 6) = 0x80
  printf "CONTROLLED_FLAG_OVERRIDE context=%u unit_flags=0x80 for 272822-272824\n", $current_context
  continue
end

break *0x100182e0
commands
  silent
  set $unit_score_calls = $unit_score_calls + 1
  if $unit_score_calls <= 128
    printf "OVERRIDE_LOCAL_SCORE call=%u context=%u unit=%u\n", $unit_score_calls, $current_context, *(unsigned int *)($esp + 4)
  end
  continue
end

break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 24
    printf "OVERRIDE_SELECTED_UNIT call=%u unit=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

# FUN_10023350 calls FUN_100230a0 to attach left/right/total span fields.
# This join point is after its full-span filtering branch and before local
# scoring, so the node list records the candidates that the scorer receives.
break *0x10023814
commands
  silent
  set $postselection_calls = $postselection_calls + 1
  set $ctx = $current_context
  set $state = $current_state
  set $record = $state + 0xae894 + $ctx * 0xfc
  set $count = *(unsigned short *)($record + 0xf4)
  set $total = *(unsigned int *)($state + 0xec624)
  printf "POST_COVERAGE_GATE call=%u context=%u total_positions=%u candidate_count=%u input_ids=", $postselection_calls, $ctx, $total, $count
  set $i = 0
  set $input_count = *(unsigned short *)$record
  while $i < $input_count && $i < 12
    printf "%u ", *(unsigned int *)($record + 4 + $i * 4)
    set $i = $i + 1
  end
  printf "\n"
  set $successors = *(unsigned int *)($current_model + 0x30)
  printf "MODEL_TABLE_30_LOOKUP context=%u 272822->%u 272823->%u 272824->%u 272825->%u\n", $ctx, *(unsigned int *)($successors + 272822 * 4), *(unsigned int *)($successors + 272823 * 4), *(unsigned int *)($successors + 272824 * 4), *(unsigned int *)($successors + 272825 * 4)
  set $unit_signature = *(unsigned int *)($current_model + 0x64)
  printf "MODEL_SIGNATURE_BYTE6 context=%u 272822=%#x 272823=%#x 272824=%#x 272825=%#x\n", $ctx, *(unsigned char *)($unit_signature + 272822 * 7 + 6), *(unsigned char *)($unit_signature + 272823 * 7 + 6), *(unsigned char *)($unit_signature + 272824 * 7 + 6), *(unsigned char *)($unit_signature + 272825 * 7 + 6)
  set $i = 0
  while $i < $count && $i < 64
    set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
    printf "POST_COVERAGE_CANDIDATE context=%u index=%u unit=%u left=%d right=%d span=%d\n", $ctx, $i, *(unsigned int *)($node + 8), *(short *)($node + 0xc), *(short *)($node + 0x10), *(short *)($node + 0xe)
    set $i = $i + 1
  end
  continue
end

continue
