set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $target_norm = ((short (*)(char *))0x1002a570)("hello")
  set $target_phon = ((short (*)(char *))0x1002a590)("HH")
  printf "USERDICT_VALIDATE target_norm_hello=%d target_phon_HH=%d\n", $target_norm, $target_phon
  set $control_text = "vtspeakprobe"
  set $control_path = "Z:/work/stage16/userdict-control.wav"
  set $control_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $control_text, $control_path, 1, -1, -1, -1, -1, -1, -1)
  printf "USERDICT_CONTROL text=vtspeakprobe dictidx=-1 return=%d\n", $control_result
  set $load_result = ((short (*)(int, char *))0x10027960)(17, "Z:/work/stage16/userdict-probe.csv")
  printf "USERDICT_LOAD index=17 filename=Z:/work/stage16/userdict-probe.csv return=%d\n", $load_result
  set $dict_text = "vtspeakprobe"
  set $dict_path = "Z:/work/stage16/userdict-enabled.wav"
  set $dict_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $dict_text, $dict_path, 1, -1, -1, -1, -1, 17, -1)
  printf "USERDICT_SYNTH text=vtspeakprobe dictidx=17 return=%d\n", $dict_result
  set $unload_result = ((short (*)(int))0x10027a80)(17)
  printf "USERDICT_UNLOAD index=17 return=%d\n", $unload_result
  continue
end

continue
