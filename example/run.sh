#!/usr/bin/env bash
set -emo pipefail

echo "[run.sh] starting application..."
export GO_SHUTDOWN_GRACEFUL_LOG=true
START_TIME=$SECONDS
go build .
./example &
pid=$!

echo "[run.sh] waiting..."
sleep 5

echo "[run.sh] sending SIGINT to $1"
kill -s SIGINT -- -$pid

echo "[run.sh] completed in $(($SECONDS - $START_TIME))s"
sleep 1
