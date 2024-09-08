package packetmanipulation

import (
	"log"
	"time"

	"obzev0/common/proto/packetManipulation"
	"obzev0/daemon/api/grpc/helper"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PacketManipulationService struct {
	packetManipulation.UnimplementedPacketManipulationServiceServer
	// metrics     MetricsData
	// metricsChan chan MetricsData
}

func (s *PacketManipulationService) StartManipulationProxy(
	ctx context.Context,
	requestForManipulationProxy *packetManipulation.RequestForManipulationProxy,
) (*packetManipulation.ResponseFromManipulationProxy, error) {
	if err := requestForManipulationProxy.Config.Validate(); err != nil {
		log.Printf("Invalid request: %v", err)
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid request: %v",
			err,
		)
	}
	config := requestForManipulationProxy.GetConfig()

	log.Printf("- Received Client: %s", config.Client)
	log.Printf("- Received Server: %s", config.Server)
	log.Printf("- Received DurationConfig: %+v", config.DurationConfig)

	proxyConfiguration := ProxyConfig{
		Server: config.Server,
		Client: config.Client,
	}

	if config.DurationConfig.DurationSeconds > 0 {
		proxyConfiguration.DropRate = float64(config.DurationConfig.DropRate)
		proxyConfiguration.Timeout = time.Duration(
			config.DurationConfig.DurationSeconds,
		) * time.Second

		go func() {
			if err := Proxy(proxyConfiguration); err != nil {
				log.Printf("Error in manipulation Proxy: %v", err)
			}
		}()

		time.Sleep(2 * time.Second)
		if err := helper.ReqSimulator(config.Server, false, proxyConfiguration.Timeout); err != nil {
			log.Printf("Request simulation error: %v", err)
		}
	} else {
		log.Println("No duration set for manipulation. Proxy not started.")
	}

	return &packetManipulation.ResponseFromManipulationProxy{
		Message: "User Space program status: Proxy manipulation started",
	}, nil
}
