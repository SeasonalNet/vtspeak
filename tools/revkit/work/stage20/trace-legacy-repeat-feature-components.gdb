set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hello-repeat-2006-feature-components-v2/legacy-feature-components-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001ea73
commands
  silent
  set $current = *(unsigned int *)($esp + 0x3c)
  set $previous = $ebp
  if ($current == 272823) && ($previous == 273369)
    set $category = *(signed char *)($esp + 0x13)
    set $weight0 = *(float *)(0x100707f0 + $category * 12)
    set $weight1 = *(float *)(0x100707f4 + $category * 12)
    set $weight2 = *(float *)(0x100707f8 + $category * 12)
    printf "LEGACY_FEATURE_COMPONENTS current=%u previous=%u category=%d metrics=%g,%g,%g weights=%g,%g,%g bias=%g weighted_sum=%g\n", $current, $previous, $category, *(float *)($esp + 0x68), *(float *)($esp + 0x74), *(float *)($esp + 0x64), $weight0, $weight1, $weight2, *(float *)0x1006116c, *(float *)($esp + 0x44)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_FEATURE_COMPONENTS_TRACE_READY\n"
  continue
end

continue
