package packetmanipulation

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	proto "obzev0/common/proto/packetManipulation"
	"obzev0/daemon/api/grpc/helper"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PacketManipulationService struct {
	proto.UnimplementedPacketManipulationServiceServer
}

// activeCancel stores the context.CancelFunc for any running experiment.
var activeCancel atomic.Value

func (s *PacketManipulationService) StartManipulationProxy(
	ctx context.Context,
	requestForManipulationProxy *proto.RequestForManipulationProxy,
) (*proto.ResponseFromManipulationProxy, error) {
	if err := requestForManipulationProxy.Config.Validate(); err != nil {
		log.Printf("Invalid request: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
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
		proxyConfiguration.Timeout = time.Duration(config.DurationConfig.DurationSeconds) * time.Second

		expCtx, cancel := context.WithCancel(context.Background())
		activeCancel.Store(cancel)

		go func() {
			if err := Proxy(proxyConfiguration, expCtx); err != nil {
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

	return &proto.ResponseFromManipulationProxy{
		Message: "User Space program status: Proxy manipulation started",
	}, nil
}

func (s *PacketManipulationService) StopManipulationProxy(
	ctx context.Context,
	req *proto.StopRequest,
) (*proto.StopResponse, error) {
	log.Printf("StopManipulationProxy called, reason: %s", req.Reason)
	if cancel, ok := activeCancel.Load().(context.CancelFunc); ok {
		cancel()
	}
	return &proto.StopResponse{Message: "packet manipulation experiment stopped"}, nil
}
