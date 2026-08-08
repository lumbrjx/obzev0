package latency

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"obzev0/common/definitions"
	"os"
	"sync"
	"time"
)

type MetricsData struct {
	BytesNumber  []int64
	ResponseTime int64
}

var (
	Mtrx = make(chan MetricsData)
	Data = &MetricsData{}
)

func handleConnection(
	conn net.Conn,
	clientConn net.Conn,
	cnf definitions.LatencyInternalConfig,
	wg *sync.WaitGroup,
) {
	defer conn.Close()
	defer clientConn.Close()
	defer wg.Done()

	start := time.Now()

	localData := &MetricsData{}

	_, err := conn.Write([]byte("TCP server response\n"))
	if err != nil {
		fmt.Println("error sending data:", err)
		return
	}

	p1 := New(time.Duration(cnf.Delays.ReqDelay))
	p2 := New(time.Duration(cnf.Delays.ResDelay))

	var innerWg sync.WaitGroup
	innerWg.Add(3)

	go func() {
		defer innerWg.Done()
		Pipe(p1, conn, "tcp", "p1", localData)
	}()
	go func() {
		defer innerWg.Done()
		Pipe(clientConn, p1, "p1", "http", localData)
	}()
	go func() {
		defer innerWg.Done()
		Pipe(p2, clientConn, "http", "p2", localData)
	}()
	Pipe(conn, p2, "p2", "tcp", localData)

	innerWg.Wait()
	localData.ResponseTime = time.Since(start).Milliseconds()

	Data.BytesNumber = append(Data.BytesNumber, localData.BytesNumber...)
	Data.ResponseTime += localData.ResponseTime
}

type inner_proxy struct {
	data chan []byte
	l    time.Duration
}

func New(t time.Duration) *inner_proxy {
	return &inner_proxy{
		data: make(chan []byte),
		l:    t,
	}
}

func (p *inner_proxy) Read(b []byte) (int, error) {
	select {
	case data := <-p.data:
		time.Sleep(p.l)
		return copy(b, data), nil
	case <-time.After(1 * time.Second):
		return 0, io.EOF
	}
}

func (p *inner_proxy) Write(b []byte) (int, error) {
	p.data <- b
	return len(b), nil
}

func Pipe(dst io.Writer, src io.Reader, r, s string, mtr *MetricsData) {
	n, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("we have an error: %s ", err)
	}
	log.Printf("copied %d bytes from %s to %s \n", n, r, s)

	mtr.BytesNumber = append(mtr.BytesNumber, n)
}

// LaunchTcp runs until ctx is cancelled (externally via StopTcpServer) or the
// built-in 10-second timeout fires, whichever comes first.
func LaunchTcp(conf definitions.LatencyInternalConfig, ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+conf.Server.Port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port " + conf.Server.Port)

	var wg sync.WaitGroup

	// Combine the caller's context with a 10-second timeout.
	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-timeoutCtx.Done():
					return
				default:
					fmt.Println("Error accepting connection:", err)
					continue
				}
			}

			clientConn, err := net.Dial("tcp", conf.Client.Port)
			if err != nil {
				log.Printf("Error connecting to client: %v", err)
			}

			fmt.Println("client connected:", conn.RemoteAddr().String())

			wg.Add(1)
			go handleConnection(conn, clientConn, conf, &wg)
		}
	}()

	<-timeoutCtx.Done()

	listener.Close()
	wg.Wait()
	fmt.Println("Server has shut down gracefully, Collecting data...")
	Mtrx <- *Data
	return nil
}
