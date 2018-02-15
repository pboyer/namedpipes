package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	nps := flag.Args()

	pid := os.Getpid()
	fmt.Printf("WriterStart @ %d\n", pid)

	// handle SIGTERM for clarity
	signalCh := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(signalCh)

	go func() {
		<-signalCh
		fmt.Printf("WriterSIGTERM @ %d. Good bye!\n", pid)
		os.Exit(0)
	}()

	// open pipes
	pipes := make([]*os.File, len(nps))
	for i, np := range nps {
		pipe, err := os.OpenFile(np, os.O_WRONLY, 0600)
		if err != nil {
			fmt.Printf("WriterErr @ %d: %v\n", pid, err)
			os.Exit(1)
		}
		pipes[i] = pipe
	}

	// write in parallel
	wg := &sync.WaitGroup{}
	wg.Add(len(pipes))
	for _, fd := range pipes {
		go func(fd *os.File) {
			// writes 10 integers to fd
			for i := 0; i < 10; i++ {
				_, err := fmt.Fprintf(fd, "%d ", i)
				if err != nil {
					fmt.Printf("WriterErr @ %d: %v\n", pid, err)
					os.Exit(1)
				}
				time.Sleep(100 * time.Millisecond) // deliberately wait 100ms so you can see what's happening
			}
			wg.Done()
			fd.Close()
		}(fd)
	}

	wg.Wait()

	fmt.Printf("WriterDone @ %d\n", pid)
}
