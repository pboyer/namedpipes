There are three executables found in the `cmd` directory:

* master
* reader
* writer

Use `./run.sh` to run the demo.

### How it works

`master` creates one `writer` child process and n `reader` child processes. For each child, `stdout` and `stderr` are piped to `master`'s stdout for visibility.

In addition to the child processes, `master` also creates n named pipes, one for each `reader`, using the `mkfifo` system call.

For more info on `mkfifo`: https://linux.die.net/man/3/mkfifo

`master` passes the name of each named pipe to the respective `reader` process. It passes the name of all named pipes to the `writer` process.

`writer` proceeds to write to each of the named pipes in parallel. It simply writes increasing integers to 10. It could write any byte stream to the named pipe.

Each `reader` process simply reads from it's respective named pipe and writes the output to `stdout`.
