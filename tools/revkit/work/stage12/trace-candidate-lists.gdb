set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10023350
commands
  silent
  set $context = *(int *)($esp + 4)
  set $engine = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  tbreak *$ret
  commands
    silent
    set $count = *(short *)($engine + 0xae988 + $context * 0xfc)
    set $list = $engine + 0xae910 + $context * 0xfc
    printf "FUN23350_CANDIDATES context=%d count=%d\n", $context, $count
    set $i = 0
    while $i < $count && $i < 100
      printf "FUN23350_UNIT context=%d index=%d unit_index=%u\n", $context, $i, *(unsigned int *)($list + $i * 4)
      set $i = $i + 1
    end
    continue
  end
  continue
end

break *0x10024680
commands
  silent
  set $engine = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  tbreak *$ret
  commands
    silent
    set $context_count = *(unsigned int *)($engine + 0xec624)
    printf "FUN24680_CONTEXTS count=%u\n", $context_count
    set $context = 0
    while $context < $context_count && $context < 100
      set $base = $engine + 0xae894 + $context * 0xfc
      set $class_count = *(short *)$base
      printf "CLASS_CANDIDATES context=%d count=%d\n", $context, $class_count
      set $i = 0
      while $i < $class_count && $i < 100
        printf "CLASS_ID context=%d index=%d value=%u\n", $context, $i, *(unsigned int *)($base + 4 + $i * 4)
        set $i = $i + 1
      end
      set $context = $context + 1
    end
    continue
  end
  continue
end

continue
