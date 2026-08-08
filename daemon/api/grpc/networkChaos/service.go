package networkchaos

import (
	"context"
	"fmt"
	"sync/atomic"

	proto "obzev0/common/proto/networkChaos"
)

// NetworkChaosService implements proto.NetworkChaosServiceServer.
type NetworkChaosService struct {
	proto.UnimplementedNetworkChaosServiceServer
	bandwidthCancel atomic.Value
	dnsCancel       atomic.Value
	tcpResetCancel  atomic.Value
}

func (s *NetworkChaosService) StartBandwidthLimit(_ context.Context, req *proto.BandwidthRequest) (*proto.BandwidthResponse, error) {
	iface := req.Interface
	if iface == "" {
		iface = "eth0"
	}
	if err := startBandwidth(iface, req.RateKbps, req.BurstKb); err != nil {
		return nil, fmt.Errorf("StartBandwidthLimit: %w", err)
	}
	s.bandwidthCancel.Store(iface)
	return &proto.BandwidthResponse{Message: fmt.Sprintf("bandwidth limited on %s at %d kbps", iface, req.RateKbps)}, nil
}

func (s *NetworkChaosService) StopBandwidthLimit(_ context.Context, req *proto.StopRequest) (*proto.StopResponse, error) {
	iface := ""
	if v := s.bandwidthCancel.Load(); v != nil {
		iface, _ = v.(string)
	}
	if iface == "" {
		iface = "eth0"
	}
	if err := stopBandwidth(iface); err != nil {
		return nil, fmt.Errorf("StopBandwidthLimit: %w", err)
	}
	return &proto.StopResponse{Message: "bandwidth limit removed"}, nil
}

func (s *NetworkChaosService) StartDNSChaos(_ context.Context, req *proto.DNSChaosRequest) (*proto.DNSChaosResponse, error) {
	listenAddr := req.ListenAddr
	if listenAddr == "" {
		listenAddr = "127.0.0.1:5353"
	}
	upstream := req.Upstream
	if upstream == "" {
		upstream = "8.8.8.8:53"
	}

	cancel, err := startDNSChaos(req.Mode, req.DelayMs, listenAddr, upstream)
	if err != nil {
		return nil, fmt.Errorf("StartDNSChaos: %w", err)
	}

	if old := s.dnsCancel.Swap(cancel); old != nil {
		old.(context.CancelFunc)()
	}
	return &proto.DNSChaosResponse{Message: fmt.Sprintf("DNS chaos (%s) on %s", req.Mode, listenAddr)}, nil
}

func (s *NetworkChaosService) StopDNSChaos(_ context.Context, _ *proto.StopRequest) (*proto.StopResponse, error) {
	if v := s.dnsCancel.Swap(nil); v != nil {
		v.(context.CancelFunc)()
	}
	return &proto.StopResponse{Message: "DNS chaos stopped"}, nil
}

func (s *NetworkChaosService) StartTCPReset(_ context.Context, req *proto.TCPResetRequest) (*proto.TCPResetResponse, error) {
	listenAddr := req.ListenAddr
	if listenAddr == "" {
		listenAddr = ":9080"
	}

	cancel, err := startTCPReset(listenAddr, req.ResetRate)
	if err != nil {
		return nil, fmt.Errorf("StartTCPReset: %w", err)
	}

	if old := s.tcpResetCancel.Swap(cancel); old != nil {
		old.(context.CancelFunc)()
	}
	return &proto.TCPResetResponse{Message: fmt.Sprintf("TCP reset proxy on %s (rate=%.2f)", listenAddr, req.ResetRate)}, nil
}

func (s *NetworkChaosService) StopTCPReset(_ context.Context, _ *proto.StopRequest) (*proto.StopResponse, error) {
	if v := s.tcpResetCancel.Swap(nil); v != nil {
		v.(context.CancelFunc)()
	}
	return &proto.StopResponse{Message: "TCP reset stopped"}, nil
}
