set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-repeat-2013-feature-components/adapted-feature-components-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x100191bc
commands
  silent
  set $context = *(int *)($ebp + 8)
  set $current = *(unsigned int *)($ebp - 8)
  set $previous_ptr = *(unsigned int *)($ebp - 0x18)
  set $previous = *(unsigned int *)$previous_ptr
  if ($context == 1) && ($current == 272823) && ($previous == 273369)
    set $category = *(signed char *)($ebp - 1)
    set $weight0 = *(float *)(0x1007c288 + $category * 12)
    set $weight1 = *(float *)(0x1007c28c + $category * 12)
    set $weight2 = *(float *)(0x1007c290 + $category * 12)
    set $metric0 = *(float *)($ebp - 0x38)
    set $metric1 = *(float *)($ebp - 0x34)
    set $metric2 = *(float *)($ebp - 0x30)
    set $feature = *(float *)($ebp - 0x4c)
    printf "NEW_FEATURE_COMPONENTS context=%d current=%u previous=%u category=%d metrics=%g,%g,%g weights=%g,%g,%g bias=%g weighted_sum=%g\n", $context, $current, $previous, $category, $metric0, $metric1, $metric2, $weight0, $weight1, $weight2, *(float *)0x1006d174, $feature
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "NEW_FEATURE_COMPONENTS_TRACE_READY\n"
  continue
end

continue
