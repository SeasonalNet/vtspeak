set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/kate-key5-score-path-hello/signature5-swap-gdb.log
set logging overwrite on
set logging enabled on

# Swap only signature byte +5 between the matched position-2 pair while
# FUN_100182e0 runs. Candidate membership, other signature bytes, scalar
# feature arrays, and transition tables remain unchanged.
break *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  if $unit == 272823 || $unit == 264072
    set $model = *(unsigned int *)($esp + 24)
    set $signatures = *(unsigned int *)($model + 0x64)
    set $old_byte = *(unsigned char *)($signatures + 272823 * 7 + 5)
    set $new_byte = *(unsigned char *)($signatures + 264072 * 7 + 5)
    set {unsigned char}($signatures + 272823 * 7 + 5) = $new_byte
    set {unsigned char}($signatures + 264072 * 7 + 5) = $old_byte
    set $return = *(unsigned int *)$esp
    printf "KATE_KEY5_SIGNATURE5_SWAP_ENTRY unit=%u legacy=%u current=%u\n", $unit, $old_byte, $new_byte
    tbreak *$return
    commands
      silent
      printf "KATE_KEY5_SIGNATURE5_SWAP_RETURN unit=%u score=%g\n", $unit, $st0
      set {unsigned char}($signatures + 272823 * 7 + 5) = $old_byte
      set {unsigned char}($signatures + 264072 * 7 + 5) = $new_byte
      continue
    end
  else
    if $unit == 272823 || $unit == 264072 || $unit == 280485 || $unit == 21433 || $unit == 264073 || $unit == 272824 || $unit == 264074 || $unit == 272825
      set $return = *(unsigned int *)$esp
      printf "KATE_KEY5_SIGNATURE5_CONTROL_ENTRY unit=%u\n", $unit
      tbreak *$return
      commands
        silent
        printf "KATE_KEY5_SIGNATURE5_CONTROL_RETURN unit=%u score=%g\n", $unit, $st0
        continue
      end
    end
  end
  continue
end

# Capture the dynamic-programming rows after the +5 intervention. This shows
# whether the repaired local score also wins once transition costs accumulate.
break *0x10018c80
commands
  silent
  set $context = *(int *)($esp + 4)
  set $state = *(unsigned int *)($esp + 8)
  set $return = *(unsigned int *)$esp
  set $record = $state + 0xae894 + $context * 0xfc
  tbreak *$return
  commands
    silent
    set $count = *(short *)($record + 0xf4)
    set $previous = $record - 0xfc
    printf "KATE_KEY5_SIGNATURE5_PATH context=%d candidates=%d previous=%d", $context, $count, *(short *)($previous + 0xf4)
    set $index = 0
    while $index < $count && $index < 40
      set $id = *(unsigned int *)($record + 0x7c + $index * 4)
      set $cost = *(float *)($state + 0x8212c + ($context * 30 + $index) * 4)
      set $pred = *(short *)($state + 0x9f664 + ($context * 30 + $index) * 2)
      set $pred_id = *(unsigned int *)($previous + 0x7c + $pred * 4)
      printf " {%u cost=%g pred=%u}", $id, $cost, $pred_id
      set $index = $index + 1
    end
    printf "\n"
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  printf "KATE_KEY5_SIGNATURE5_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "KATE_KEY5_SIGNATURE5_TRACE_READY\n"
  continue
end

continue
