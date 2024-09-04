package helper

import (
	"log"
	"net"
	"time"
)

func ReqSimulator(
	targetAddr string,
	duration time.Duration,
) error {
	endTime := time.Now().Add(duration)

	for time.Now().Before(endTime) {
		conn, err := net.Dial("tcp", ":"+targetAddr)
		if err != nil {
			log.Printf("Failed to connect to target: %v", err)
			return err
		}

		request := "GET / HTTP/1.1\r\n" +
			"Host: " + targetAddr + "\r\n" +
			"Connection: close\r\n" +
			"\r\n"

		_, err = conn.Write([]byte(request))
		if err != nil {
			log.Printf("Failed to write to target: %v", err)
			conn.Close()
			return err
		}

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Failed to read from target: %v", err)
			conn.Close()
			return err
		}

		log.Printf("Received response: %s", string(buf[:n]))

		conn.Close()

		time.Sleep(1 * time.Second)
	}

	return nil
}
