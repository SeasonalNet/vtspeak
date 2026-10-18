set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $apple_pair_calls = 0

# Sample the paired 2013 transition scorer for the same old/new candidate
# combinations captured in the native 2006 path trace.
hbreak *0x100191bc
commands
  silent
  set $pair_context = *(int *)($ebp + 8)
  set $pair_current = *(unsigned int *)($ebp - 8)
  set $pair_previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $pair_previous = *(unsigned int *)$pair_previous_ptr
  if (($pair_context == 2) && (($pair_current == 100180) || ($pair_current == 64719)) && ($pair_previous == 177774)) || (($pair_context == 3) && ((($pair_current == 59558) && ($pair_previous == 100180)) || (($pair_current == 2554) && ($pair_previous == 64719)) || (($pair_current == 82036) && ($pair_previous == 64719)))) || (($pair_context == 4) && ((($pair_current == 59559) && ($pair_previous == 59558)) || (($pair_current == 82037) && ($pair_previous == 82036)) || (($pair_current == 2555) && ($pair_previous == 2554))))
    set $apple_pair_calls = $apple_pair_calls + 1
    set $pair_feature = *(float *)($ebp - 0x4c)
    set $pair_divisor = *(int *)($ebp - 0x10)
    set $pair_penalty = *(int *)($ebp - 0xc)
    set $pair_previous_cost_ptr = *(unsigned int *)($ebp - 0x54)
    set $pair_previous_cost = *(float *)$pair_previous_cost_ptr
    set $pair_local_cost = *(float *)($ebp - 0x74)
    set $pair_raw_distance = *(float *)($ebp - 0x38)
    set $pair_feature_a = *(float *)($ebp - 0x34)
    set $pair_feature_b = *(float *)($ebp - 0x30)
    set $pair_weight_group = *(signed char *)($ebp - 1)
    set $pair_raw_weight = *(float *)(0x1007c288 + $pair_weight_group * 12)
    set $pair_a_weight = *(float *)(0x1007c290 + $pair_weight_group * 12)
    set $pair_b_weight = *(float *)(0x1007c28c + $pair_weight_group * 12)
    set $pair_current_metric_group = *(short *)($ebp - 0x14)
    set $pair_previous_metric_group = *(short *)($ebp - 0x28)
    set $pair_current_metric_array = *(unsigned int *)($ebx + $pair_current_metric_group * 4 + 0x3c)
    set $pair_previous_metric_array = *(unsigned int *)($ebx + $pair_previous_metric_group * 4 + 0x3c)
    set $pair_current_metric_word = *(unsigned short *)($pair_current_metric_array + $pair_current * 2)
    set $pair_previous_metric_word = *(unsigned short *)($pair_previous_metric_array + $pair_previous * 2)
    printf "APPLE_2013_SHARED_EDGE n=%u context=%d current=%u previous=%u previous_cost=%g local_cost=%g raw=%g derived_a=%g derived_b=%g weights_raw_a_b=%g,%g,%g feature_subtotal=%g divisor=%d penalty=%d metric_groups=%d,%d metric_words=%u,%u total=%g\n", $apple_pair_calls, $pair_context, $pair_current, $pair_previous, $pair_previous_cost, $pair_local_cost, $pair_raw_distance, $pair_feature_a, $pair_feature_b, $pair_raw_weight, $pair_a_weight, $pair_b_weight, $pair_feature, $pair_divisor, $pair_penalty, $pair_current_metric_group, $pair_previous_metric_group, $pair_current_metric_word, $pair_previous_metric_word, $pair_feature / $pair_divisor + $pair_penalty + $pair_previous_cost + $pair_local_cost
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_2013_SHARED_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2013_SHARED_PATH_TRACE_READY\n"
  continue
end

continue
