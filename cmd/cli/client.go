package main

import (
	"context"
	"fmt"
	"log"
	"obzev0/common/definitions"
	"obzev0/common/proto/latency"
	"obzev0/common/proto/packetManipulation"
	"obzev0/common/proto/tcAnalyser"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func client(addr string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}

func apply(rpcConfig definitions.Config) {
	conn := client(rpcConfig.ServerAddr)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	defer conn.Close()

	if rpcConfig.LatencySvcConfig.Enabled {

		c := latency.NewLatencyServiceClient(conn)
		response, err := c.StartTcpServer(
			ctx,
			&latency.RequestForTcp{Config: &latency.TcpConfig{
				ReqDelay: int32(rpcConfig.LatencySvcConfig.ReqDelay),
				ResDelay: int32(rpcConfig.LatencySvcConfig.ReqDelay),
				Server:   rpcConfig.LatencySvcConfig.Server,
				Client:   rpcConfig.LatencySvcConfig.Client,
			}},
		)
		if err != nil {
			log.Printf("Error calling StartTcpServer: %v", err)
			log.Fatal("error calling StartTcpServer: %w", err)
		}
		fmt.Printf(
			"Response from LatencyService gRPC server: %s\n",
			response.Message,
		)
	}
	if rpcConfig.TcAnalyserSvcConfig.Enabled {

		client2 := tcAnalyser.NewTcAnalyserServiceClient(conn)
		rsp, err := client2.StartUserSpace(
			ctx,
			&tcAnalyser.RequestForUserSpace{Config: &tcAnalyser.TcConfig{
				Interface: rpcConfig.TcAnalyserSvcConfig.NetIFace,
			}},
		)
		if err != nil {
			log.Fatal("error calling StartUserSpace: %w", err)
		}
		fmt.Printf(
			"Response from TcAnalyserService gRPC server: %s\n",
			rsp.Message,
		)
	}

	if rpcConfig.PacketManipulationSvcConfig.Enabled {
		client3 := packetManipulation.NewPacketManipulationServiceClient(conn)
		s := strToFlt(rpcConfig.PacketManipulationSvcConfig.DropRate)
		r := strToFlt(rpcConfig.PacketManipulationSvcConfig.CorruptRate)

		d, err := client3.StartManipulationProxy(
			ctx,
			&packetManipulation.RequestForManipulationProxy{
				Config: &packetManipulation.PctmConfig{
					Server: rpcConfig.PacketManipulationSvcConfig.Server,
					Client: rpcConfig.PacketManipulationSvcConfig.Client,

					DurationConfig: &packetManipulation.DurationConfig{
						DurationSeconds: int32(
							rpcConfig.PacketManipulationSvcConfig.DurationSeconds,
						),
						DropRate:    s,
						CorruptRate: r,
					},
				},
			},
		)
		if err != nil {
			log.Fatal("error calling StartManipulationProxy: %w", err)
		}
		fmt.Printf(
			"Response from packetManipulationService gRPC server: %s\n",
			d.Message,
		)
	}

}
