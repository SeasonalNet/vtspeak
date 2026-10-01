set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/kate-key5-score-path-hello/feature35-swap-gdb.log
set logging overwrite on
set logging enabled on

# Swap only the two feature arrays previously shown to affect local cost,
# between old-selected rows and the rows chosen by this key5 adapter layout.
break *0x100182e0
commands
  silent
  set $unit = *(unsigned int *)($esp + 4)
  set $a = 0
  set $b = 0
  if $unit == 272822 || $unit == 280485
    set $a = 272822
    set $b = 280485
  else
    if $unit == 21433
      set $a = 272822
      set $b = 21433
    else
      if $unit == 272823 || $unit == 264072
        set $a = 272823
        set $b = 264072
      else
        if $unit == 272824 || $unit == 264073
          set $a = 272824
          set $b = 264073
        else
          if $unit == 272825 || $unit == 264074
            set $a = 272825
            set $b = 264074
          end
        end
      end
    end
  end
  if $a != 0
    set $model = *(unsigned int *)($esp + 24)
    set $p3 = *(unsigned int *)($model + 0x54)
    set $p5 = *(unsigned int *)($model + 0x5c)
    set $p3a = *(unsigned char *)($p3 + $a)
    set $p3b = *(unsigned char *)($p3 + $b)
    set $p5a = *(unsigned char *)($p5 + $a)
    set $p5b = *(unsigned char *)($p5 + $b)
    set {unsigned char}($p3 + $a) = $p3b
    set {unsigned char}($p3 + $b) = $p3a
    set {unsigned char}($p5 + $a) = $p5b
    set {unsigned char}($p5 + $b) = $p5a
    set $return = *(unsigned int *)$esp
    printf "KATE_KEY5_FEATURE35_SWAP_ENTRY unit=%u pair=%u,%u p3=%u/%u p5=%u/%u\n", $unit, $a, $b, $p3a, $p3b, $p5a, $p5b
    tbreak *$return
    commands
      silent
      printf "KATE_KEY5_FEATURE35_SWAP_RETURN unit=%u score=%g\n", $unit, $st0
      set {unsigned char}($p3 + $a) = $p3a
      set {unsigned char}($p3 + $b) = $p3b
      set {unsigned char}($p5 + $a) = $p5a
      set {unsigned char}($p5 + $b) = $p5b
      continue
    end
  end
  continue
end

break *0x1001b200
commands
  silent
  printf "KATE_KEY5_FEATURE35_SELECTED id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

hbreak *0x408187
commands
  silent
  printf "KATE_KEY5_FEATURE35_TRACE_READY\n"
  continue
end

continue
