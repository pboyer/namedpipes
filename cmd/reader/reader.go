package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	namedPipe := flag.Args()[0]

	pid := os.Getpid()
	fmt.Printf("ReaderStart @ %d\n", pid)

	// handle SIGTERM for clarity
	signalCh := make(chan os.Signal, syscall.SIGTERM)
	signal.Notify(signalCh)

	go func() {
		<-signalCh
		fmt.Printf("ReaderSIGTERM @ %d. Good bye!\n", pid)
		os.Exit(0)
	}()

	pipe, _ := os.OpenFile(namedPipe, os.O_RDONLY, 0600)

	// buffer to read input
	buf := make([]byte, 20)
	for {
		n, err := pipe.Read(buf)

		if err == io.EOF {
			fmt.Printf("ReaderEOF @ %d\n", pid)
			os.Exit(0)
		}

		if err != nil {
			fmt.Printf("ReaderError @ %d: %v\n", pid, err)
			os.Exit(1)
		}

		fmt.Printf("ReaderRead @ %d: %s\n", pid, string(buf[0:n]))
	}
}
