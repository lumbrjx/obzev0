package httpfault

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	proto "obzev0/common/proto/httpFault"
)

// HTTPFaultService implements proto.HTTPFaultServiceServer.
type HTTPFaultService struct {
	proto.UnimplementedHTTPFaultServiceServer
	activeServer atomic.Value // *http.Server
}

func (s *HTTPFaultService) StartHTTPFault(_ context.Context, req *proto.HTTPFaultRequest) (*proto.HTTPFaultResponse, error) {
	// Stop any running proxy first.
	s.stopServer()

	listenAddr := req.ListenAddr
	if listenAddr == "" {
		listenAddr = ":8880"
	}
	targetURL := req.TargetUrl
	if targetURL == "" {
		return nil, fmt.Errorf("target_url is required")
	}

	proxy, err := newFaultProxy(targetURL, req.ErrorRate, req.ErrorCode, req.DelayMs, req.AbortRate)
	if err != nil {
		return nil, fmt.Errorf("StartHTTPFault: %w", err)
	}

	srv := &http.Server{Addr: listenAddr, Handler: proxy}
	s.activeServer.Store(srv)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Server closed via Shutdown is normal.
		}
	}()

	return &proto.HTTPFaultResponse{
		Message: fmt.Sprintf("HTTP fault proxy listening on %s → %s (error_rate=%.2f, delay=%dms, abort_rate=%.2f)",
			listenAddr, targetURL, req.ErrorRate, req.DelayMs, req.AbortRate),
	}, nil
}

func (s *HTTPFaultService) StopHTTPFault(_ context.Context, _ *proto.StopRequest) (*proto.StopResponse, error) {
	s.stopServer()
	return &proto.StopResponse{Message: "HTTP fault proxy stopped"}, nil
}

func (s *HTTPFaultService) stopServer() {
	if v := s.activeServer.Swap(nil); v != nil {
		srv := v.(*http.Server)
		srv.Shutdown(context.Background())
	}
}
