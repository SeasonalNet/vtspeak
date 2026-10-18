set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1005e330
commands
  silent
  printf "USERDICT_OPEN file="
  x/s *(unsigned int *)($esp + 4)
  continue
end

break *0x1005e4a4
commands
  silent
  printf "USERDICT_RECORD field_count=%d\n", $eax
  continue
end

break *0x1005e4d1
commands
  silent
  printf "USERDICT_FIELDS source="
  x/s *(unsigned int *)($ebp - 0x28)
  printf " target="
  x/s *(unsigned int *)($ebp - 0x14)
  printf " type="
  x/s $edx
  continue
end

break *0x1005e51a
commands
  silent
  printf "USERDICT_SOURCE_NORM result=%d\n", $eax
  continue
end

break *0x1005e680
commands
  silent
  printf "USERDICT_TARGET_PHON result=%d\n", $eax
  continue
end

break *0x1001da50
commands
  silent
  set $load = ((short (*)(int, char *))0x10027960)(18, "Z:/work/stage16/userdict-probe-quoted.csv")
  printf "USERDICT_LOAD_QUOTED index=18 return=%d\n", $load
  disable 6
  continue
end

continue
