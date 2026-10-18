set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set logging file /probe/stage20/hi-2006-row-gain-gdb.log
set logging overwrite on
set logging enabled on
set $hi_gain_calls = 0

hbreak *0x10023a10
commands
  silent
  set $hi_gain_calls = $hi_gain_calls + 1
  set $state = *(unsigned int *)($esp + 8)
  set $slot = *(short *)($esp + 12)
  if $slot >= 0 && $slot < 32
    printf "LEGACY_HI_ROW_GAIN call=%u slot=%d pcm_count=%d volume_scale=%d pitch_scale=%d duration_scale=%d\n", $hi_gain_calls, $slot, *(int *)($state + 11000 + $slot * 4), *(int *)($state + 0xce44 + $slot * 4), *(int *)($state + 0xabe4 + $slot * 4), *(int *)($state + 0xbd14 + $slot * 4)
  else
    printf "LEGACY_HI_ROW_GAIN call=%u slot=%d out_of_range=1\n", $hi_gain_calls, $slot
  end
  continue
end

hbreak *0x408187
commands
  silent
  printf "LEGACY_HI_ROW_GAIN_TRACE_READY\n"
  continue
end

continue
