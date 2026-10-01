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
  set $lead = 0
  set $total = 0
  set $all_hash = 2166136261
  while $lead < 256
    set $trail = 0
    set $row_count = 0
    set $row_hash = 2166136261
    while $trail < 256
      set {unsigned char}($pair) = $lead
      set {unsigned char}($pair + 1) = $trail
      set $accepted = ((int (*)(unsigned char *))0x1001c900)($pair)
      set $row_count = $row_count + $accepted
      set $total = $total + $accepted
      set $row_hash = (($row_hash ^ $accepted) * 16777619) & 0xffffffff
      set $all_hash = (($all_hash ^ $accepted) * 16777619) & 0xffffffff
      set $trail = $trail + 1
    end
    printf "CLASSIFIER_ROW lead=%02x accepted=%d digest=%08x\n", $lead, $row_count, $row_hash
    set $lead = $lead + 1
  end
  printf "CLASSIFIER_TOTAL pairs=65536 accepted=%d digest=%08x\n", $total, $all_hash
  continue
end

continue
