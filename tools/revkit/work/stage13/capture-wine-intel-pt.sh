#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat >&2 <<'EOF'
Usage:
  capture-wine-intel-pt.sh OUTPUT.perf.data MODE FIXTURE

Example:
  capture-wine-intel-pt.sh /tmp/vtspeak-g83-control.perf.data \
    acoustic-control numeric-context-g83-can-anybody-help

MODE and FIXTURE are passed to the existing Stage 12 run-effects.sh runner.
The script records a single fixture and requires sudo access for host perf.
EOF
}

if [[ $# -ne 3 ]]; then
  usage
  exit 2
fi

output=$1
mode=$2
fixture=$3

if [[ ! $fixture =~ ^[A-Za-z0-9][A-Za-z0-9._-]*$ ]]; then
  printf 'invalid fixture name: %s\n' "$fixture" >&2
  exit 2
fi

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
repo_root=$(cd -- "$script_dir/../../../.." && pwd)
work_root="$repo_root/tools/revkit/work"
stage_dir="$work_root/stage13"
compose_file="$work_root/stage8/compose.yaml"

if [[ ! -f "$stage_dir/pt-target.sh" || ! -f "$compose_file" ]]; then
  printf 'required Stage 13 probe files or Stage 8 Compose file are missing\n' >&2
  exit 1
fi

if [[ ! -f "$stage_dir/../stage12/inputs/$fixture.txt" ]]; then
  printf 'Stage 12 fixture does not exist: %s\n' "$fixture" >&2
  exit 2
fi

if [[ $output != /* ]]; then
  output="$repo_root/$output"
fi
output_dir=$(dirname -- "$output")
mkdir -p -- "$output_dir"
if [[ -e $output ]]; then
  printf 'refusing to overwrite existing trace: %s\n' "$output" >&2
  exit 2
fi
if [[ -e $output.wine.log || -e $output.wav ]]; then
  printf 'refusing to overwrite an existing probe sidecar for: %s\n' "$output" >&2
  exit 2
fi

command -v docker >/dev/null || { printf 'docker is required\n' >&2; exit 1; }
command -v sudo >/dev/null || { printf 'sudo is required\n' >&2; exit 1; }
command -v perf >/dev/null || { printf 'perf is required on the host\n' >&2; exit 1; }

sudo -v

session_id="pt-$(date +%s)-$$"
control_dir="$stage_dir/control-$session_id"
control_rel="stage13/control-$session_id"
mkdir -m 700 -- "$control_dir"
container_name="vtspeak-$session_id"
container_id=
perf_pid=
perf_log="$control_dir/perf.log"

cleanup() {
  local status=$?
  trap - EXIT INT TERM HUP

  if [[ -n $perf_pid ]] && kill -0 "$perf_pid" 2>/dev/null; then
    kill -INT "$perf_pid" 2>/dev/null || true
    wait "$perf_pid" 2>/dev/null || true
  fi

  if [[ -n $container_id ]]; then
    docker rm --force "$container_id" >/dev/null 2>&1 || true
  fi

  rm -rf -- "$control_dir/stage5" "$control_dir/captures" "$control_dir/scratch"
  rm -f -- "$control_dir/ready" "$control_dir/start" \
    "$control_dir/perf.log" "$control_dir/runner.log" "$control_dir/output.wav"
  rmdir -- "$control_dir" 2>/dev/null || true
  exit "$status"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
trap 'exit 129' HUP

cp -a -- "$work_root/stage5" "$control_dir/stage5"
container_id=$(docker compose -f "$compose_file" run --detach \
  --volume "$control_dir/stage5:/work/stage5" \
  --name "$container_name" runtime \
  /bin/bash /work/stage13/pt-target.sh "$control_rel" "$mode" "$fixture")

if [[ -z $container_id ]]; then
  printf 'docker compose did not return a container ID\n' >&2
  exit 1
fi

ready=0
for _ in $(seq 1 300); do
  if [[ -f $control_dir/ready ]]; then
    ready=1
    break
  fi
  state=$(docker inspect --format '{{.State.Status}}' "$container_id" 2>/dev/null || true)
  if [[ $state != running ]]; then
    printf 'runtime container stopped before reaching the PT gate (state: %s)\n' "${state:-unknown}" >&2
    docker logs "$container_id" >&2 || true
    exit 1
  fi
  sleep 0.1
done

if [[ $ready -ne 1 ]]; then
  printf 'timed out waiting for runtime container PT gate\n' >&2
  exit 1
fi

host_pid=$(docker inspect --format '{{.State.Pid}}' "$container_id")
if [[ ! $host_pid =~ ^[1-9][0-9]*$ ]]; then
  printf 'docker returned an invalid host PID: %s\n' "$host_pid" >&2
  exit 1
fi

printf 'Container: %s\nHost PID to trace: %s\n' "$container_id" "$host_pid"
printf 'Starting host Intel PT capture; the existing Stage 12 runner starts after perf attaches.\n'

sudo perf record -vv -o "$output" -e intel_pt//u -p "$host_pid" \
  >"$perf_log" 2>&1 &
perf_pid=$!

attached=0
for _ in $(seq 1 100); do
  if ! kill -0 "$perf_pid" 2>/dev/null; then
    break
  fi
  if grep -Eq 'Permission denied|Operation not permitted|not supported|Error:' "$perf_log"; then
    break
  fi
  if grep -q 'sys_perf_event_open:' "$perf_log"; then
    sleep 0.2
    if kill -0 "$perf_pid" 2>/dev/null &&
      ! grep -Eq 'Permission denied|Operation not permitted|not supported|Error:' "$perf_log"; then
      attached=1
    fi
    break
  fi
  sleep 0.1
done

if [[ $attached -ne 1 ]]; then
  printf 'perf did not confirm opening the Intel PT event; refusing to start the probe\n' >&2
  cat "$perf_log" >&2
  exit 1
fi

touch -- "$control_dir/start"
printf 'Probe started: %s / %s\n' "$mode" "$fixture"

container_status=$(docker wait "$container_id")
if [[ $container_status != 0 ]]; then
  printf 'VoiceText runtime exited with status %s\n' "$container_status" >&2
  cat "$perf_log" >&2
  exit 1
fi

if wait "$perf_pid"; then
  perf_pid=
else
  perf_status=$?
  perf_pid=
  printf 'perf exited with status %s\n' "$perf_status" >&2
  cat "$perf_log" >&2
  exit "$perf_status"
fi

if [[ ! -s $output ]]; then
  printf 'perf reported success but produced no trace data\n' >&2
  cat "$perf_log" >&2
  exit 1
fi

sudo chown "$(id -u):$(id -g)" "$output"
mv -- "$control_dir/runner.log" "$output.wine.log"
mv -- "$control_dir/output.wav" "$output.wav"
cat "$perf_log"
printf 'Trace saved: %s\n' "$output"
printf 'VoiceText GDB log: %s\n' "$output.wine.log"
printf 'VoiceText WAV: %s\n' "$output.wav"
printf 'Decode with: sudo perf script -i %q > %q.script\n' "$output" "$output"
