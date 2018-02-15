package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

func main() {
	var writerpath string
	var readerpath string
	var numReaders uint64
	flag.StringVar(&readerpath, "reader", "./reader", "path to reader exec")
	flag.StringVar(&writerpath, "writer", "./writer", "path to writer exec")
	flag.Uint64Var(&numReaders, "numReaders", 2, "number of reader processes")
	flag.Parse()

	pid := os.Getpid()
	fmt.Printf("MasterStart @ %d\n", pid)

	// handle SIGTERM for clarity
	signalCh := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(signalCh)

	go func() {
		<-signalCh
		fmt.Printf("MasterSIGTERM @ %d. Good bye!\n", pid)
		os.Exit(0)
	}()

	tmpDir, _ := ioutil.TempDir("", "named-pipes")

	// make pipes
	pipes := make([]string, numReaders)

	for i := range pipes {
		namedPipe := filepath.Join(tmpDir, fmt.Sprintf("stdout%d", i))
		if err := syscall.Mkfifo(namedPipe, 0600); err != nil {
			fmt.Printf("MasterErr @ %d: %v\n", pid, err)
			os.Exit(1)
		}
		pipes[i] = namedPipe
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(pipes))

	// make readers
	for _, pipe := range pipes {
		cmd := exec.Command(readerpath, pipe)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		go func() {
			cmd.Run()
			wg.Done()
		}()
	}

	// make writer
	wcmd := exec.Command(writerpath, pipes...)
	wcmd.Stdout = os.Stdout
	wcmd.Stderr = os.Stderr

	wcmd.Run()
	wg.Wait()

	fmt.Printf("MasterDone @ %d\n", pid)
}
