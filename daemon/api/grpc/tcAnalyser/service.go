package tcanalyser

import (
	"context"
	"log"

	"obzev0/common/proto/tcAnalyser"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TcAnalyserService struct {
	tcAnalyser.UnimplementedTcAnalyserServiceServer
}

func (s *TcAnalyserService) StartUserSpace(
	ctx context.Context,
	requestUserSpace *tcAnalyser.RequestForUserSpace,
) (*tcAnalyser.ResponseFromUserSpace, error) {
	if err := requestUserSpace.Config.Validate(); err != nil {
		log.Printf("Invalid request: %v", err)
		return nil, status.Errorf(
			codes.InvalidArgument,
			"Invalid request: %v",
			err,
		)
	}

	config := requestUserSpace.GetConfig()
	log.Printf("- Received TC Configuration: Interface=%s", config.Interface)

	go func() {
		if err := bpfLoader(config.Interface); err != nil {
			log.Printf(
				"Error loading BPF program: %v",
				err,
			)
		}
	}()

	return &tcAnalyser.ResponseFromUserSpace{
		Message: "User Space program status: Loading",
	}, nil
}
