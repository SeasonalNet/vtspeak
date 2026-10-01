set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/kate-key5-score-path-hello/score-fields-gdb.log
set logging overwrite on
set logging enabled on

# Capture target/candidate signatures and the six scalar feature arrays used
# by FUN_100182e0 for legacy-selected rows and the current selected rows.
break *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if (($unit >= 272822) && ($unit <= 272825)) || (($unit == 280485) || ($unit == 21433) || ($unit >= 264072 && $unit <= 264074))
    set $model = *(unsigned int *)($esp + 24)
    set $target = *(unsigned int *)($esp + 12)
    set $feature_view = *(unsigned int *)($esp + 16)
    set $signature = *(unsigned int *)($model + 0x64) + $unit * 7
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
    printf "KATE_KEY5_SCORE_ENTRY unit=%u target=", $unit
    x/7ub $target
    printf "KATE_KEY5_SCORE_TARGET_VIEW:"
    x/10ub $feature_view
    printf "KATE_KEY5_SCORE_SIGNATURE:"
    x/7ub $signature
    printf "KATE_KEY5_SCORE_ARRAYS attr_a=%u attr_b=%u metrics=", *(unsigned char *)($attr_a + $unit), *(unsigned char *)($attr_b + $unit)
    printf "%u,%u,%u features=", *(unsigned short *)($metric_0 + $unit * 2), *(unsigned short *)($metric_1 + $unit * 2), *(unsigned short *)($metric_2 + $unit * 2)
    printf "%u,%u,%u,%u,%u,%u\n", *(unsigned char *)($feature_0 + $unit), *(unsigned char *)($feature_1 + $unit), *(unsigned char *)($feature_2 + $unit), *(unsigned char *)($feature_3 + $unit), *(unsigned char *)($feature_4 + $unit), *(unsigned char *)($feature_5 + $unit)
    tbreak *$return
    commands
      silent
      printf "KATE_KEY5_SCORE_RETURN unit=%u score=%g\n", $unit, $st0
      continue
    end
  end
  continue
end

break *0x1001b200
commands
  silent
  printf "KATE_KEY5_SCORE_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "KATE_KEY5_SCORE_FIELDS_TRACE_READY\n"
  continue
end

continue
