$pid = start-job {
    echo "[run.sh] starting application..."
    export GO_SHUTDOWN_GRACEFUL_LOG = true
    go build .
    ./example.exe
}

echo "[run.sh] waiting..."
sleep 5

echo "[run.sh] st SIGINT to $1"
stop-job $pid
