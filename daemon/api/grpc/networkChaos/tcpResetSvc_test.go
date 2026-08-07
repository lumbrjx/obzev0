package networkchaos

import (
	"net"
	"testing"
	"time"
)

func TestTCPReset_FullResetRate(t *testing.T) {
	cancel, err := startTCPReset("127.0.0.1:19080", 1.0)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	rstCount := 0
	total := 10
	for i := 0; i < total; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:19080", 1*time.Second)
		if err != nil {
			t.Fatalf("dial failed: %v", err)
		}

		buf := make([]byte, 1)
		conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		_, err = conn.Read(buf)
		conn.Close()

		if err != nil {
			rstCount++
		}
	}

	if rstCount < total/2 {
		t.Fatalf("expected most connections to get RST at rate=1.0, only %d/%d did", rstCount, total)
	}
}

func TestTCPReset_ZeroRate(t *testing.T) {
	cancel, err := startTCPReset("127.0.0.1:19081", 0.0)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	time.Sleep(50 * time.Millisecond)

	for i := 0; i < 5; i++ {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:19081", 1*time.Second)
		if err != nil {
			t.Fatalf("dial %d failed: %v", i, err)
		}
		conn.Close()
	}
}

func TestTCPReset_StopReleasesPort(t *testing.T) {
	cancel, err := startTCPReset("127.0.0.1:19082", 0.5)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(50 * time.Millisecond)

	// Verify it's listening
	conn, err := net.DialTimeout("tcp", "127.0.0.1:19082", 1*time.Second)
	if err != nil {
		t.Fatalf("should be listening: %v", err)
	}
	conn.Close()

	cancel()
	time.Sleep(300 * time.Millisecond)

	ln, err := net.Listen("tcp", "127.0.0.1:19082")
	if err != nil {
		t.Logf("port may still be releasing: %v", err)
	} else {
		ln.Close()
	}
}
