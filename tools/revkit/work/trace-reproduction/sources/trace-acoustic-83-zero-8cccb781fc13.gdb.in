set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $active_unit_id = 0
set $active_unit = 0
set $desc_calls = 0

break *0x1000e0c0
commands
  silent
  set $code = *(unsigned char *)($esp + 8)
  if $code == 83
    printf "TPP_NUMERIC_ABLATE83 code=%u\n", $code
    set {unsigned char}($esp + 8) = 0
  end
  continue
end

break *0x1002d1e0
commands
  silent
  set $active_unit_id = *(unsigned int *)($esp + 16)
  set $active_unit = 1
  set $call_return = *(unsigned int *)$esp
  printf "ACOUSTIC_UNIT_BEGIN id=%u\n", $active_unit_id
  tbreak *$call_return
  commands
    silent
    set $active_unit = 0
    printf "ACOUSTIC_UNIT_END id=%u\n", $active_unit_id
    continue
  end
  continue
end

break *0x1002c120
commands
  silent
  if $active_unit
    set $desc_calls = $desc_calls + 1
    set $desc = *(unsigned int *)($esp + 4)
    set $desc_return = *(unsigned int *)$esp
    tbreak *$desc_return
    commands
      silent
      printf "ACOUSTIC_DESC id=%u n=%u mode=%u bank=%u dat_offset=%u dat_samples=%u upm_offset=%u upm_periods=%u edges=%u/%u\n", $active_unit_id, $desc_calls, *(unsigned char *)($desc + 0x26), *(unsigned short *)($desc + 0x12), *(unsigned int *)($desc + 0x14), *(unsigned int *)($desc + 0x0c), *(unsigned int *)($desc + 0x1c), *(unsigned char *)($desc + 0x20), *(unsigned short *)($desc + 0x22), *(unsigned short *)($desc + 0x24)
      continue
    end
  end
  continue
end

continue
