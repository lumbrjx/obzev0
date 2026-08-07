package controller

import (
	"context"
	"fmt"
	"log"
	"obzev0/common/proto/latency"
	"time"

	pb "obzev0/common/proto/latency"
	httpfaultproto "obzev0/common/proto/httpFault"
	netchaoproto "obzev0/common/proto/networkChaos"
	pca "obzev0/common/proto/packetManipulation"
	tca "obzev0/common/proto/tcAnalyser"
	v "obzev0/controller/api/v1"

	"google.golang.org/grpc"
)

type GrpcServiceConfig struct {
	LatencyConfig    v.TcpConfig
	TcAConfig        v.TcAnalyserConfig
	PctmConfig       v.PacketManipulationConfig
	NetworkChaosConf v.NetworkChaosConfig
	HTTPFaultConf    v.HTTPFaultConfig
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
				Client: config.PctmConfig.Client,

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

	if config.NetworkChaosConf.Enabled {
		ncClient := netchaoproto.NewNetworkChaosServiceClient(conn)

		if config.NetworkChaosConf.BandwidthRateKbps > 0 {
			resp, err := ncClient.StartBandwidthLimit(ctx, &netchaoproto.BandwidthRequest{
				Interface: config.NetworkChaosConf.Interface,
				RateKbps:  config.NetworkChaosConf.BandwidthRateKbps,
				BurstKb:   uint32(config.NetworkChaosConf.BandwidthRateKbps / 8),
			})
			if err != nil {
				log.Printf("Error calling StartBandwidthLimit: %v", err)
			} else {
				fmt.Printf("Response from NetworkChaosService (bandwidth): %s\n", resp.Message)
			}
		}

		if config.NetworkChaosConf.DNSChaosMode != "" {
			resp, err := ncClient.StartDNSChaos(ctx, &netchaoproto.DNSChaosRequest{
				Mode:       config.NetworkChaosConf.DNSChaosMode,
				DelayMs:    config.NetworkChaosConf.DNSDelayMs,
				ListenAddr: config.NetworkChaosConf.DNSListenAddr,
				Upstream:   config.NetworkChaosConf.DNSUpstream,
			})
			if err != nil {
				log.Printf("Error calling StartDNSChaos: %v", err)
			} else {
				fmt.Printf("Response from NetworkChaosService (DNS): %s\n", resp.Message)
			}
		}

		if config.NetworkChaosConf.TCPResetListenAddr != "" {
			resp, err := ncClient.StartTCPReset(ctx, &netchaoproto.TCPResetRequest{
				ListenAddr: config.NetworkChaosConf.TCPResetListenAddr,
				ResetRate:  config.NetworkChaosConf.TCPResetRate,
			})
			if err != nil {
				log.Printf("Error calling StartTCPReset: %v", err)
			} else {
				fmt.Printf("Response from NetworkChaosService (TCP RST): %s\n", resp.Message)
			}
		}
	}

	if config.HTTPFaultConf.Enabled {
		hfClient := httpfaultproto.NewHTTPFaultServiceClient(conn)
		resp, err := hfClient.StartHTTPFault(ctx, &httpfaultproto.HTTPFaultRequest{
			ListenAddr: config.HTTPFaultConf.ListenAddr,
			TargetUrl:  config.HTTPFaultConf.TargetURL,
			ErrorRate:  config.HTTPFaultConf.ErrorRate,
			ErrorCode:  config.HTTPFaultConf.ErrorCode,
			DelayMs:    config.HTTPFaultConf.DelayMs,
			AbortRate:  config.HTTPFaultConf.AbortRate,
		})
		if err != nil {
			log.Printf("Error calling StartHTTPFault: %v", err)
		} else {
			fmt.Printf("Response from HTTPFaultService: %s\n", resp.Message)
		}
	}

	return nil
}
