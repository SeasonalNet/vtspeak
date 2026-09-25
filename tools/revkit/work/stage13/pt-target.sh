#!/usr/bin/env bash
set -euo pipefail

control_rel=${1:?usage: pt-target.sh CONTROL_RELATIVE_PATH MODE FIXTURE}
mode=${2:?usage: pt-target.sh CONTROL_RELATIVE_PATH MODE FIXTURE}
fixture=${3:?usage: pt-target.sh CONTROL_RELATIVE_PATH MODE FIXTURE}

if [[ ! $control_rel =~ ^stage13/control-pt-[A-Za-z0-9-]+$ ]]; then
  printf 'invalid Stage 13 control path: %s\n' "$control_rel" >&2
  exit 2
fi

control_dir="/work/$control_rel"
printf '%s\n' "$$" > "$control_dir/ready"

for _ in $(seq 1 600); do
  if [[ -f $control_dir/start ]]; then
    rm -f -- "$control_dir/start"
    mkdir -p -- "$control_dir/captures" "$control_dir/scratch"
    VTSPEAK_CAPTURE_DIR="$control_dir/captures" \
      VTSPEAK_SCRATCH_DIR="$control_dir/scratch" \
      /bin/bash /work/stage12/run-effects.sh "$mode" "$fixture"
    runner_log="$control_dir/captures/$fixture-effect-$mode.log"
    runner_wav="$control_dir/captures/$fixture-effect-$mode.wav"
    if ! grep -Fq '[Inferior 1 (Remote target) exited normally]' "$runner_log"; then
      printf 'Stage 12 GDB inferior did not exit normally; see %s\n' "$runner_log" >&2
      exit 1
    fi
    case "$mode" in
      acoustic-control)
        if ! grep -q '^ACOUSTIC_DESC ' "$runner_log"; then
          printf 'acoustic-control produced no descriptor records; see %s\n' "$runner_log" >&2
          exit 1
        fi
        ;;
      acoustic-83-zero)
        if ! grep -q '^TPP_NUMERIC_ABLATE83 code=83$' "$runner_log" ||
          ! grep -q '^ACOUSTIC_DESC ' "$runner_log"; then
          printf 'G83-zero probe is missing its intervention or descriptor records; see %s\n' "$runner_log" >&2
          exit 1
        fi
        ;;
    esac
    if [[ $mode == acoustic-control || $mode == acoustic-83-zero ]] &&
      cmp -s /work/stage5/output.wav "$runner_wav"; then
      printf '%s output WAV is unchanged from the Stage 5 seed; refusing this capture\n' "$mode" >&2
      exit 1
    fi
    cp "$runner_log" "$control_dir/runner.log"
    cp "$runner_wav" "$control_dir/output.wav"
    exit 0
  fi
  sleep 0.1
done

printf 'timed out waiting for host Intel PT collector\n' >&2
exit 124
