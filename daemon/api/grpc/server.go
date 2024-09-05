package main

import (
	"log"
	"net"
	"net/http"
	"obzev0/daemon/api/grpc/interceptors"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	rpcLogger *logrus.Entry
)

func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to start on port 50051: ", err)
	}

	rpcLogger := Logger()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(
					interceptors.RecoveryHandler(rpcLogger),
				),
			),
		),
	)
	reflection.Register(grpcServer)
	serviceAgent(grpcServer, rpcLogger)

	go waitForMetrics()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatal("Failed to serve metrics endpoint: ", err)
		}
	}()

	log.Printf("server listening at %v", l.Addr())
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal("Failed to serve gRPC over 50051: ", err)
	}
}
