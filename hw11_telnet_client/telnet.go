package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		in:      in,
		out:     out,
		timeout: timeout,
		address: address,
	}
}

func (t *telnetClient) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("telnetClient.Connect: %w", err)
	}

	return nil
}

func (t telnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("telnetClient.Send: %w", err)
	}

	return nil
}

func (t telnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("telnetClient.Receive: %w", err)
	}

	return nil
}

var ErrNoConnection = errors.New("no connection")

func (t telnetClient) Close() error {
	if t.conn == nil {
		return fmt.Errorf("telnetClient.Close: %w", ErrNoConnection)
	}
	err := t.conn.Close()
	if err != nil {
		return fmt.Errorf("telnetClient.Close: %w", err)
	}

	return nil
}
