package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"sync"
	"time"
)

func handleConnection(conn net.Conn, cnf Config) {
	defer conn.Close()
	fmt.Println("client connected:", conn.RemoteAddr().String())

	_, err := conn.Write([]byte("TCP server response\n"))
	if err != nil {
		fmt.Println("error sending data:", err)
		return
	}

	reader := bufio.NewReader(conn)
	responseChan := make(chan string)

	var wg sync.WaitGroup

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("disconnected:", conn.RemoteAddr().String())
			return
		}

		// here goes the reader routine
		if method, url, ok := ExtractURL(message); ok {
			wg.Add(1)
			go p1(method, url, responseChan, &wg, cnf.Delays.ReqDelay)
		}

		resp := <-responseChan

		wg.Add(1)
		go p2(resp, responseChan, &wg, cnf.Delays.ResDelay)

		resp = <-responseChan
		response := fmt.Sprintf(
			"response: %s", resp)

		fmt.Println(response)
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("error sending data:", err)
			return
		}
		go func() {
			wg.Wait()
			close(responseChan)
		}()

		break
	}
}

func p1(method, urL string, rChan chan string, wg *sync.WaitGroup, c int) {
	defer wg.Done()
	parsedURL, err := url.Parse(urL)
	if err != nil {
		panic(err)
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()

	if port == "" {
		if parsedURL.Scheme == "http" {
			port = "80"
		} else if parsedURL.Scheme == "https" {
			port = "443"
		} else {
			panic("Unsupported URL scheme")
		}
	}
	address := net.JoinHostPort(host, port)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request := fmt.Sprintf(
		"GET / HTTP/1.1\r\nHost: %s\r\nConnection: close\r\n\r\n",
		host,
	)

	_, err = conn.Write([]byte(request))
	if err != nil {
		panic(err)
	}

	response := ""
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		response += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Print the response
	rs := response
	rChan <- rs
}

func p2(resp string, rChan chan string, wg *sync.WaitGroup, c int) {
	defer wg.Done()

	// Some processing on the response
	time.Sleep(time.Duration(c) * time.Second)
	rChan <- resp
}

func main() {

	cnf, err := LoadConfig("tonConf.yaml")
	listener, err := net.Listen("tcp", ":"+cnf.Server.Port)
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		os.Exit(1)
	}
	defer listener.Close()

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

		go handleConnection(conn, cnf)
	}
}
