set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x1001da50
commands
  silent
  disable 1
  set $pair = (unsigned char *)malloc(3)
  set {unsigned char}($pair + 2) = 0
  set $total = 0
  set $hash = 2166136261
  set $lead = 0xa1
  while $lead <= 0xad
    set $trail = 0xa1
    set $row_count = 0
    while $trail <= 0xfe
      set {unsigned char}($pair) = $lead
      set {unsigned char}($pair + 1) = $trail
      set $accepted = ((int (*)(unsigned char *))0x1001c900)($pair)
      if $accepted != 1
        printf "CLASSIFIER_MISMATCH lead=%02x trail=%02x expected=1 actual=%d\n", $lead, $trail, $accepted
      end
      set $row_count = $row_count + $accepted
      set $total = $total + $accepted
      set $hash = (($hash ^ $accepted) * 16777619) & 0xffffffff
      set $trail = $trail + 1
    end
    printf "CLASSIFIER_ROW lead=%02x accepted=%d\n", $lead, $row_count
    set $lead = $lead + 1
  end
  set {unsigned char}($pair) = 0xae
  set $trail = 0xa1
  set $row_count = 0
  while $trail <= 0xc2
    set {unsigned char}($pair + 1) = $trail
    set $accepted = ((int (*)(unsigned char *))0x1001c900)($pair)
    if $accepted != 1
      printf "CLASSIFIER_MISMATCH lead=ae trail=%02x expected=1 actual=%d\n", $trail, $accepted
    end
    set $row_count = $row_count + $accepted
    set $total = $total + $accepted
    set $hash = (($hash ^ $accepted) * 16777619) & 0xffffffff
    set $trail = $trail + 1
  end
  printf "CLASSIFIER_ROW lead=ae accepted=%d\n", $row_count
  set {unsigned char}($pair) = 0xfd
  set {unsigned char}($pair + 1) = 0xfe
  set $accepted = ((int (*)(unsigned char *))0x1001c900)($pair)
  if $accepted != 1
    printf "CLASSIFIER_MISMATCH lead=fd trail=fe expected=1 actual=%d\n", $accepted
  end
  set $total = $total + $accepted
  set $hash = (($hash ^ $accepted) * 16777619) & 0xffffffff
  printf "CLASSIFIER_SPECIAL lead=fd trail=fe accepted=%d\n", $accepted
  printf "CLASSIFIER_ACCEPTED_SET total=%d expected=1257 digest=%08x\n", $total, $hash
  continue
end

continue
