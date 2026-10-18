set pagination off
set confirm off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/continuity-bit7-unit-major-overlay/2013-continuity-bit7-unit-major-gdb.log
set logging overwrite on
set logging enabled on
set $postselection_calls = 0
set $current_context = 0
set $current_state = 0
set $current_model = 0
set $last_sum = -1
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

break *0x10023350
commands
  silent
  set $current_context = *(int *)($esp + 4)
  set $current_state = *(unsigned int *)($esp + 8)
  set $current_model = *(unsigned int *)($esp + 12)
  continue
end

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
  set $unit_signature = *(unsigned int *)($current_model + 0x64)
  printf "MODEL_SIGNATURE_BYTE6 context=%u 272822=%#x 272823=%#x 272824=%#x 272825=%#x\n", $ctx, *(unsigned char *)($unit_signature + 272822 * 7 + 6), *(unsigned char *)($unit_signature + 272823 * 7 + 6), *(unsigned char *)($unit_signature + 272824 * 7 + 6), *(unsigned char *)($unit_signature + 272825 * 7 + 6)
  set $i = 0
  set $span1 = 0
  set $span2 = 0
  set $span3 = 0
  set $span4 = 0
  set $target_mask = 0
  set $target_class_1 = 0
  set $target_class_2 = 0
  set $target_class_3 = 0
  set $target_class_4 = 0
  set $class_table = *(unsigned int *)($current_model + 0x30)
  while $i < $count && $i < 512
    set $node = *(unsigned int *)($state + 0x477ac + $i * 4)
    set $unit = *(unsigned int *)($node + 8)
    set $span = *(short *)($node + 0xe)
    if $span == 1
      set $span1 = $span1 + 1
    end
    if $span == 2
      set $span2 = $span2 + 1
    end
    if $span == 3
      set $span3 = $span3 + 1
    end
    if $span == 4
      set $span4 = $span4 + 1
    end
    if $unit == 272822
      set $target_mask = $target_mask | 1
    end
    if $unit == 272823
      set $target_mask = $target_mask | 2
    end
    if $unit == 272824
      set $target_mask = $target_mask | 4
    end
    if $unit == 272825
      set $target_mask = $target_mask | 8
    end
    if $unit >= 272822 && $unit <= 272825
      printf "TARGET_CANDIDATE context=%u unit=%u left=%d right=%d span=%d\n", $ctx, $unit, *(short *)($node + 0xc), *(short *)($node + 0x10), $span
    end
    set $class_id = *(unsigned int *)($class_table + $unit * 4)
    if $class_id == 32928
      set $target_class_1 = $target_class_1 + 1
      printf "TARGET_CLASS_MEMBER context=%u unit=%u class=%u span=%d\n", $ctx, $unit, $class_id, $span
    end
    if $class_id == 13478
      set $target_class_2 = $target_class_2 + 1
      printf "TARGET_CLASS_MEMBER context=%u unit=%u class=%u span=%d\n", $ctx, $unit, $class_id, $span
    end
    if $class_id == 9147
      set $target_class_3 = $target_class_3 + 1
      printf "TARGET_CLASS_MEMBER context=%u unit=%u class=%u span=%d\n", $ctx, $unit, $class_id, $span
    end
    if $class_id == 18857
      set $target_class_4 = $target_class_4 + 1
      printf "TARGET_CLASS_MEMBER context=%u unit=%u class=%u span=%d\n", $ctx, $unit, $class_id, $span
    end
    set $i = $i + 1
  end
  printf "POST_SPAN_DISTRIBUTION context=%u scanned=%u span1=%u span2=%u span3=%u span4=%u target_mask=%#x target_classes=%u,%u,%u,%u\n", $ctx, $i, $span1, $span2, $span3, $span4, $target_mask, $target_class_1, $target_class_2, $target_class_3, $target_class_4
  continue
end

continue
