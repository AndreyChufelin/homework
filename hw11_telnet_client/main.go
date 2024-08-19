package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 5*time.Second, "timeout")
	flag.Parse()

	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	err := tc.Connect()
	if err != nil {
		log.Fatal("Connection error: ", err)
	}
	defer tc.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := tc.Send(); err != nil {
			log.Print("Send error: ", err)
			stop()
		}
	}()
	go func() {
		if err := tc.Receive(); err != nil {
			log.Print("Receive error: ", err)
		}
	}()

	<-ctx.Done()
}
