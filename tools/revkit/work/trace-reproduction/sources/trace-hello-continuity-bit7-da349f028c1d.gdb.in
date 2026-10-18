set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $postselection_calls = 0
set $current_context = 0
set $current_state = 0
set $current_model = 0
set $last_sum = -1

# Accept nonempty 2013 shortlists with sums above the legacy >2 cutoff.
break *0x10024060
commands
  silent
  set $slot = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $target_row = $state + ($slot * 3 + 0x76314) * 2
  set $bank = *(unsigned char *)$target_row
  set $leaf = *(unsigned char *)($target_row + 2)
  set $voice = *(unsigned int *)($state + 0x4c)
  set $target_signature = $voice + $bank * 0x3c0 + 0x6e2 + $leaf * 7
  set $key0 = *(unsigned char *)(0x1007b7ec + *(unsigned char *)($target_signature + 1))
  set $key1 = *(unsigned char *)(0x1007b788 + *(unsigned char *)($target_signature + 2))
  set $key2 = *(unsigned char *)(0x1007b850 + *(unsigned char *)($target_signature + 3))
  printf "HELLO_TARGET_SIGNATURE slot=%u bank=%u leaf=%u:", $slot, $bank, $leaf
  x/7ub $target_signature
  printf "HELLO_TARGET_CLASS_KEY slot=%u key=%u,%u,%u,%u,%u\n", $slot, $key0, $key1, $key2, *(unsigned char *)($target_signature + 5), (*(unsigned char *)($target_signature + 6) & 0x20)
  set $return = *(unsigned int *)$esp
  set $last_sum = -1
  tbreak *$return
  commands
    silent
    set $native = $al
    if $native == 0 && $last_sum > 2
      printf "HELLO_CUTOFF_OVERRIDE slot=%u sum=%u native=%u overridden=1\n", $slot, $last_sum, $native
      set $eax = 1
    end
    printf "HELLO_TARGET_SPLIT_RETURN slot=%u native=%u final=%u measured_sum=%d\n", $slot, $native, $al, $last_sum
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
  printf "HELLO_POST_COVERAGE_GATE call=%u context=%u total_positions=%u candidate_count=%u input_ids=", $postselection_calls, $ctx, $total, $count
  set $i = 0
  set $input_count = *(unsigned short *)$record
  while $i < $input_count && $i < 12
    printf "%u ", *(unsigned int *)($record + 4 + $i * 4)
    set $i = $i + 1
  end
  printf "\n"
  set $unit_signature = *(unsigned int *)($current_model + 0x64)
  printf "HELLO_MARKER_BYTES context=%u ids=272822:%#x,272823:%#x,272824:%#x,272825:%#x\n", $ctx, *(unsigned char *)($unit_signature + 272822 * 7 + 6), *(unsigned char *)($unit_signature + 272823 * 7 + 6), *(unsigned char *)($unit_signature + 272824 * 7 + 6), *(unsigned char *)($unit_signature + 272825 * 7 + 6)
  set $class_table = *(unsigned int *)($current_model + 0x30)
  printf "HELLO_SOURCE_CLASS_IDS context=%u 272822:%u,272823:%u,272824:%u,272825:%u\n", $ctx, *(unsigned int *)($class_table + 272822 * 4), *(unsigned int *)($class_table + 272823 * 4), *(unsigned int *)($class_table + 272824 * 4), *(unsigned int *)($class_table + 272825 * 4)
  set $i = 0
  set $span1 = 0
  set $span2 = 0
  set $span3 = 0
  set $span4 = 0
  set $target_mask = 0
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
      printf "HELLO_TARGET_CANDIDATE context=%u unit=%u class=%u left=%d right=%d span=%d\n", $ctx, $unit, *(unsigned int *)($class_table + $unit * 4), *(short *)($node + 0xc), *(short *)($node + 0x10), $span
    end
    set $i = $i + 1
  end
  printf "HELLO_SPANS context=%u scanned=%u span1=%u span2=%u span3=%u span4=%u target_mask=%#x\n", $ctx, $i, $span1, $span2, $span3, $span4, $target_mask
  continue
end

set $selected_calls = 0
break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  printf "HELLO_SELECTED_UNIT n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  continue
end

break *0x408187
commands
  silent
  printf "HELLO_CONTINUITY_TRACE_READY\n"
  continue
end

continue
