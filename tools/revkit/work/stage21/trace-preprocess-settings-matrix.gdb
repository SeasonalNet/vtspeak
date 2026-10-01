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
  set {unsigned char}($path + 17) = 116
  set {unsigned char}($path + 18) = 102
  set {unsigned char}($path + 19) = 51
  set {unsigned char}($path + 20) = 118
  set {unsigned char}($path + 21) = 48
  set {unsigned char}($path + 22) = 48
  set {unsigned char}($path + 23) = 0

  set $variant = 0
  while $variant < 11
    set $pitch = -1
    set $speed = -1
    set $volume = -1
    set $pause = -1
    if $variant == 1
      set $pitch = 50
    end
    if $variant == 2
      set $pitch = 200
    end
    if $variant == 3
      set $speed = 50
    end
    if $variant == 4
      set $speed = 400
    end
    if $variant == 5
      set $volume = 0
    end
    if $variant == 6
      set $volume = 500
    end
    if $variant == 7
      set $pause = 0
    end
    if $variant == 8
      set $pause = 250
    end
    if $variant == 9
      set $pause = 65535
    end

    set {unsigned char}($path + 19) = 51
    set {unsigned char}($path + 21) = 48 + ($variant / 10)
    set {unsigned char}($path + 22) = 48 + ($variant % 10)
    set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 3, 1, $pitch, $speed, $volume, $pause, 0, 0)
    printf "PREPROCESS_SETTINGS flag=3 variant=%d pitch=%d speed=%d volume=%d pause=%d raw=%#x path=%s\n", $variant, $pitch, $speed, $volume, $pause, $raw, $path

    set {unsigned char}($path + 19) = 53
    set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 5, 1, $pitch, $speed, $volume, $pause, 0, 0)
    printf "PREPROCESS_SETTINGS flag=5 variant=%d pitch=%d speed=%d volume=%d pause=%d raw=%#x path=%s\n", $variant, $pitch, $speed, $volume, $pause, $raw, $path

    set {unsigned char}($path + 19) = 55
    set {unsigned int}0x1007cfb0 = 0
    set {unsigned int}0x1009fc4c = 8
    set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 7, 1, $pitch, $speed, $volume, $pause, 0, 0)
    printf "PREPROCESS_SETTINGS flag=7 variant=%d pitch=%d speed=%d volume=%d pause=%d state=%u raw=%#x path=%s\n", $variant, $pitch, $speed, $volume, $pause, *(unsigned int *)0x1009fc4c, $raw, $path

    set {unsigned char}($path + 19) = 97
    set $raw = ((unsigned int (*)(char *, char *, unsigned char, int, int, int, int, int, int, int))0x1001dd60)((char *)$text, $path, 10, 1, $pitch, $speed, $volume, $pause, 0, 0)
    printf "PREPROCESS_SETTINGS flag=10 variant=%d pitch=%d speed=%d volume=%d pause=%d raw=%#x path=%s\n", $variant, $pitch, $speed, $volume, $pause, $raw, $path

    set $variant = $variant + 1
  end
  continue
end

continue
