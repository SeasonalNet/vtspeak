#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/vtspeak-speaker-metadata-input1.txt
backup_output=/tmp/vtspeak-speaker-metadata-output.wav
restore() {
  cp "$backup_input" "$work/input1.txt"
  cp "$backup_output" "$work/output.wav"
  cmp -s "$backup_input" "$work/input1.txt"
  cmp -s "$backup_output" "$work/output.wav"
  if [ -n "$xpid" ]; then kill "$xpid" 2>/dev/null || true; fi
}

cp "$work/input1.txt" "$backup_input"
cp "$work/output.wav" "$backup_output"
trap restore EXIT
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-speaker-metadata-contract.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-speaker-metadata-contract.gdb" > "$probe/speaker-metadata-contract-api.log" 2>&1

test "$(grep -c '^PATH_KEY slot=' "$probe/speaker-metadata-contract-api.log")" -eq 6
test "$(grep -c '^SPEAKERS_INFO_CONTRACT ' "$probe/speaker-metadata-contract-api.log")" -eq 6
grep -Fq 'PATH_KEY slot=0 pointer=0x100fe6e0 value=SOFTWARE\VW\VT\Kate\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'PATH_KEY slot=1 pointer=0x100fe6e0 value=SOFTWARE\VW\VT\Paul\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'PATH_KEY slot=2 pointer=0x100fe6e0 value=SOFTWARE\VW\VT\em001\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'PATH_KEY slot=3 pointer=0x100fe6e0 value=SOFTWARE\VW\VT\Julie\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'PATH_KEY slot=4 pointer=0x100fe6e0 value=SOFTWARE\VW\VT\James\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'PATH_KEY slot=5 pointer=0x100fe6e0 value=SOFTWARE\VW\VT\Ashley\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'PATH_KEY_BUFFER same_pointer=1 saved_first=SOFTWARE\VW\VT\Kate\M16 current_after_sweep=SOFTWARE\VW\VT\Ashley\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'PATH_KEY_EDGE input=INT_MIN value=SOFTWARE\VW\VT\Kate\M16' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=0 result=6 name_len=4 name_nul=0 name_next=0xa5 name_tail_bad=0' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=1 result=6 name_len=4 name_nul=0 name_next=0xa5 name_tail_bad=0' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=2 result=6 name_len=5 name_nul=0 name_next=0xa5 name_tail_bad=0' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=3 result=6 name_len=5 name_nul=0 name_next=0xa5 name_tail_bad=0' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=4 result=6 name_len=5 name_nul=0 name_next=0xa5 name_tail_bad=0' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=5 result=6 name_len=6 name_nul=0 name_next=0xa5 name_tail_bad=0' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'path_nul=0 path_next=0xa5 path_tail_bad=0' "$probe/speaker-metadata-contract-api.log"
test "$(grep -Ec '^SPEAKERS_INFO_CONTRACT .*path_nul=0 path_next=0xa5 path_tail_bad=0 ' "$probe/speaker-metadata-contract-api.log")" -eq 6
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=0 result=6 name_len=4 name_nul=0 name_next=0xa5 name_tail_bad=0 path_len=20 path_nul=0 path_next=0xa5 path_tail_bad=0 name=kate path=d:/eng/db/susan/pcm/' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=1 result=6 name_len=4 name_nul=0 name_next=0xa5 name_tail_bad=0 path_len=20 path_nul=0 path_next=0xa5 path_tail_bad=0 name=paul path=d:/eng/db/isaac/pcm/' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=2 result=6 name_len=5 name_nul=0 name_next=0xa5 name_tail_bad=0 path_len=20 path_nul=0 path_next=0xa5 path_tail_bad=0 name=em001 path=d:/eng/db/em001/pcm/' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=3 result=6 name_len=5 name_nul=0 name_next=0xa5 name_tail_bad=0 path_len=23 path_nul=0 path_next=0xa5 path_tail_bad=0 name=julie path=d:/eng/db/jennifer/pcm/' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=4 result=6 name_len=5 name_nul=0 name_next=0xa5 name_tail_bad=0 path_len=18 path_nul=0 path_next=0xa5 path_tail_bad=0 name=james path=d:/eng/db/lee/pcm/' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_CONTRACT slot=5 result=6 name_len=6 name_nul=0 name_next=0xa5 name_tail_bad=0 path_len=20 path_nul=0 path_next=0xa5 path_tail_bad=0 name=ashley path=d:/eng/db/casey/pcm/' "$probe/speaker-metadata-contract-api.log"
grep -Fq 'SPEAKERS_INFO_EDGE input=INT_MIN result=6 name=kate path=d:/eng/db/susan/pcm/' "$probe/speaker-metadata-contract-api.log"
