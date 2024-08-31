package controller

import (
	"context"
	"fmt"
	"obzev0/common/proto/latency"
	"time"

	pb "obzev0/common/proto/latency"
	tca "obzev0/common/proto/tcAnalyser"
	v "obzev0/controller/api/v1"

	"google.golang.org/grpc"
)

type GrpcServiceConfig struct {
	LatencyConfig v.TcpConfig
	TcAConfig     v.TcAnalyserConfig

	// Add more fields as needed
}

func callGrpcServices(
	conn *grpc.ClientConn,
	config GrpcServiceConfig,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Handle LatencyService gRPC call
	client := pb.NewLatencyServiceClient(conn)
	response, err := client.StartTcpServer(
		ctx,
		&pb.RequestForTcp{Config: &latency.TcpConfig{
			ReqDelay: config.LatencyConfig.ReqDelay,
			ResDelay: config.LatencyConfig.ResDelay,
			Server:   config.LatencyConfig.Server,
			Client:   config.LatencyConfig.Client,
		}},
	)
	if err != nil {
		return fmt.Errorf("error calling StartTcpServer: %w", err)
	}
	fmt.Printf("Response from LatencyService gRPC server: %s\n", response.Message)

	// Handle TcAnalyserService gRPC call
	client2 := tca.NewTcAnalyserServiceClient(conn)
	rsp, err := client2.StartUserSpace(
		ctx,
		&tca.RequestForUserSpace{Config: &tca.TcConfig{
			Interface: config.TcAConfig.NetIFace,
		}},
	)
	if err != nil {
		return fmt.Errorf("error calling StartUserSpace: %w", err)
	}
	fmt.Printf("Response from TcAnalyserService gRPC server: %s\n", rsp.Message)

	return nil
}
