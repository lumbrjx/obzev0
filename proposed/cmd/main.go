package main

import (
	"flag"
	"io"
	"log"
	"net"
	"time"
)

var (
	DefaultTimeOut = 10 * time.Second
)

func main() {
	var listen_addr = flag.String("addr", "127.0.0.1:9090", "tcp addresse that our tcp server listen to")
	var upstream_addr = flag.String("upstream", "127.0.0.1:8080", "upstream addr")
	flag.Parse()
	listner, err := net.Listen("tcp", *listen_addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listner.Close()
	for {
		tcpConn, err := listner.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func(con net.Conn) {
			defer tcpConn.Close()
			p1 := New(2 * time.Millisecond)
			p2 := New(1 * time.Millisecond)
			httpConn, err := net.Dial("tcp", *upstream_addr)
			if err != nil {
				log.Println(err)
				return
			}

			defer httpConn.Close()

			deadline := time.Now().Add(DefaultTimeOut)
			tcpConn.SetDeadline(deadline)
			httpConn.SetDeadline(deadline)

			go tunnel(p1, tcpConn, "proxy-1", "tcp-proxy")
			go tunnel(httpConn, p1, "http-server", "proxy-1")
			go tunnel(p2, httpConn, "proxy-2", "http-server")
			tunnel(tcpConn, p2, "tcp-prxy", "peoxy-2")

		}(tcpConn)

	}
}

type inner_proxy struct {
	data chan []byte
	// we assume we are using this param only in the reading phase
	l time.Duration
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
	case <-time.After(DefaultTimeOut):
		return 0, io.EOF

	}

}
func (p *inner_proxy) Write(b []byte) (int, error) {
	p.data <- b
	return len(b), nil
}

func tunnel(dst io.Writer, src io.Reader, srcName, dstName string) {

	n, err := io.Copy(dst, src)
	if err != nil {

		log.Printf("error : %v  [%s] --> [%s] ", err, srcName, dstName)
	}
	log.Printf("copied %d bytes  [%s] --> [%s]  \n", n, srcName, dstName)

}
