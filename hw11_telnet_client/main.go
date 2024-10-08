package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var duration time.Duration

func init() {
	flag.DurationVar(&duration, "timeout", 10*time.Second, "timeout")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("host and port are required")
		return
	}

	address := net.JoinHostPort(args[0], args[1])
	client := ConnectClient(address, duration, io.NopCloser(bufio.NewReader(os.Stdin)), os.Stdout)
	defer client.Close()

	os.Stderr.WriteString(fmt.Sprintf("...Connected to %s\n", address))
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	go func() {
		err := client.Send()
		if err != nil {
			cancel()
		}
	}()

	go func() {
		err := client.Receive()
		if err != nil {
			cancel()
		}
	}()

	<-ctx.Done()
	os.Stderr.WriteString("...Connection was closed by peer\n")
}

func ConnectClient(address string, duration time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	client := NewTelnetClient(address, duration, in, out)
	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	return client
}
