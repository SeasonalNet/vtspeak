set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-plain-crosswalk-score-fields/adapted-score-fields-gdb.log
set logging overwrite on
set logging enabled on

# Retain the Stage 20 >2 split intervention. The per-position scorer is
# evaluated only after the same four-position chain is constructed.
set $last_sum = -1
break *0x10024060
commands
  silent
  set $last_sum = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    if $al == 0 && $last_sum > 2
      printf "SCORE_FIELDS_CUTOFF_OVERRIDE slot=%u sum=%d\n", $slot, $last_sum
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
    tbreak *$return
    commands
      silent
      set $last_sum = $eax
      continue
    end
  end
  continue
end

break *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if ($unit == 272822 || $unit == 272823 || $unit == 272824 || $unit == 272825 || $unit == 52542 || $unit == 273369 || $unit == 273370 || $unit == 273371 || $unit == 255282 || $unit == 264072 || $unit == 264073 || $unit == 264074)
    set $model = *(unsigned int *)($esp + 24)
    set $target = *(unsigned int *)($esp + 12)
    set $feature_view = *(unsigned int *)($esp + 16)
    set $candidate_signature = *(unsigned int *)($model + 0x64) + $unit * 7
    set $attr_a = *(unsigned int *)($model + 0x68)
    set $attr_b = *(unsigned int *)($model + 0x60)
    set $metric_0 = *(unsigned int *)($model + 0x3c)
    set $metric_1 = *(unsigned int *)($model + 0x40)
    set $metric_2 = *(unsigned int *)($model + 0x44)
    set $feature_0 = *(unsigned int *)($model + 0x48)
    set $feature_1 = *(unsigned int *)($model + 0x4c)
    set $feature_2 = *(unsigned int *)($model + 0x50)
    set $feature_3 = *(unsigned int *)($model + 0x54)
    set $feature_4 = *(unsigned int *)($model + 0x58)
    set $feature_5 = *(unsigned int *)($model + 0x5c)
    set $return = *(unsigned int *)$esp
    printf "SCORE_FIELDS_ENTRY unit=%u target=", $unit
    x/7ub $target
    printf "SCORE_FIELDS_TARGET_VIEW:"
    x/10ub $feature_view
    printf "SCORE_FIELDS_CANDIDATE_SIGNATURE:"
    x/7ub $candidate_signature
    printf "SCORE_FIELDS_ATTRS attr_a=%u attr_b=%u metrics=", *(unsigned char *)($attr_a + $unit), *(unsigned char *)($attr_b + $unit)
    printf "%u,%u,%u features=", *(unsigned short *)($metric_0 + $unit * 2), *(unsigned short *)($metric_1 + $unit * 2), *(unsigned short *)($metric_2 + $unit * 2)
    printf "%u,%u,%u,%u,%u,%u\n", *(unsigned char *)($feature_0 + $unit), *(unsigned char *)($feature_1 + $unit), *(unsigned char *)($feature_2 + $unit), *(unsigned char *)($feature_3 + $unit), *(unsigned char *)($feature_4 + $unit), *(unsigned char *)($feature_5 + $unit)
    tbreak *$return
    commands
      silent
      printf "SCORE_FIELDS_RETURN unit=%u score=%g\n", $unit, $st0
      continue
    end
  end
  continue
end

set $selected_calls = 0
break *0x1001b200
commands
  silent
  set $selected_calls = $selected_calls + 1
  if $selected_calls <= 16
    printf "SCORE_FIELDS_SELECTED n=%u id=%u\n", $selected_calls, *(unsigned int *)($esp + 8)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "SCORE_FIELDS_TRACE_READY\n"
  continue
end

continue
