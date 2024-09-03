package controller

import (
	"context"
	"fmt"
	"log"
	"obzev0/common/proto/latency"
	"time"

	pb "obzev0/common/proto/latency"
	pca "obzev0/common/proto/packetManipulation"
	tca "obzev0/common/proto/tcAnalyser"
	v "obzev0/controller/api/v1"

	"google.golang.org/grpc"
)

type GrpcServiceConfig struct {
	LatencyConfig v.TcpConfig
	TcAConfig     v.TcAnalyserConfig
	PctmConfig    v.PacketManipulationConfig

	// Add more fields as needed
}

func callGrpcServices(
	conn *grpc.ClientConn,
	config GrpcServiceConfig,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()

	// Handle LatencyService gRPC call
	if config.LatencyConfig.Enabled {
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
			log.Printf("Error calling StartTcpServer: %v", err)
			return fmt.Errorf("error calling StartTcpServer: %w", err)
		}
		fmt.Printf(
			"Response from LatencyService gRPC server: %s\n",
			response.Message,
		)

	}

	// Handle TcAnalyserService gRPC call
	if config.TcAConfig.Enabled {
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
		fmt.Printf(
			"Response from TcAnalyserService gRPC server: %s\n",
			rsp.Message,
		)
	}

	if config.PctmConfig.Enabled {
		client3 := pca.NewPacketManipulationServiceClient(conn)
		s := strToFlt(config.PctmConfig.DropRate)
		r := strToFlt(config.PctmConfig.CorruptRate)

		d, err := client3.StartManipulationProxy(
			ctx,
			&pca.RequestForManipulationProxy{Config: &pca.PctmConfig{
				Server: config.PctmConfig.Server,
				Client: config.LatencyConfig.Server,

				DurationConfig: &pca.DurationConfig{
					DurationSeconds: config.PctmConfig.DurationSeconds,
					DropRate:        s,
					CorruptRate:     r,
				},
			}},
		)
		if err != nil {
			return fmt.Errorf("error calling StartManipulationProxy: %w", err)
		}
		fmt.Printf(
			"Response from packetManipulationService gRPC server: %s\n",
			d.Message,
		)
	}

	return nil
}
