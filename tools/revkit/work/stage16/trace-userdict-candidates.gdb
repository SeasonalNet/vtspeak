set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $plain = ((short (*)(int, char *))0x10027960)(20, "Z:/work/stage16/userdict-hello-plain.csv")
  printf "USERDICT_LOAD plain index=20 low_ax=%d\n", $plain
  set $quoted = ((short (*)(int, char *))0x10027960)(21, "Z:/work/stage16/userdict-hello-quoted.csv")
  printf "USERDICT_LOAD quoted index=21 low_ax=%d\n", $quoted
  set $fourfield = ((short (*)(int, char *))0x10027960)(22, "Z:/work/stage16/userdict-hello-fourfield.csv")
  printf "USERDICT_LOAD fourfield index=22 low_ax=%d\n", $fourfield
  continue
end

continue
