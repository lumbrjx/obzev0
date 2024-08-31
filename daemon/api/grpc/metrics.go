package main

import (
	"log"
	"obzev0/daemon/api/grpc/latency"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
