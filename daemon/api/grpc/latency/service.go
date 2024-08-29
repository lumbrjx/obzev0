package latency

import (
	"log"

	"obzev0/common/definitions"
	"obzev0/common/proto/latency"

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
	log.Printf("recived %s", config.Client)

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
	go LaunchTcp(conf)
	return &latency.ResponseFromTcp{
		Message: "Tcp server started",
	}, nil
}
