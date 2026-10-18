set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $load = ((short (*)(int, char *))0x10027960)(19, "Z:/work/stage16/userdict-alpha.csv")
  printf "USERDICT_LOAD_ALPHA index=19 filename=Z:/work/stage16/userdict-alpha.csv return=%d\n", $load
  if $load == 1
    set $unload = ((short (*)(int))0x10027a80)(19)
    printf "USERDICT_UNLOAD_ALPHA index=19 return=%d\n", $unload
  end
  continue
end

continue
