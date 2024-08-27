// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: common/proto/tcAnalyser/tcAnalyser.proto

package tcAnalyser

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	TcAnalyserService_StartUserSpace_FullMethodName = "/tcAnalyserServ.TcAnalyserService/StartUserSpace"
)

// TcAnalyserServiceClient is the client API for TcAnalyserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TcAnalyserServiceClient interface {
	StartUserSpace(ctx context.Context, in *RequestForUserSpace, opts ...grpc.CallOption) (*ResponseFromUserSpace, error)
}

type tcAnalyserServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTcAnalyserServiceClient(cc grpc.ClientConnInterface) TcAnalyserServiceClient {
	return &tcAnalyserServiceClient{cc}
}

func (c *tcAnalyserServiceClient) StartUserSpace(ctx context.Context, in *RequestForUserSpace, opts ...grpc.CallOption) (*ResponseFromUserSpace, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResponseFromUserSpace)
	err := c.cc.Invoke(ctx, TcAnalyserService_StartUserSpace_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TcAnalyserServiceServer is the server API for TcAnalyserService service.
// All implementations must embed UnimplementedTcAnalyserServiceServer
// for forward compatibility.
type TcAnalyserServiceServer interface {
	StartUserSpace(context.Context, *RequestForUserSpace) (*ResponseFromUserSpace, error)
	mustEmbedUnimplementedTcAnalyserServiceServer()
}

// UnimplementedTcAnalyserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedTcAnalyserServiceServer struct{}

func (UnimplementedTcAnalyserServiceServer) StartUserSpace(context.Context, *RequestForUserSpace) (*ResponseFromUserSpace, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartUserSpace not implemented")
}
func (UnimplementedTcAnalyserServiceServer) mustEmbedUnimplementedTcAnalyserServiceServer() {}
func (UnimplementedTcAnalyserServiceServer) testEmbeddedByValue()                           {}

// UnsafeTcAnalyserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TcAnalyserServiceServer will
// result in compilation errors.
type UnsafeTcAnalyserServiceServer interface {
	mustEmbedUnimplementedTcAnalyserServiceServer()
}

func RegisterTcAnalyserServiceServer(s grpc.ServiceRegistrar, srv TcAnalyserServiceServer) {
	// If the following call pancis, it indicates UnimplementedTcAnalyserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&TcAnalyserService_ServiceDesc, srv)
}

func _TcAnalyserService_StartUserSpace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestForUserSpace)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TcAnalyserServiceServer).StartUserSpace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TcAnalyserService_StartUserSpace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TcAnalyserServiceServer).StartUserSpace(ctx, req.(*RequestForUserSpace))
	}
	return interceptor(ctx, in, info, handler)
}

// TcAnalyserService_ServiceDesc is the grpc.ServiceDesc for TcAnalyserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TcAnalyserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tcAnalyserServ.TcAnalyserService",
	HandlerType: (*TcAnalyserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartUserSpace",
			Handler:    _TcAnalyserService_StartUserSpace_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "common/proto/tcAnalyser/tcAnalyser.proto",
}