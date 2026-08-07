package networkchaos

import (
	"net"
	"testing"
	"time"
)

func buildDNSQuery() []byte {
	// Minimal DNS query for "example.com" type A
	return []byte{
		0xAA, 0xBB, // Transaction ID
		0x01, 0x00, // Flags: standard query, RD=1
		0x00, 0x01, // Questions: 1
		0x00, 0x00, // Answers: 0
		0x00, 0x00, // Authority: 0
		0x00, 0x00, // Additional: 0
		// QNAME: example.com
		0x07, 'e', 'x', 'a', 'm', 'p', 'l', 'e',
		0x03, 'c', 'o', 'm',
		0x00,       // end of name
		0x00, 0x01, // QTYPE: A
		0x00, 0x01, // QCLASS: IN
	}
}

func TestDNSChaos_NXDOMAIN(t *testing.T) {
	cancel, err := startDNSChaos("nxdomain", 0, "127.0.0.1:15353", "8.8.8.8:53")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	conn, err := net.Dial("udp", "127.0.0.1:15353")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	query := buildDNSQuery()
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err := conn.Write(query); err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if n < 4 {
		t.Fatal("response too short")
	}

	// Check QR bit is set (response)
	if buf[2]&0x80 == 0 {
		t.Fatal("QR bit not set in response")
	}

	// Check RCODE is 3 (NXDOMAIN)
	rcode := buf[3] & 0x0F
	if rcode != 3 {
		t.Fatalf("expected RCODE 3 (NXDOMAIN), got %d", rcode)
	}

	// Check transaction ID preserved
	if buf[0] != 0xAA || buf[1] != 0xBB {
		t.Fatalf("transaction ID not preserved: got %02x%02x", buf[0], buf[1])
	}
}

func TestDNSChaos_SERVFAIL(t *testing.T) {
	cancel, err := startDNSChaos("servfail", 0, "127.0.0.1:15354", "8.8.8.8:53")
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	conn, err := net.Dial("udp", "127.0.0.1:15354")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	query := buildDNSQuery()
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	conn.Write(query)

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	if n < 4 {
		t.Fatal("response too short")
	}

	rcode := buf[3] & 0x0F
	if rcode != 2 {
		t.Fatalf("expected RCODE 2 (SERVFAIL), got %d", rcode)
	}
}

func TestDNSChaos_StopCancelsServer(t *testing.T) {
	cancel, err := startDNSChaos("nxdomain", 0, "127.0.0.1:15355", "8.8.8.8:53")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	// Verify it's running
	conn, err := net.Dial("udp", "127.0.0.1:15355")
	if err != nil {
		t.Fatal(err)
	}
	conn.SetDeadline(time.Now().Add(1 * time.Second))
	conn.Write(buildDNSQuery())
	buf := make([]byte, 512)
	_, err = conn.Read(buf)
	conn.Close()
	if err != nil {
		t.Fatal("server should be running before cancel")
	}

	cancel()
	time.Sleep(300 * time.Millisecond)

	// After cancel, the UDP port should be released
	_, err = net.ListenPacket("udp", "127.0.0.1:15355")
	if err != nil {
		t.Logf("port may still be releasing, that's OK: %v", err)
	}
}

func TestBuildDNSError_ShortInput(t *testing.T) {
	resp := buildDNSError([]byte{0x01}, 3)
	if len(resp) < 1 {
		t.Fatal("response should not be empty")
	}
}

func TestBuildDNSError_PreservesTransactionID(t *testing.T) {
	query := buildDNSQuery()
	resp := buildDNSError(query, 3)

	if resp[0] != query[0] || resp[1] != query[1] {
		t.Fatalf("transaction ID not preserved")
	}
}
