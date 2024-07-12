package main

import (
	"log"
	"net"
	"net/http"
	ltc "obzev0/common/proto/latency"
	"obzev0/daemon/api/grpc/latency"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var (
	bytesHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_bytes",
			Help:    "Bytes processed by the gRPC server",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	responseTimeHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_response_time_seconds",
			Help:    "Response time of gRPC server methods",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func recordMetrics(method string, bytes int64, responseTime time.Duration) {
	log.Printf(
		"Recording metrics: method=%s, bytes=%d, responseTime=%s",
		method,
		bytes,
		responseTime,
	)
	bytesHistogram.WithLabelValues(method).Observe(float64(bytes))
	responseTimeHistogram.WithLabelValues(method).Observe(responseTime.Seconds())
}

func waitForMetrics() {
	for {
		data := <-latency.Mtrx
		log.Printf("Received data: %+v", data)

		for _, bytes := range data.BytesNumber {
			recordMetrics(
				"LatencyService",
				bytes,
				time.Duration(data.ResponseTime),
			)
		}
	}
}

func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to start on port 50051: ", err)
	}

	s := latency.LatencyService{}
	grpcServer := grpc.NewServer()
	ltc.RegisterLatencyServiceServer(grpcServer, &s)

	healthSrv := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthSrv)

	healthSrv.SetServingStatus(
		"grpc.health.v1.Health",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)
	healthSrv.SetServingStatus(
		"obzev0.common.proto.latency.LatencyService",
		grpc_health_v1.HealthCheckResponse_SERVING,
	)

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
