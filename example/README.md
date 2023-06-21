# Example

Execute the following command to run the example:

```bash
./run.sh
```

The example runs the following steps:

    - Build and run example the application
    - Sleep for 5 seconds
    - Send a SIGINT signal to running application

You should see the logs similar to the following:

```log
[run.sh] starting application...
[run.sh] waiting...
[example] 2023/06/21 12:00:03 started
[run.sh] sending SIGINT to 11926
[go-shutdown-graceful] 2023/06/21 12:00:08 received interrupt(2)! shutting down
[go-shutdown-graceful] 2023/06/21 12:00:08 waiting 1 for services/ routines to finish
[go-shutdown-graceful] 2023/06/21 12:00:08 all observers closed
[example] 2023/06/21 12:00:08 graceful shutdown complete
[run.sh] completed in 6s
```
