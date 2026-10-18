set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  set $text = *(unsigned int *)($esp + 8)
  disable 1
  set $path = ((char *(*)(unsigned int))0x1001d9c0)(64)
  set {unsigned char}($path + 0) = 90
  set {unsigned char}($path + 1) = 58
  set {unsigned char}($path + 2) = 47
  set {unsigned char}($path + 3) = 119
  set {unsigned char}($path + 4) = 111
  set {unsigned char}($path + 5) = 114
  set {unsigned char}($path + 6) = 107
  set {unsigned char}($path + 7) = 47
  set {unsigned char}($path + 8) = 115
  set {unsigned char}($path + 9) = 116
  set {unsigned char}($path + 10) = 97
  set {unsigned char}($path + 11) = 103
  set {unsigned char}($path + 12) = 101
  set {unsigned char}($path + 13) = 50
  set {unsigned char}($path + 14) = 49
  set {unsigned char}($path + 15) = 47
  set {unsigned char}($path + 16) = 112
  set {unsigned char}($path + 17) = 99
  set {unsigned char}($path + 18) = 109
  set {unsigned char}($path + 19) = 0
  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=00 pitch=-1 speed=-1 volume=-1 pause=-1 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v00.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, 50, -1, -1, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=01 pitch=50 speed=-1 volume=-1 pause=-1 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v01.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, 200, -1, -1, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=02 pitch=200 speed=-1 volume=-1 pause=-1 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v02.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, 50, -1, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=03 pitch=-1 speed=50 volume=-1 pause=-1 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v03.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, 400, -1, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=04 pitch=-1 speed=400 volume=-1 pause=-1 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v04.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, -1, 0, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=05 pitch=-1 speed=-1 volume=0 pause=-1 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v05.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, -1, 500, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=06 pitch=-1 speed=-1 volume=500 pause=-1 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v06.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, -1, -1, 0, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=07 pitch=-1 speed=-1 volume=-1 pause=0 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v07.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, -1, -1, 250, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=08 pitch=-1 speed=-1 volume=-1 pause=250 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v08.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, -1, -1, 65535, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=09 pitch=-1 speed=-1 volume=-1 pause=65535 raw=%#x\n", $raw
  shell cp /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v09.pcm

  set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 6, 1, -1, -1, -1, -1, 0, 0)
  printf "PREPROCESS_PCM_SETTINGS variant=10 pitch=-1 speed=-1 volume=-1 pause=-1 raw=%#x\n", $raw
  shell mv /work/stage5/test.pcm /work/stage21/preprocess-pcm-settings-v10.pcm

  continue
end

continue
