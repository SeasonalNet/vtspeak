set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/legacy-map-comparison/legacy-class-span-gdb.log
set logging overwrite on
set logging enabled on
set $legacy_calls = 0

hbreak *0x1001cd60
commands
  silent
  set $legacy_calls = $legacy_calls + 1
  if $legacy_calls <= 4
    set $position = *(unsigned int *)($esp + 4)
    set $state = *(unsigned int *)($esp + 8)
    set $model = *(unsigned int *)($esp + 12)
    set $target = $state + 0xbb154 + $position * 0xfc
    set $target_count = *(unsigned short *)$target
    set $table = *(unsigned int *)($model + 0x30)
    set $unit_flags = *(unsigned int *)($model + 0x68)
    printf "LEGACY_TARGET_CLASS_POOL position=%u count=%u classes=", $position, $target_count
    set $i = 0
    while $i < $target_count && $i < 32
      printf "%u ", *(unsigned int *)($target + 4 + $i * 4)
      set $i = $i + 1
    end
    printf "\nLEGACY_MODEL_TABLE_30 position=%u 272822->%u 272823->%u 272824->%u 272825->%u\n", $position, *(unsigned int *)($table + 272822 * 4), *(unsigned int *)($table + 272823 * 4), *(unsigned int *)($table + 272824 * 4), *(unsigned int *)($table + 272825 * 4)
    printf "LEGACY_MODEL_UNIT_FLAGS position=%u 272822=%#x 272823=%#x 272824=%#x 272825=%#x\n", $position, *(unsigned char *)($unit_flags + 272822), *(unsigned char *)($unit_flags + 272823), *(unsigned char *)($unit_flags + 272824), *(unsigned char *)($unit_flags + 272825)
    set $return = *(unsigned int *)$esp
    thbreak *$return
    commands
      silent
      set $record = $state + 0xbb154 + $position * 0xfc
      set $count = *(unsigned short *)($record + 0xf4)
      set $candidate_array = $state + 0x528fc
      printf "LEGACY_SPAN_RESULT position=%u candidates=%u\n", $position, $count
      set $i = 0
      while $i < $count && $i < 12
        set $candidate = *(unsigned int *)($candidate_array + $i * 4)
        printf "LEGACY_SPAN_CANDIDATE index=%u unit=%u span=%u left=%u right=%u\n", $i, *(unsigned int *)($candidate + 8), *(unsigned short *)($candidate + 0x10), *(unsigned short *)($candidate + 0xc), *(unsigned short *)($candidate + 0xe)
        set $i = $i + 1
      end
      continue
    end
  end
  continue
end

break *0x408187
commands
  silent
  printf "LEGACY_CLASS_SPAN_TRACE_READY\n"
  continue
end

continue
