package main

import (
	ltc "obzev0/common/proto/latency"
	"obzev0/common/proto/packetManipulation"
	tcanl "obzev0/common/proto/tcAnalyser"
	"obzev0/daemon/api/grpc/latency"
	packetmanipulation "obzev0/daemon/api/grpc/packetManipulation"
	tcanalyser "obzev0/daemon/api/grpc/tcAnalyser"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func serviceAgent(grpcServer *grpc.Server, rpcLogger *logrus.Entry) {

	// Latency Service
	s := latency.LatencyService{}
	ltc.RegisterLatencyServiceServer(grpcServer, &s)

	// Traffic Controller Service
	tc := tcanalyser.TcAnalyserService{}
	tcanl.RegisterTcAnalyserServiceServer(grpcServer, &tc)

	pct := packetmanipulation.PacketManipulationService{}
	packetManipulation.RegisterPacketManipulationServiceServer(grpcServer, &pct)

	// Health Checking Serivce
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
	rpcLogger.Log(logrus.DebugLevel, "gRpc services have been established")
}
