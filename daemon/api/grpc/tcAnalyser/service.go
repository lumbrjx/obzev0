package tcanalyser

import (
	"context"
	"log"
	"sync/atomic"

	proto "obzev0/common/proto/tcAnalyser"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TcAnalyserService struct {
	proto.UnimplementedTcAnalyserServiceServer
}

// activeCancel stores the context.CancelFunc for any running eBPF program.
var activeCancel atomic.Value

func (s *TcAnalyserService) StartUserSpace(
	ctx context.Context,
	requestUserSpace *proto.RequestForUserSpace,
) (*proto.ResponseFromUserSpace, error) {
	if err := requestUserSpace.Config.Validate(); err != nil {
		log.Printf("Invalid request: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	config := requestUserSpace.GetConfig()
	log.Printf("- Received TC Configuration: Interface=%s", config.Interface)

	expCtx, cancel := context.WithCancel(context.Background())
	activeCancel.Store(cancel)

	go func() {
		if err := bpfLoader(config.Interface, expCtx); err != nil {
			log.Printf("Error loading BPF program: %v", err)
		}
	}()

	return &proto.ResponseFromUserSpace{Message: "User Space program status: Loading"}, nil
}

func (s *TcAnalyserService) StopUserSpace(
	ctx context.Context,
	req *proto.StopRequest,
) (*proto.StopResponse, error) {
	log.Printf("StopUserSpace called, reason: %s", req.Reason)
	if cancel, ok := activeCancel.Load().(context.CancelFunc); ok {
		cancel()
	}
	return &proto.StopResponse{Message: "tcAnalyser experiment stopped"}, nil
}
