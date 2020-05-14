package udp

import (
	"context"
	"fmt"
	"net"
	"time"
)

// ResponseHandler - Response handler function
type ResponseHandler func(response Response) error

// Send : Send message over UDP
var Send = func(conf *Config, message []byte, handler ResponseHandler) error {
	ch := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.TimeoutInSeconds)*time.Second)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				ch <- fmt.Errorf("send-recovered: recovered from %s", r)
			}
		}()
		send(ch, conf, message, handler)
	}()

	select {
	case <-ctx.Done():
		close(ch)
		return ctx.Err()
	case err := <-ch:
		close(ch)
		return err
	}
}

func send(ch chan error, conf *Config, message []byte, handler ResponseHandler) {
	addr, err := net.ResolveUDPAddr(Network, conf.ServerAddress())
	if err != nil {
		ch <- err
		return
	}

	conn, err := net.DialUDP(Network, nil, addr)
	if err != nil {
		ch <- err
		return
	}

	defer conn.Close() //nolint:errcheck

	_, err = conn.Write(message)

	if err != nil {
		ch <- err
		return
	}

	if handler != nil {
		buffer := make([]byte, 1024)
		size, addr, err := conn.ReadFromUDP(buffer)

		address := ""
		if addr != nil {
			address = addr.String()
		}

		ch <- handler(Response{Size: size, Address: address, Body: buffer, Err: err})
	} else {
		ch <- nil
	}
}
