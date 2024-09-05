package main

import (
	"log"
	"obzev0/daemon/api/grpc/latency"
	packetmanipulation "obzev0/daemon/api/grpc/packetManipulation"
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
	droppedPacketsHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_dropped_packets_count",
			Help:    "Dropped packets count during packet manipulation",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	corruptedPacketsHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_server_corrupted_packets_count",
			Help:    "Corrupted packets count during packet manipulation",
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

func processPacketManipulationData() {
	select {
	case dropData := <-packetmanipulation.Mtrx:
		log.Printf("Received packet manipulation data: %+v", dropData)
		droppedPacketsHistogram.WithLabelValues("PacketManipulationService").
			Observe(float64(dropData.DropedCount))
		corruptedPacketsHistogram.WithLabelValues("PacketManipulationService").
			Observe(float64(dropData.CorruptedCount))
	default:

	}
}

func processLatencyData() {
	select {
	case data := <-latency.Mtrx:
		log.Printf("Received latency data: %+v", data)
		for _, bytes := range data.BytesNumber {
			recordMetrics(
				"LatencyService",
				bytes,
				time.Duration(data.ResponseTime),
			)
		}
	default:
	}
}

func waitForMetrics() {
	for {
		processPacketManipulationData()
		processLatencyData()
		time.Sleep(1 * time.Second)
	}
}
