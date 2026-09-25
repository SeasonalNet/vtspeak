set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

break *0x10024680
commands
  silent
  set $word_index = *(short *)($esp + 4)
  set $engine = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "CANDIDATE_BUILD_ENTER word_index=%d engine=%#x\n", $word_index, $engine
  tbreak *$ret
  commands
    silent
    set $count = *(unsigned int *)($engine + 0xec624)
    printf "CANDIDATE_BUILD_RETURN word_index=%d count=%u eax=%#x\n", $word_index, $count, $eax
    set $i = 0
    while $i < $count && $i < 120
      printf "CANDIDATE_DESC index=%u bytes:", $i
      x/6bx ($engine + 0xec628 + $i * 6)
      set $i = $i + 1
    end
    continue
  end
  continue
end

break *0x10024060
commands
  silent
  set $candidate = *(unsigned int *)($esp + 4)
  set $engine = *(unsigned int *)($esp + 8)
  set $descriptor = $engine + ($candidate * 3 + 0x76314) * 2
  set $ret = *(unsigned int *)$esp
  printf "CANDIDATE_FILTER_A candidate=%u descriptor=", $candidate
  x/6bx $descriptor
  tbreak *$ret
  commands
    silent
    printf "CANDIDATE_FILTER_A_RETURN candidate=%u eax=%u\n", $candidate, $eax
    continue
  end
  continue
end

break *0x100242a0
commands
  silent
  set $candidate = *(unsigned int *)($esp + 4)
  set $engine = *(unsigned int *)($esp + 8)
  set $descriptor = $engine + ($candidate * 3 + 0x76314) * 2
  set $ret = *(unsigned int *)$esp
  printf "CANDIDATE_FILTER_B candidate=%u descriptor=", $candidate
  x/6bx $descriptor
  tbreak *$ret
  commands
    silent
    printf "CANDIDATE_FILTER_B_RETURN candidate=%u eax=%u\n", $candidate, $eax
    continue
  end
  continue
end

set $context_hits = 0
break *0x10018770
commands
  silent
  set $context_hits = $context_hits + 1
  if $context_hits <= 300
    printf "CANDIDATE_CONTEXT call=%u args=", $context_hits
    x/5wx $esp
    bt 2
  end
  continue
end

break *0x10023350
commands
  silent
  set $candidate = *(int *)($esp + 4)
  set $engine = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "CANDIDATE_SCORE_ENTER candidate=%d\n", $candidate
  tbreak *$ret
  commands
    silent
    printf "CANDIDATE_SCORE_RETURN candidate=%d count=%u metrics:", $candidate, *(short *)($engine + 0xae988 + $candidate * 0xfc)
    x/30fw ($engine + 0x8212c + $candidate * 0x78)
    continue
  end
  continue
end

break *0x10018c80
commands
  silent
  set $candidate = *(int *)($esp + 4)
  set $engine = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "CEP_SCORE_ENTER candidate=%d\n", $candidate
  tbreak *$ret
  commands
    silent
    printf "CEP_SCORE_RETURN candidate=%d metrics:", $candidate
    x/30fw ($engine + 0x8212c + ($candidate - 1) * 0x78)
    continue
  end
  continue
end

break *0x10024510
commands
  silent
  set $candidate = *(int *)($esp + 4)
  set $engine = *(unsigned int *)($esp + 8)
  set $ret = *(unsigned int *)$esp
  printf "CANDIDATE_TRIM_ENTER candidate=%d\n", $candidate
  tbreak *$ret
  commands
    silent
    printf "CANDIDATE_TRIM_RETURN candidate=%d metrics:", $candidate
    x/30fw ($engine + 0x8212c + $candidate * 0x78)
    continue
  end
  continue
end

break *0x10024900
commands
  silent
  set $first = *(int *)($esp + 4)
  set $last = *(int *)($esp + 8)
  set $engine = *(unsigned int *)($esp + 12)
  set $ret = *(unsigned int *)$esp
  printf "CANDIDATE_RANK_ENTER first=%d last=%d\n", $first, $last
  tbreak *$ret
  commands
    silent
    printf "CANDIDATE_RANK_RETURN first=%d last=%d selected=%d metrics:", $first, $last, *(short *)($engine + 0xae0c4 + $last * 2)
    x/30fw ($engine + 0x8212c + $last * 0x78)
    continue
  end
  continue
end

break *0x1001b200
commands
  silent
  printf "SELECTED_UNIT id=%u\n", *(unsigned int *)($esp + 8)
  continue
end

continue
