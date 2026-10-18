set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass
set $crosswalk_calls = 0
set $crosswalk_context = -1
set $crosswalk_state = 0
set $crosswalk_model = 0

# Save the selector context and model before the helper computes coverage.
hbreak *0x10023350
commands
  silent
  set $crosswalk_context = *(int *)($esp + 4)
  set $crosswalk_state = *(unsigned int *)($esp + 8)
  set $crosswalk_model = *(unsigned int *)($esp + 12)
  continue
end

# FUN_10023350 reaches this point after FUN_100230a0 has attached coverage
# metadata and applied its full-span filtering branch.
hbreak *0x10023814
commands
  silent
  set $crosswalk_calls = $crosswalk_calls + 1
  set $ctx = $crosswalk_context
  set $state = $crosswalk_state
  set $model = $crosswalk_model
  set $record = $state + 0xae894 + $ctx * 0xfc
  set $count = *(unsigned short *)($state + 0xae988 + $ctx * 0xfc)
  set $signature = *(unsigned int *)($model + 0x64)
  set $index = 0
  printf "APPLE_COVERAGE_CONTEXT call=%u context=%d positions=%u candidates=%u\n", $crosswalk_calls, $ctx, *(unsigned int *)($state + 0xec624), $count
  while $index < $count && $index < 60
    set $node = *(unsigned int *)($state + 0x477ac + $index * 4)
    set $unit = *(unsigned int *)($node + 8)
    set $unit_signature = $signature + $unit * 7
    printf "APPLE_COVERAGE_ROW context=%d rank=%u unit=%u right_matches=%d span=%d weighted_coverage=%d byte0=%#x byte6=%#x\n", $ctx, $index, $unit, *(short *)($node + 0xc), *(short *)($node + 0xe), *(short *)($node + 0x10), *(unsigned char *)$unit_signature, *(unsigned char *)($unit_signature + 6)
    set $index = $index + 1
  end
  continue
end

hbreak *0x1001b200
commands
  silent
  printf "APPLE_COVERAGE_SELECTED unit=%u\n", *(unsigned int *)($esp + 8)
  continue
end

continue
