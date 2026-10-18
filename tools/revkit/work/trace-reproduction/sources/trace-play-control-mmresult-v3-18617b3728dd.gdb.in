set pagination off
set confirm off
set debuginfod enabled off
handle SIGSEGV nostop noprint pass

set $phase = 0
set $pause_count = 0
set $restart_count = 0
set $pause_dispatched = 0
set $restart_dispatched = 0

break *0x10027010
commands
  silent
  set $pause_count = $pause_count + 1
  printf "PLAY_CONTROL pause_enter=%d handle=%#x\n", $pause_count, *(unsigned int *)0x100a7490
  continue
end

break *0x1002701a
commands
  silent
  set $pause_dispatched = 1
  printf "PLAY_CONTROL pause_winmm_call=%d handle=%#x\n", $pause_count, *(unsigned int *)0x100a7490
  continue
end

break *0x10027020
commands
  silent
  if $pause_dispatched == 1
    printf "PLAY_CONTROL pause_winmm_result=%d result=%u\n", $pause_count, $eax
    set $pause_dispatched = 0
  else
    printf "PLAY_CONTROL pause_no_handle_branch=%d eax=%#x\n", $pause_count, $eax
  end
  continue
end

break *0x10027030
commands
  silent
  set $restart_count = $restart_count + 1
  printf "PLAY_CONTROL restart_enter=%d handle=%#x\n", $restart_count, *(unsigned int *)0x100a7490
  continue
end

break *0x1002703a
commands
  silent
  set $restart_dispatched = 1
  printf "PLAY_CONTROL restart_winmm_call=%d handle=%#x\n", $restart_count, *(unsigned int *)0x100a7490
  continue
end

break *0x10027040
commands
  silent
  if $restart_dispatched == 1
    printf "PLAY_CONTROL restart_winmm_result=%d result=%u\n", $restart_count, $eax
    set $restart_dispatched = 0
  else
    printf "PLAY_CONTROL restart_no_handle_branch=%d eax=%#x\n", $restart_count, $eax
  end
  continue
end

break *0x4016d9
commands
  silent
  if $phase == 0
    set $phase = 1
    printf "PLAY_CONTROL null_handle_wrappers_returned handle=%#x\n", *(unsigned int *)0x100a7490
    set $esp = $esp - 4
    set {unsigned int}$esp = 0x4016d9
    set {int}($esp + 4) = 0
    set {int}($esp + 8) = 0
    set {int}($esp + 12) = $play_text
    set {int}($esp + 16) = 1
    set {int}($esp + 20) = -1
    set {int}($esp + 24) = -1
    set {int}($esp + 28) = -1
    set {int}($esp + 32) = -1
    set {int}($esp + 36) = -1
    set {int}($esp + 40) = 0
    set $eip = 0x10027050
    printf "PLAY_CONTROL redirected_to_play speaker=1\n"
    continue
  else
    if $phase == 1
      set $phase = 2
      printf "PLAY_CONTROL active_handle_play_return eax=%d handle=%#x state=%d\n", (short)$eax, *(unsigned int *)0x100a7490, *(int *)0x100a7498
      set $esp = $esp - 16
      set {unsigned int}$esp = 0x10027010
      set {unsigned int}($esp + 4) = 0x10027030
      set {unsigned int}($esp + 8) = 0x10027030
      set {unsigned int}($esp + 12) = 0x4016d9
      set $eip = 0x10027010
      continue
    else
      printf "PLAY_CONTROL repeated_calls_complete pause_calls=%d restart_calls=%d handle=%#x\n", $pause_count, $restart_count, *(unsigned int *)0x100a7490
      disable 7
      set $esp = $esp - 4
      set {unsigned int}$esp = 0x4016f9
      set $eip = 0x10027560
      continue
    end
  end
end

break *0x4016f9
commands
  silent
  printf "PLAY_CONTROL stopped handle=%#x state=%d\n", *(unsigned int *)0x100a7490, *(int *)0x100a7498
  disable 8
  continue
end

break *0x1001da50
commands
  silent
  set $play_text = *(unsigned int *)($esp + 8)
  set $phase = 0
  disable 9
  set $esp = $esp - 8
  set {unsigned int}$esp = 0x10027030
  set {unsigned int}($esp + 4) = 0x4016d9
  set $eip = 0x10027010
  printf "PLAY_CONTROL starting_null_handle=%#x text=%#x\n", *(unsigned int *)0x100a7490, $play_text
  continue
end

continue
