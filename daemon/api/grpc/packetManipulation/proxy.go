package packetmanipulation

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type MetricsData struct {
	DropedCount    int64
	CorruptedCount int64
}

var (
	Mtrx = make(chan MetricsData)
	Data = &MetricsData{}
)

type ProxyConfig struct {
	Server      string
	Client      string
	DropRate    float64
	CorruptRate float64
	Timeout     time.Duration
}

func Proxy(conf ProxyConfig) {
	listener, err := net.Listen("tcp", ":"+conf.Server)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port " + conf.Server)

	var wg sync.WaitGroup

	ctx, cancel := context.WithTimeout(
		context.Background(),
		conf.Timeout*time.Second,
	)
	defer cancel()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-ctx.Done():
					return
				default:
					fmt.Println("Error accepting connection:", err)
					continue
				}
			}

			clientConn, err := net.Dial("tcp", ":"+conf.Client)
			if err != nil {
				log.Printf("Error connecting to client: %v", err)
				conn.Close()
				continue
			}

			fmt.Println("Client connected:", conn.RemoteAddr().String())

			wg.Add(1)
			go DropPackets(
				conn,
				clientConn,
				&wg,
				conf.DropRate,
				conf.CorruptRate,
			)
		}
	}()

	<-ctx.Done()
	listener.Close()

	fmt.Println("Timeout reached. Shutting down server...")
	Mtrx <- *Data
	wg.Wait()

}
