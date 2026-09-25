set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1003a7b0
commands
  silent
  set $out = *(unsigned int *)($esp + 4)
  set $key = *(unsigned int *)($esp + 8)
  set $selector = *(unsigned int *)($esp + 12)
  if $selector == 0x46
    set $length = 0
    while *(unsigned char *)($key + $length) != 0 && $length < 64
      set $length = $length + 1
    end
    set $starts_casa = *(unsigned int *)$key == 0x61736143
    set $starts_saint = *(unsigned int *)$key == 0x6e696153
    if $length == 23 && ($starts_casa || $starts_saint)
      printf "FORCED_E_LOOKUP original_selector=%#x key=", $selector
      x/s $key
      set {unsigned int}($esp + 12) = 0x45
      set $return_address = *(unsigned int *)$esp
      tbreak *$return_address
      commands
        silent
        printf "FORCED_E_RESULT eax=%#x output=", $eax
        x/s $out
        printf "FORCED_E_BYTES:"
        x/16bx $out
        set {unsigned char}$out = 0
        set $eax = 0
        continue
      end
    end
  end
  continue
end

continue
