#!/usr/bin/env bash
set -emo pipefail

function stop_main() {
  echo "[run.sh] waiting..."
  sleep 5
  echo "[run.sh] sending SIGINT to $1"
  kill -s SIGINT -- -$1
}

export GO_SHUTDOWN_GRACEFUL_LOG=true
START_TIME=$SECONDS
echo "[run.sh] starting application..."
go build .
./example &
stop_main $!
wait
echo "[run.sh] completed in $(($SECONDS - $START_TIME))s"
sleep 1
