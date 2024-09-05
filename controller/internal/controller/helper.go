package controller

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func strToFlt(str string) float32 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
		return 0
	}
	return float32(f)
}

func LoggingInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// Log the RPC method name and request
	log.Printf("Invoking gRPC method: %s with request: %v", method, req)

	// Call the RPC
	err := invoker(ctx, method, req, reply, cc, opts...)

	// Log the result of the RPC call
	if err != nil {
		log.Printf("Error calling gRPC method: %s with error: %v", method, err)
	} else {
		log.Printf("gRPC method: %s returned response: %v", method, reply)
	}

	return err
}

func CheckConnection(conn *grpc.ClientConn) {
	state := conn.GetState()
	switch state {
	case connectivity.Idle:
		log.Println("Connection is idle")
	case connectivity.Connecting:
		log.Println("Connection is in the process of connecting")
	case connectivity.Ready:
		log.Println("Connection is ready and open")
	case connectivity.TransientFailure:
		log.Println("Connection is experiencing a transient failure")
	case connectivity.Shutdown:
		log.Println("Connection is closed or shutting down")
	default:
		log.Println("Unknown connection state")
	}
}
