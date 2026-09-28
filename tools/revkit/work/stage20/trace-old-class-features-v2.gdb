set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /work/corpus-parity/stage20/old-class-features-v6-gdb.log
set logging overwrite on
set logging enabled on

break *0x408187
commands
  silent
  printf "OLD_FEATURE_TRACE_V2_READY\n"

break *0x1001da80
commands
  silent
  set $lookup_caller = *(unsigned int *)$esp
  if $lookup_caller >= 0x1001dfb0 && $lookup_caller < 0x1001e110
    set $lookup_key = *(unsigned int *)($esp + 4)
    set $lookup_slot = *(unsigned int *)($esp + 8)
    set $lookup_state = *(unsigned int *)($esp + 12)
    set $lookup_model = *(unsigned int *)($esp + 16)
    set $lookup_mode = *(unsigned short *)($esp + 20)
    set $lookup_return = $lookup_caller
    printf "OLD_CLASS_LOOKUP_ENTER slot=%u mode=%u key:", $lookup_slot, $lookup_mode
    x/5ub $lookup_key
    tbreak *$lookup_return
    commands
      silent
      set $lookup_count = *(unsigned short *)($lookup_state + 0xf8e64)
      set $lookup_features = *(unsigned int *)($lookup_model + 0x98)
      printf "OLD_CLASS_LOOKUP_RETURN slot=%u mode=%u appended=%u total=%u ids:", $lookup_slot, $lookup_mode, $eax, $lookup_count
      set $lookup_i = 0
      while $lookup_i < $lookup_count && $lookup_i < 32
        set $lookup_id = *(unsigned int *)($lookup_state + 0xf89b4 + $lookup_i * 4)
        printf " %u", $lookup_id
        printf "\nOLD_CLASS_FEATURE_BASE model=%#x base=%#x id=%u bytes:", $lookup_model, $lookup_features, $lookup_id
        x/7ub ($lookup_features + $lookup_id * 7)
        set $lookup_i = $lookup_i + 1
      end
      printf "\n"
      continue
    end
  end
  continue
end

continue
end

continue
