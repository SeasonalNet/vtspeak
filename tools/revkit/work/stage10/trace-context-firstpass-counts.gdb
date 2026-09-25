set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# FUN_10024060: after FUN_10018770 candidate generation.
break *0x10024119
commands
  silent
  printf "FIRSTPASS_GENERATED index=%u flag=%u count=%d\n", $edi, *(unsigned char *)($ebp + 0xb), *(short *)($esi + 0xec620)
  continue
end

# Return from FUN_10023af0, which ranks/prunes the first candidate list.
break *0x10024164
commands
  silent
  printf "FIRSTPASS_RANKED index=%u output_count=%d\n", $edi, (short)$eax
  continue
end

# Return from FUN_10024010 after the other candidate list is deduplicated.
break *0x10024181
commands
  silent
  printf "FIRSTPASS_OTHER_LIST index=%u count=%d ranked_count=%d\n", $edi, (short)$eax, *(short *)($esi + 0xec620)
  continue
end

# FUN_10023060 result immediately before FUN_10024060's boolean threshold.
break *0x10024271
commands
  silent
  set $final_rec = *(unsigned int *)($ebp - 4)
  set $final_list = *(unsigned int *)($ebp - 0xc)
  set $final_count = *(short *)$final_rec
  set $final_table = *(unsigned int *)(*(unsigned int *)($ebp + 0x10) + 0x8c)
  printf "FIRSTPASS_FINAL index=%u flag=%u score_sum=%d candidate_count=%d distinct_count=%d\n", $edi, *(unsigned char *)($ebp + 0xb), (short)$eax, $final_count, *(short *)($final_rec + 2)
  if $edi == 1 || $edi == 3
    set $i = 0
    while $i < $final_count
      set $candidate_id = *(unsigned int *)($final_list + $i * 4)
      set $candidate_weight = *(unsigned short *)($final_table + $candidate_id * 2)
      printf "FIRSTPASS_WEIGHT index=%u slot=%u candidate=%u weight=%u\n", $edi, $i, $candidate_id, $candidate_weight
      set $i = $i + 1
    end
  end
  continue
end

continue
