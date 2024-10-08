package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrNotConnected = fmt.Errorf("not connected")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type Telnet struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (t *Telnet) Connect() error {
	var err error
	dialer := &net.Dialer{}
	dialer.Timeout = t.timeout
	t.conn, err = dialer.Dial("tcp", t.address)
	if err != nil {
		return err
	}

	return nil
}

func (t *Telnet) Close() error {
	if t.conn == nil {
		return ErrNotConnected
	}

	err := t.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (t *Telnet) Send() error {
	if t.conn == nil {
		return ErrNotConnected
	}
	return processIO(t.in, t.conn)
}

func (t *Telnet) Receive() error {
	if t.conn == nil {
		return ErrNotConnected
	}
	return processIO(t.conn, t.out)
}

func processIO(src io.Reader, dst io.Writer) error {
	buf := bufio.NewScanner(src)
	for buf.Scan() {
		_, err := dst.Write([]byte(fmt.Sprintf("%s\n", buf.Text())))
		if err != nil {
			return err
		}
	}
	return nil
}
