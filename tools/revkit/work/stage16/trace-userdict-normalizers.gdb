set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $source_norm = (char *)malloc(256)
  set $source_result = ((int (*)(char *, char *))0x1005f2e0)($source_norm, "vtspeakprobe")
  printf "USERDICT_SOURCE input=vtspeakprobe result=%d normalized=", $source_result
  x/s $source_norm
  set $target_norm_result = ((short (*)(char *))0x1005f3b0)("hello")
  printf "USERDICT_TARGET_A input=hello result=%d\n", $target_norm_result
  set $target_phon_result = ((short (*)(char *))0x1005f5c0)("HH")
  printf "USERDICT_TARGET_P_VALIDATE input=HH result=%d\n", $target_phon_result
  set $target_phon = (char *)malloc(256)
  set $target_phon_convert = ((short (*)(char *, char *))0x1005f710)($target_phon, "HH")
  printf "USERDICT_TARGET_P_CONVERT input=HH result=%d output=", $target_phon_convert
  x/s $target_phon
  continue
end

continue
