set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

# Returns from the base and first synthesized-signature FUN_10023f90 calls.
break *0x1001896a
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $base_state = $edi
    printf "BYTE_BASE_LOOKUP_RETURN value=%d candidate_count=%d signature=", (short)$eax, *(short *)($base_state + 0xec620)
    x/7ub $ebp-0x54
  end
  continue
end

break *0x10018a45
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $base_state = $edi
    printf "BYTE_COMBINED_LOOKUP_RETURN value=%d candidate_count=%d signature=", (short)$eax, *(short *)($base_state + 0xec620)
    x/7ub $ebp-0x54
    if (short)$eax > 0
      set $emitted_slot = 0
      while $emitted_slot < *(short *)($base_state + 0xec620)
        set $emitted_id = *(unsigned int *)($base_state + 0xec170 + $emitted_slot * 4)
        printf "BYTE_LOOKUP_EMITTED slot=%u id=%u\n", $emitted_slot, $emitted_id
        set $emitted_slot = $emitted_slot + 1
      end
    end
  end
  continue
end

break *0x10018ae4
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_input = *(unsigned int *)($ebp + 8)
    set $branch_state = $edi
    printf "BYTE_BRANCH_PLUS2_ZERO input2=%u mapped=%u\n", *(unsigned char *)($branch_input + 2), *(unsigned char *)(0x1007c2c4 + *(unsigned char *)($branch_input + 2))
  end
  continue
end

break *0x10018af1
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_input = *(unsigned int *)($ebp + 8)
    printf "BYTE_BRANCH_PLUS1_TEST input1=%u mapped=%u\n", *(unsigned char *)($branch_input + 1), $al
  end
  continue
end

break *0x10018b03
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_input = *(unsigned int *)($ebp + 8)
    set $branch_state = $edi
    set $branch_raw3 = *(unsigned char *)($branch_input + 3)
    set $branch_mapped3 = *(unsigned char *)(0x1007c2c4 + $branch_raw3)
    printf "BYTE_BRANCH_PLUS3_TEST input3=%u mapped=%u candidate_count=%d\n", $branch_raw3, $branch_mapped3, *(short *)($branch_state + 0xec620)
  end
  continue
end

break *0x10018b05
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_input = *(unsigned int *)($ebp + 8)
    printf "BYTE_BRANCH_PLUS3_NONZERO_PATH input3=%u\n", *(unsigned char *)($branch_input + 3)
  end
  continue
end

break *0x10018b19
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_state = $edi
    printf "BYTE_BRANCH_PLUS3_CALL1_RETURN value=%d candidate_count=%d\n", (short)$eax, *(short *)($branch_state + 0xec620)
  end
  continue
end

break *0x10018b3d
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_state = $edi
    printf "BYTE_BRANCH_PLUS3_CALL2_RETURN value=%d candidate_count=%d\n", (short)$eax, *(short *)($branch_state + 0xec620)
  end
  continue
end

break *0x10018b70
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_state = $edi
    printf "BYTE_BRANCH_PLUS3_CALL3_RETURN value=%d candidate_count=%d\n", (short)$eax, *(short *)($branch_state + 0xec620)
  end
  continue
end

break *0x10018b77
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_input = *(unsigned int *)($ebp + 8)
    printf "BYTE_BRANCH_PLUS3_ZERO_PATH input3=%u plus1_mapped=%u\n", *(unsigned char *)($branch_input + 3), $al
  end
  continue
end

break *0x10018b8f
commands
  silent
  if *(unsigned int *)($ebp + 4) == 0x10024119 && *(unsigned int *)($ebp + 0xc) == 1
    set $branch_state = $edi
    printf "BYTE_BRANCH_PLUS1_CALL_RETURN value=%d candidate_count=%d\n", (short)$eax, *(short *)($branch_state + 0xec620)
  end
  continue
end

continue
