package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConnection(conn net.Conn, clientConn net.Conn, cnf Config) {
	defer conn.Close()
	_, err := conn.Write([]byte("TCP server response\n"))
	if err != nil {
		fmt.Println("error sending data:", err)
		return
	}

	deadline := time.Now().Add(40 * time.Second)
	conn.SetDeadline(deadline)
	clientConn.SetDeadline(deadline)

	p1 := New(time.Duration(cnf.Delays.ReqDelay))
	p2 := New(time.Duration(cnf.Delays.ResDelay))

	go Pipe(p1, conn, "tcp", "p1")
	go Pipe(clientConn, p1, "p1", "http")
	go Pipe(p2, clientConn, "http", "p2")
	Pipe(conn, p2, "p2", "tcp")

}

type inner_proxy struct {
	data chan []byte
	l    time.Duration
}

func New(t time.Duration) *inner_proxy {
	return &inner_proxy{
		make(chan []byte),
		t,
	}
}

func (p *inner_proxy) Read(b []byte) (int, error) {
	select {
	case data := <-p.data:
		time.Sleep(p.l)
		return copy(b, data), nil
	case <-time.After(4 * time.Second):
		return 0, io.EOF

	}

}
func (p *inner_proxy) Write(b []byte) (int, error) {
	p.data <- b
	return len(b), nil
}

func Pipe(
	dst io.Writer,
	src io.Reader,
	r, s string,
) {
	n, err := io.Copy(dst, src)
	if err != nil {
		log.Printf("we have an error: %s ", err)
	}
	log.Printf("copied %d bytes from %s to %s \n", n, r, s)

}

func main() {

	cnf, err := LoadConfig("tonConf.yaml")
	// tcp
	listener, err := net.Listen("tcp", ":"+cnf.Server.Port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()
	// client

	fmt.Println("TCP server listening on port " + cnf.Server.Port)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		clientConn, err := net.Dial("tcp", ":"+cnf.Client.Port)
		if err != nil {
			panic(err)
		}
		defer clientConn.Close()

		fmt.Println("client connected:", conn.RemoteAddr().String())

		go handleConnection(conn, clientConn, cnf)
	}
}
