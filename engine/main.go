package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
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
	responseChan := make(chan ResponseT)

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
			"response: %s %s %d %s",
			resp.method,
			resp.url,
			resp.statusCode,
			resp.body,
		)

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

type ResponseT struct {
	method     string
	url        string
	statusCode int
	body       string
}

func p1(method, url string, rChan chan ResponseT, wg *sync.WaitGroup, c int) {
	defer wg.Done()
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	time.Sleep(time.Duration(c) * time.Second)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return
	}

	rs := ResponseT{
		method:     method,
		url:        url,
		body:       string(body),
		statusCode: resp.StatusCode,
	}
	rChan <- rs
}

func p2(resp ResponseT, rChan chan ResponseT, wg *sync.WaitGroup, c int) {
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
