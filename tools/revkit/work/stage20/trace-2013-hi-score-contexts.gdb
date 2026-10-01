set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hi-2013-score-contexts/adapted-score-contexts-gdb.log
set logging overwrite on
set logging enabled on
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
      printf "HI_SCORE_CONTEXT_CUTOFF_OVERRIDE slot=%u sum=%d\n", $slot, $last_sum
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
  if $unit == 272820 || $unit == 65331 || $unit == 272821 || $unit == 280449 || $unit == 279499 || $unit == 279498
    set $context = *(unsigned int *)($esp + 12)
    printf "HI_SCORE_CONTEXT unit=%u bytes=%02x,%02x,%02x,%02x,%02x,%02x\n", $unit, *(unsigned char *)$context, *(unsigned char *)($context + 1), *(unsigned char *)($context + 2), *(unsigned char *)($context + 3), *(unsigned char *)($context + 4), *(unsigned char *)($context + 5)
    set $return = *(unsigned int *)$esp
    tbreak *$return
    commands
      silent
      printf "HI_SCORE_RESULT unit=%u score=%g\n", $unit, $st0
      continue
    end
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "HI_SCORE_CONTEXT_TRACE_READY\n"
  continue
end

continue
