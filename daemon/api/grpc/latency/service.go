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
	if requestForTcp == nil || requestForTcp.Config == nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"RequestForTcp or TcpConfig cannot be nil",
		)
	}

	config := requestForTcp.GetConfig()
	if config.ReqDelay == 0 || config.ResDelay == 0 || config.Server == "" ||
		config.Client == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"All fields in TcpConfig must be provided",
		)
	}
	log.Printf("recived %s", requestForTcp.Config.Client)

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
