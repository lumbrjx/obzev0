package main

import (
	"log"
	"net"
	"net/http"
	ltc "obzev0/common/proto/latency"
	tcanl "obzev0/common/proto/tcAnalyser"
	"obzev0/daemon/api/grpc/interceptors"
	"obzev0/daemon/api/grpc/latency"
	tcanalyser "obzev0/daemon/api/grpc/tcAnalyser"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
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

var (
	rpcLogger *logrus.Entry
)

func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to start on port 50051: ", err)
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logger.SetOutput(os.Stdout)

	rpcLogger = logger.WithFields(logrus.Fields{
		"service": "gRPC/server",
	})
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(
					interceptors.RecoveryHandler(rpcLogger),
				),
			),
		),
	)
	s := latency.LatencyService{}
	ltc.RegisterLatencyServiceServer(grpcServer, &s)

	tc := tcanalyser.TcAnalyserService{}
	tcanl.RegisterTcAnalyserServiceServer(grpcServer, &tc)

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
