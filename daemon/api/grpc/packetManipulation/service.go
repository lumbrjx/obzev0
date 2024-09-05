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
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid request: %v",
			err,
		)
	}
	config := requestForManipulationProxy.GetConfig()

	log.Printf("recived %s", config.Client)
	log.Printf("recived %s", config.Server)
	log.Printf("recived %s", config.DurationConfig)

	proxyConfiguration := ProxyConfig{
		Server: config.Server,
		Client: config.Client,
	}

	// Start proxy in a goroutine
	if config.DurationConfig.DurationSeconds != 0 {
		proxyConfiguration.DropRate = float64(config.DurationConfig.DropRate)
		proxyConfiguration.Timeout = time.Duration(
			config.DurationConfig.DurationSeconds,
		)

		go func() {
			err := Proxy(proxyConfiguration)
			if err != nil {
				log.Printf("Error in manipulation Proxy: %v", err)
			}
		}()
		time.Sleep(2 * time.Second)
		err := helper.ReqSimulator(
			config.Server,
			false,
			time.Duration(
				config.DurationConfig.DurationSeconds,
			)*time.Second,
		)
		if err != nil {
			log.Printf("Request simulation error: %v", err)
		}
	}

	return &packetManipulation.ResponseFromManipulationProxy{
		Message: "User Space program staus :",
	}, nil
}
