set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  set $path = *(unsigned int *)($esp + 12)
  set $speaker = *(int *)($esp + 16)
  set $pitch = *(int *)($esp + 20)
  set $speed = *(int *)($esp + 24)
  set $volume = *(int *)($esp + 28)
  set $pause = *(int *)($esp + 32)
  set $dict = *(int *)($esp + 36)
  disable 1

  set $alias = (char *)malloc(64)
  set {char}($alias + 0) = 60
  set {char}($alias + 1) = 118
  set {char}($alias + 2) = 116
  set {char}($alias + 3) = 109
  set {char}($alias + 4) = 108
  set {char}($alias + 5) = 95
  set {char}($alias + 6) = 115
  set {char}($alias + 7) = 117
  set {char}($alias + 8) = 98
  set {char}($alias + 9) = 32
  set {char}($alias + 10) = 97
  set {char}($alias + 11) = 108
  set {char}($alias + 12) = 105
  set {char}($alias + 13) = 97
  set {char}($alias + 14) = 115
  set {char}($alias + 15) = 61
  set {char}($alias + 16) = 34
  set {char}($alias + 17) = 72
  set {char}($alias + 18) = 101
  set {char}($alias + 19) = 108
  set {char}($alias + 20) = 108
  set {char}($alias + 21) = 111
  set {char}($alias + 22) = 32
  set {char}($alias + 23) = 119
  set {char}($alias + 24) = 111
  set {char}($alias + 25) = 114
  set {char}($alias + 26) = 108
  set {char}($alias + 27) = 100
  set {char}($alias + 28) = 46
  set {char}($alias + 29) = 34
  set {char}($alias + 30) = 62
  set {char}($alias + 31) = 120
  set {char}($alias + 32) = 60
  set {char}($alias + 33) = 47
  set {char}($alias + 34) = 118
  set {char}($alias + 35) = 116
  set {char}($alias + 36) = 109
  set {char}($alias + 37) = 108
  set {char}($alias + 38) = 95
  set {char}($alias + 39) = 115
  set {char}($alias + 40) = 117
  set {char}($alias + 41) = 98
  set {char}($alias + 42) = 62
  set {char}($alias + 43) = 0
  printf "MARKUP_SAMPLE=<vtml_sub alias=Hello-world.>x</vtml_sub>\n"

  set $plain_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, (char *)$text, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 0)
  printf "MARKUP_PLAIN return=%d\n", $plain_result
  shell cp /work/stage5/output.wav /work/stage16/markup-plain.wav

  set $jeita_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $alias, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 4)
  printf "MARKUP_JEITA return=%d\n", $jeita_result
  shell cp /work/stage5/output.wav /work/stage16/markup-jeita.wav

  set $jeita_plus_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $alias, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 6)
  printf "MARKUP_JEITA_PLUS return=%d\n", $jeita_plus_result
  shell cp /work/stage5/output.wav /work/stage16/markup-jeita-plus.wav

  set $umd_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $alias, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 8)
  printf "MARKUP_UMD return=%d\n", $umd_result
  shell cp /work/stage5/output.wav /work/stage16/markup-umd.wav

  set $child = (char *)malloc(2)
  set {char}($child + 0) = 120
  set {char}($child + 1) = 0
  set $child_result = ((short (*)(int, char *, char *, int, int, int, int, int, int, int))0x1001da50)(4, $child, (char *)$path, $speaker, $pitch, $speed, $volume, $pause, $dict, 0)
  printf "MARKUP_CHILD_ONLY return=%d\n", $child_result
  shell cp /work/stage5/output.wav /work/stage16/markup-child-only.wav
  continue
end

continue
