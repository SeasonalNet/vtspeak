set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $forced_null = 0
set $nested_alloc_calls = 0

break syncinfo_probe_alloc_start
commands
  silent
  disable 1
  break *0x10026360
  commands
    silent
    printf "SYNC_ALLOC_EXPORT_ENTER address=0x10026360\n"
    disable 2
    continue
  end
  break *0x1001d9c0
  commands
    silent
    if *(unsigned int *)$esp == 0x100263be
      set $nested_alloc_calls = $nested_alloc_calls + 1
      if $nested_alloc_calls == 1
        printf "SYNC_ALLOC_NESTED_HELPER_ENTER size=%d return=0x100263be\n", *(unsigned int *)($esp + 4)
      end
    end
    continue
  end
  break *0x1001d9d9
  commands
    silent
    if *(unsigned int *)($ebp + 4) == 0x100263be && $forced_null == 0
      set $eax = 0
      set $forced_null = 1
      printf "SYNC_ALLOC_FORCED_TRANSIENT_NULL callsite=0x100263be size=%d\n", $esi
    end
    continue
  end
  break *0x1001d9f1
  commands
    silent
    if *(unsigned int *)($ebp + 4) == 0x100263be && $forced_null == 1
      printf "SYNC_ALLOC_RETRY_RETURN callsite=0x100263be size=%d malloc_result_nonnull=%d\n", $esi, $eax != 0
      disable 5
    end
    continue
  end
  break syncinfo_probe_alloc_done
  commands
    silent
    printf "SYNC_ALLOC_INJECTED_OBJECT nested_alloc_calls=%d valid=%d rows=%d width=%d row_array=%d first_nested=%d last_nested=%d distinct_nested=%d\n", $nested_alloc_calls, *(int *)($ebp + 8), *(int *)($ebp + 12), *(int *)($ebp + 16), *(int *)($ebp + 20), *(int *)($ebp + 24), *(int *)($ebp + 28), *(int *)($ebp + 32)
    disable 6
    continue
  end
  continue
end

continue
