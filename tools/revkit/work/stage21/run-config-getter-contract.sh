#!/bin/bash
set -euo pipefail

work=/work/stage5
probe=/work/stage21
xpid=
backup_input=/tmp/vtspeak-config-getter-input1.txt
backup_output=/tmp/vtspeak-config-getter-output.wav
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
Xvfb :99 -screen 0 1280x1024x24 -nolisten tcp > "$probe/xvfb-config-getter-contract.log" 2>&1 &
xpid=$!
sleep 1
export DISPLAY=:99
WINEPREFIX=/work/stage2-copy/wineprefix WINEDEBUG=-all timeout 180s \
  winedbg --gdb /samples/voicetext_paul.exe \
  < "$probe/trace-config-getter-contract.gdb" > "$probe/config-getter-contract-api.log" 2>&1

test "$(grep -c '^PSV_NULLMASK ' "$probe/config-getter-contract-api.log")" -eq 16
test "$(grep -c '^COMMA_NULLMASK ' "$probe/config-getter-contract-api.log")" -eq 2
test "$(grep -c '^PSV_SLOT ' "$probe/config-getter-contract-api.log")" -eq 10
test "$(grep -c '^COMMA_SLOT ' "$probe/config-getter-contract-api.log")" -eq 10
grep -Fq 'PSV_NULLMASK mask=0 ret=1 fields=11111111,22222222,33333333,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=1 ret=1 fields=00000064,22222222,33333333,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=2 ret=1 fields=11111111,00000064,33333333,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=3 ret=1 fields=00000064,00000064,33333333,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=4 ret=1 fields=11111111,22222222,000000c8,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=5 ret=1 fields=00000064,22222222,000000c8,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=6 ret=1 fields=11111111,00000064,000000c8,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=7 ret=1 fields=00000064,00000064,000000c8,44444444' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=8 ret=1 fields=11111111,22222222,33333333,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=9 ret=1 fields=00000064,22222222,33333333,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=10 ret=1 fields=11111111,00000064,33333333,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=11 ret=1 fields=00000064,00000064,33333333,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=12 ret=1 fields=11111111,22222222,000000c8,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=13 ret=1 fields=00000064,22222222,000000c8,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=14 ret=1 fields=11111111,00000064,000000c8,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_NULLMASK mask=15 ret=1 fields=00000064,00000064,000000c8,0000039d' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_NULLMASK mask=0 ret=1 value=55667788' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_NULLMASK mask=1 ret=1 value=000000c8' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=-2147483648 ret=1 fields=100,100,200,925' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=-1 ret=1 fields=100,100,200,925' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=0 ret=-1 fields=286331153,572662306,858993459,1145324612' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=1 ret=1 fields=100,100,200,925' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=2 ret=-1 fields=286331153,572662306,858993459,1145324612' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=3 ret=-1 fields=286331153,572662306,858993459,1145324612' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=4 ret=-1 fields=286331153,572662306,858993459,1145324612' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=5 ret=-1 fields=286331153,572662306,858993459,1145324612' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=6 ret=1 fields=100,100,200,925' "$probe/config-getter-contract-api.log"
grep -Fq 'PSV_SLOT selector=2147483647 ret=1 fields=100,100,200,925' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=-2147483648 ret=1 value=200' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=-1 ret=1 value=200' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=0 ret=-1 value=1432778632' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=2 ret=-1 value=1432778632' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=3 ret=-1 value=1432778632' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=4 ret=-1 value=1432778632' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=5 ret=-1 value=1432778632' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=1 ret=1 value=200' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=6 ret=1 value=200' "$probe/config-getter-contract-api.log"
grep -Fq 'COMMA_SLOT selector=2147483647 ret=1 value=200' "$probe/config-getter-contract-api.log"
