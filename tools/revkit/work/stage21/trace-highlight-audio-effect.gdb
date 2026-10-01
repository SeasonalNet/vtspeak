set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $return_address = *(unsigned int *)$esp
  set $format = *(int *)($esp + 4)
  set $text = *(unsigned int *)($esp + 8)
  set $path = *(unsigned int *)($esp + 12)
  set $speaker = *(int *)($esp + 16)
  set $arg5 = *(int *)($esp + 20)
  set $arg6 = *(int *)($esp + 24)
  set $arg7 = *(int *)($esp + 28)
  set $arg8 = *(int *)($esp + 32)
  set $arg9 = *(int *)($esp + 36)
  set $arg10 = *(int *)($esp + 40)
  set $runtime_state = *(unsigned int *)0x100a0460
  set $highlight_initial = *(unsigned char *)($runtime_state + 0x20424)
  set $path_off = (char *)malloc(64)
  set {char[64]}$path_off = "highlight-off.wav"
  set $path_on = (char *)malloc(64)
  set {char[64]}$path_on = "highlight-on.wav"
  printf "HIGHLIGHT_AUDIO_INPUT format=%d text=%#x path=%s speaker=%d args=%d,%d,%d,%d,%d,%d initial=%u\n", $format, $text, (char *)$path_off, $speaker, $arg5, $arg6, $arg7, $arg8, $arg9, $arg10, $highlight_initial
  call ((void (*)(unsigned char))0x10028360)(0)
  set {unsigned int}($esp + 12) = (unsigned int)$path_off
  tbreak *$return_address
  commands
    silent
    printf "HIGHLIGHT_AUDIO_OFF_RETURN result=%d flag=%u\n", (short)$eax, *(unsigned char *)($runtime_state + 0x20424)
    call ((void (*)(unsigned char))0x10028360)(1)
    set $result_on = ((int (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)($format, (char *)$text, $path_on, $speaker, $arg5, $arg6, $arg7, $arg8, $arg9, $arg10)
    printf "HIGHLIGHT_AUDIO_ON_RETURN result=%d flag=%u\n", $result_on, *(unsigned char *)($runtime_state + 0x20424)
    call ((void (*)(unsigned char))0x10028360)($highlight_initial)
    printf "HIGHLIGHT_AUDIO_RESTORED value=%u\n", *(unsigned char *)($runtime_state + 0x20424)
    continue
  end
  continue
end

continue
