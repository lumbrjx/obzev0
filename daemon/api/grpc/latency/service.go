package latency

import (
	"log"
	"time"

	"obzev0/common/definitions"
	"obzev0/common/proto/latency"
	"obzev0/daemon/api/grpc/helper"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LatencyService struct {
	latency.UnimplementedLatencyServiceServer
	metrics     MetricsData
	metricsChan chan MetricsData
}

func (s *LatencyService) StartTcpServer(
	ctx context.Context,
	requestForTcp *latency.RequestForTcp,
) (*latency.ResponseFromTcp, error) {
	if err := requestForTcp.Config.Validate(); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid request: %v",
			err,
		)
	}

	config := requestForTcp.GetConfig()
	log.Printf("Received config: %v", config)

	conf := definitions.Config{
		Delays: definitions.DelaysConfig{
			ReqDelay: config.ReqDelay,
			ResDelay: config.ResDelay,
		},
		Server: definitions.ServerConfig{
			Port: config.Server,
		},
		Client: definitions.ClientConfig{
			Port: config.Client,
		},
	}

	go func() {
		if err := LaunchTcp(conf); err != nil {
			log.Printf("Error in LaunchTcp: %v", err)
		}
	}()

	go func() {
		time.Sleep(2 * time.Second)
		err := helper.ReqSimulator(
			config.Server,
			true,
			time.Duration(0)*time.Second,
		)
		if err != nil {

			log.Printf("Error in ReqSimulator: %v", err)
		}
	}()

	return &latency.ResponseFromTcp{
		Message: "TCP server Request Completed",
	}, nil
}
