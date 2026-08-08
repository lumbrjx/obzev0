package networkchaos

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

var activeDNSCancel atomic.Value

// startDNSChaos starts a UDP DNS proxy that intercepts queries according to mode:
//
//	"nxdomain"  — respond NXDOMAIN to every query
//	"servfail"  — respond SERVFAIL to every query
//	"delay"     — forward to upstream after delayMs sleep
//
// Returns a cancel function that stops the server.
func startDNSChaos(mode string, delayMs int32, listenAddr, upstream string) (context.CancelFunc, error) {
	pc, err := net.ListenPacket("udp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("DNS listen on %s: %w", listenAddr, err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer pc.Close()
		buf := make([]byte, 512)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			pc.SetDeadline(time.Now().Add(200 * time.Millisecond))
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				continue
			}

			query := make([]byte, n)
			copy(query, buf[:n])

			go handleDNSQuery(pc, addr, query, mode, delayMs, upstream)
		}
	}()

	return cancel, nil
}

func handleDNSQuery(pc net.PacketConn, addr net.Addr, query []byte, mode string, delayMs int32, upstream string) {
	switch mode {
	case "nxdomain":
		resp := buildDNSError(query, 3) // NXDOMAIN rcode=3
		pc.WriteTo(resp, addr)

	case "servfail":
		resp := buildDNSError(query, 2) // SERVFAIL rcode=2
		pc.WriteTo(resp, addr)

	default: // "delay" or anything else — forward after optional sleep
		if delayMs > 0 {
			time.Sleep(time.Duration(delayMs) * time.Millisecond)
		}
		resp, err := forwardDNS(query, upstream)
		if err != nil {
			resp = buildDNSError(query, 2)
		}
		pc.WriteTo(resp, addr)
	}
}

// buildDNSError constructs a minimal DNS error response from the query bytes.
// It copies the transaction ID and sets the QR bit + rcode.
func buildDNSError(query []byte, rcode byte) []byte {
	if len(query) < 2 {
		return []byte{0, 0, 0x80 | rcode, 0}
	}
	resp := make([]byte, len(query))
	copy(resp, query)
	// byte 2: QR=1, opcode=0, AA=0, TC=0, RD=copy
	resp[2] = 0x80 | (query[2] & 0x01) // preserve RD bit
	// byte 3: RA=0, Z=0, RCODE
	resp[3] = rcode & 0x0F
	return resp
}

func forwardDNS(query []byte, upstream string) ([]byte, error) {
	conn, err := net.DialTimeout("udp", upstream, 3*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err := conn.Write(query); err != nil {
		return nil, err
	}
	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	return buf[:n], err
}
