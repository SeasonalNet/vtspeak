set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/apple-selector-byte0-old-shared-path-components-2026-09-29/legacy-shared-path-components-gdb.log
set logging overwrite on
set logging enabled on

# FUN_1001e470 computes the legacy pair-feature term before adding it to the
# previous path cost and local unit score. Sample the aligned Apple edges.
hbreak *0x1001ea73
commands
  silent
  set $context = *(int *)($esp + 0x2c)
  set $current = *(unsigned int *)($esp + 0x3c)
  set $previous = $ebp
  if (($context == 2) && (($current == 100180) || ($current == 64719)) && ($previous == 177774)) || (($context == 3) && ((($current == 59558) && ($previous == 100180)) || (($current == 2554) && ($previous == 64719)) || (($current == 82036) && ($previous == 64719)))) || (($context == 4) && ((($current == 59559) && ($previous == 59558)) || (($current == 82037) && ($previous == 82036)) || (($current == 2555) && ($previous == 2554))))
    set $category = *(signed char *)($esp + 0x13)
    set $weight0 = *(float *)(0x100707f0 + $category * 12)
    set $weight1 = *(float *)(0x100707f4 + $category * 12)
    set $weight2 = *(float *)(0x100707f8 + $category * 12)
    printf "APPLE_2006_SHARED_EDGE context=%d current=%u previous=%u category=%d metrics=%g,%g,%g weights=%g,%g,%g bias=%g weighted_sum=%g\n", $context, $current, $previous, $category, *(float *)($esp + 0x68), *(float *)($esp + 0x74), *(float *)($esp + 0x64), $weight0, $weight1, $weight2, *(float *)0x1006116c, *(float *)($esp + 0x44)
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "APPLE_2006_SHARED_PATH_COMPONENT_TRACE_READY\n"
  continue
end

continue
