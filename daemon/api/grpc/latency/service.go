package latency

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"obzev0/common/definitions"
	proto "obzev0/common/proto/latency"
	"obzev0/daemon/api/grpc/helper"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LatencyService struct {
	proto.UnimplementedLatencyServiceServer
	metrics     MetricsData
	metricsChan chan MetricsData
}

// activeCancel stores the context.CancelFunc for any running experiment so
// StopTcpServer can cancel it remotely.
var activeCancel atomic.Value

func (s *LatencyService) StartTcpServer(
	ctx context.Context,
	requestForTcp *proto.RequestForTcp,
) (*proto.ResponseFromTcp, error) {
	if err := requestForTcp.Config.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	config := requestForTcp.GetConfig()
	log.Printf("Received config: %v", config)

	conf := definitions.LatencyInternalConfig{
		Delays: definitions.DelaysConfig{
			ReqDelay: config.ReqDelay,
			ResDelay: config.ResDelay,
		},
		Server: definitions.ServerConfig{Port: config.Server},
		Client: definitions.ClientConfig{Port: config.Client},
	}

	expCtx, cancel := context.WithCancel(context.Background())
	activeCancel.Store(cancel)

	go func() {
		if err := LaunchTcp(conf, expCtx); err != nil {
			log.Printf("Error in LaunchTcp: %v", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		if err := helper.ReqSimulator(config.Server, true, time.Duration(0)*time.Second); err != nil {
			log.Printf("Error in ReqSimulator: %v", err)
		}
	}()

	return &proto.ResponseFromTcp{Message: "TCP server Request Completed"}, nil
}

func (s *LatencyService) StopTcpServer(
	ctx context.Context,
	req *proto.StopRequest,
) (*proto.StopResponse, error) {
	log.Printf("StopTcpServer called, reason: %s", req.Reason)
	if cancel, ok := activeCancel.Load().(context.CancelFunc); ok {
		cancel()
	}
	return &proto.StopResponse{Message: "latency experiment stopped"}, nil
}
