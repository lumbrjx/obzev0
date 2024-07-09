package main

import (
	"fmt"
	"log"
	"net"
	ltc "obzev0/common/proto/latency"
	"obzev0/daemon/api/grpc/latency"
	"os"
	"time"

	"google.golang.org/grpc"
)

func waitForMetrics() error {

	data := <-latency.Mtrx
	file, err := os.Create(
		"../../../latencyMetrics-" + time.Now().
			UTC().
			Format("01-06-02-15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("failed to create or open file: %w", err)
	}
	defer file.Close()

	bytesString := "Bytes number: "
	for _, num := range data.BytesNumber {
		bytesString += fmt.Sprintf("%d ", num)
	}
	responseTimeString := fmt.Sprintf(
		"Response time: %d ms\n",
		data.ResponseTime,
	)

	dataString := fmt.Sprintf("%s\n%s", bytesString, responseTimeString)

	_, err = file.WriteString(dataString)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil

}
func main() {

	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to start on port 9000: ", err)
	}

	s := latency.LatencyService{}
	grpcServer := grpc.NewServer()
	ltc.RegisterLatencyServiceServer(grpcServer, &s)
	go waitForMetrics()
	log.Printf("server listening at %v", l.Addr())
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Failed to serve grpc over 9000 ", err)
	}
}
