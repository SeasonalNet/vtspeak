set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/legacy-map-comparison/legacy-unit-map-gdb.log
set logging overwrite on
set logging enabled on

hbreak *0x1001cd60
commands
  silent
  set $position = *(unsigned int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $model = *(unsigned int *)($esp + 12)
  set $table = *(unsigned int *)($model + 0x30)
  if $position < 4
    printf "LEGACY_MODEL_TABLE_30 position=%u unit272822=%u unit272823=%u unit272824=%u unit272825=%u\n", $position, *(unsigned int *)($table + 272822 * 4), *(unsigned int *)($table + 272823 * 4), *(unsigned int *)($table + 272824 * 4), *(unsigned int *)($table + 272825 * 4)
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_UNIT_MAP_TRACE_READY\n"
  continue
end

continue
