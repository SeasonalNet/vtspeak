set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/hello-index-field-origin-gdb.log
set logging overwrite on
set logging enabled on

set $index_calls = 0
break *0x10019940
commands
  silent
  set $index_calls = $index_calls + 1
  set $index = *(unsigned int *)($esp + 4)
  set $return = *(unsigned int *)$esp
  set $count = *(unsigned int *)($index + 0x4c)
  printf "INDEX_ENTRY call=%u units=%u object=%#x col_a=%#x key=%#x col_b=%#x\n", $index_calls, $count, $index, *(unsigned int *)($index + 0x48), *(unsigned int *)($index + 0x44), *(unsigned int *)($index + 0x40)
  tbreak *$return
  commands
    silent
    printf "INDEX_RETURN call=%u col_a_prefix:", $index_calls
    x/8bx *(unsigned int *)($index + 0x48)
    printf "INDEX_RETURN call=%u key_prefix:", $index_calls
    x/16bx *(unsigned int *)($index + 0x44)
    printf "INDEX_RETURN call=%u col_b_prefix:", $index_calls
    x/8bx *(unsigned int *)($index + 0x40)
    if $index_calls == 4
      set $row = 1286
      printf "INDEX_RETURN row=%u col_a:", $row
      x/1bx (*(unsigned int *)($index + 0x48) + $row)
      printf "INDEX_RETURN row=%u signature:", $row
      x/7bx (*(unsigned int *)($index + 0x44) + $row * 7)
      printf "INDEX_RETURN row=%u col_b:", $row
      x/1bx (*(unsigned int *)($index + 0x40) + $row)
      set $watch_addr = *(unsigned int *)($index + 0x48) + $row
      set $watch_hits = 0
      watch *(unsigned char *)$watch_addr
      commands 7
        silent
        set $watch_hits = $watch_hits + 1
        printf "ATTR_A_WRITE row=%u address=%#x value=%u pc=%#x hits=%u\n", $row, $watch_addr, *(unsigned char *)$watch_addr, $pc, $watch_hits
        bt 8
        if $watch_hits >= 4
          disable 7
        end
        continue
      end
    end
    continue
  end
  continue
end

set $score_calls = 0
break *0x100182e0
commands
  silent
  set $score_calls = $score_calls + 1
  if $score_calls <= 80
    set $unit = *(unsigned int *)($esp + 4)
    set $state = *(unsigned int *)($esp + 24)
    set $sig_ptr = *(unsigned int *)($state + 0x64)
    set $a_ptr = *(unsigned int *)($state + 0x68)
    set $b_ptr = *(unsigned int *)($state + 0x60)
    printf "SCORE_FIELD_ORIGIN call=%u unit=%u state=%#x sig_ptr=%#x attr_a_ptr=%#x attr_b_ptr=%#x values_a_b=", $score_calls, $unit, $state, $sig_ptr, $a_ptr, $b_ptr
    printf "%u,%u signature=", *(unsigned char *)($a_ptr + $unit), *(unsigned char *)($b_ptr + $unit)
    x/7bx ($sig_ptr + $unit * 7)
  end
  continue
end

continue
