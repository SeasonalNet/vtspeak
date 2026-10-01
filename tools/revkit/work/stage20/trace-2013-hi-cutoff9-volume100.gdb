set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-query-cutoff9-volume100-gdb.log
set logging overwrite on
set logging enabled on

set $hi_metric = -1
set $hi_selected = 0
set $hi_rows = 0

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
    printf "HI_VOLUME100_METRIC_INPUT count=%u id_weight=", $count
    set $i = 0
    while $i < $count && $i < 96
      set $id = *(unsigned int *)($ids + $i * 4)
      printf "%u:%u ", $id, *(unsigned short *)($weights + $id * 2)
      set $i = $i + 1
    end
    printf "\n"
    tbreak *$return
    commands
      silent
      set $hi_metric = $eax
      printf "HI_VOLUME100_METRIC_SUM value=%u\n", $eax
      continue
    end
  end
  continue
end

break *0x10024060
commands
  silent
  set $hi_metric = -1
  set $slot = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    set $native = $al
    if $slot == 1 && $hi_metric == 9 && $native == 0
      printf "HI_VOLUME100_CUTOFF_OVERRIDE metric=%u\n", $hi_metric
      set $eax = 1
    end
    continue
  end
  continue
end

hbreak *0x10027fe0
commands
  silent
  set $old_volume = *(int *)($esp + 12)
  printf "HI_VOLUME100_SETTER original_volume=%d replacement_volume=100\n", $old_volume
  set *(int *)($esp + 12) = 100
  continue
end

hbreak *0x1002c8b0
commands
  silent
  set $row = *(unsigned int *)($esp + 8)
  set $hi_rows = $hi_rows + 1
  printf "HI_VOLUME100_DECODED_ROW call=%u samples=%d sample_scale=%d\n", $hi_rows, *(int *)($row + 0xc), *(int *)($row + 8)
  continue
end

hbreak *0x10024680
commands
  silent
  set $return = *(unsigned int *)$esp
  tbreak *$return
  commands
    silent
    printf "HI_VOLUME100_POSITION_BUILDER positions=%u\n", $eax
    continue
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  set $hi_selected = $hi_selected + 1
  printf "HI_VOLUME100_SELECTED n=%u id=%u\n", $hi_selected, *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_CUTOFF9_VOLUME100_TRACE_READY\n"
  continue
end

continue
