package tcpService

import (
	"log"
	"obzev0/shared/structs"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedTcpServiceServer
}

func (s *Server) StartTcpServer(
	ctx context.Context,
	requestForTcp *RequestForTcp,
) (*ResponseFromTcp, error) {
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

	cnf := structs.Config{
		Delays: structs.DelaysConfig{
			ReqDelay: config.ReqDelay,
			ResDelay: config.ResDelay,
		},
		Server: structs.ServerConfig{
			Port: config.Server,
		},
		Client: structs.ClientConfig{
			Port: config.Client,
		},
	}
	go LaunchTcp(cnf)
	return &ResponseFromTcp{
		Metric: "Tcp server started",
	}, nil
}
