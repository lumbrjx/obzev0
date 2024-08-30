package tcanalyser

import (
	"log"

	"obzev0/common/proto/tcAnalyser"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TcAnalyserService struct {
	tcAnalyser.UnimplementedTcAnalyserServiceServer
	// metrics     MetricsData
	// metricsChan chan MetricsData
}

func (s *TcAnalyserService) StartUserSpace(
	ctx context.Context,
	requestUserSpace *tcAnalyser.RequestForUserSpace,
) (*tcAnalyser.ResponseFromUserSpace, error) {
	if err := requestUserSpace.Config.Validate(); err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid request: %v",
			err,
		)
	}
	config := requestUserSpace.GetConfig()
	log.Printf("recived %s", config.Interface)

	go bpfLoader(config.Interface)
	// conf := definitions.Config{
	// 	Delays: definitions.DelaysConfig{
	// 		ReqDelay: config.ReqDelay,
	// 		ResDelay: config.ResDelay,
	// 	},
	// 	Server: definitions.ServerConfig{
	// 		Port: config.Server,
	// 	},
	// 	Client: definitions.ClientConfig{
	// 		Port: config.Client,
	// 	},
	// }
	// go LaunchTcp(conf)
	return &tcAnalyser.ResponseFromUserSpace{
		Message: "User Space program staus :",
	}, nil
}
