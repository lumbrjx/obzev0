package networkchaos

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync/atomic"
	"time"
)

var activeTCPResetCancel atomic.Value

// startTCPReset listens on listenAddr and forces TCP RST on a fraction of incoming connections.
// resetRate is a value in [0.0, 1.0]; connections above the threshold are proxied normally
// (connect succeeds but immediately closes, since we don't know the upstream here).
func startTCPReset(listenAddr string, resetRate float32) (context.CancelFunc, error) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("TCP reset listen on %s: %w", listenAddr, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer ln.Close()
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			ln.(*net.TCPListener).SetDeadline(time.Now().Add(200 * time.Millisecond))
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			go resetOrClose(conn, resetRate)
		}
	}()

	return cancel, nil
}

func resetOrClose(conn net.Conn, resetRate float32) {
	tc, ok := conn.(*net.TCPConn)
	if !ok {
		conn.Close()
		return
	}
	if rand.Float32() < resetRate {
		// SO_LINGER=0 forces RST on close instead of FIN.
		tc.SetLinger(0)
	}
	tc.Close()
}
