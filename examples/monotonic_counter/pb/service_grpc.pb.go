// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CounterClient is the client API for Counter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CounterClient interface {
	// Merges counts from another node
	Peer(ctx context.Context, opts ...grpc.CallOption) (Counter_PeerClient, error)
	Value(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*ValueResponse, error)
}

type counterClient struct {
	cc grpc.ClientConnInterface
}

func NewCounterClient(cc grpc.ClientConnInterface) CounterClient {
	return &counterClient{cc}
}

func (c *counterClient) Peer(ctx context.Context, opts ...grpc.CallOption) (Counter_PeerClient, error) {
	stream, err := c.cc.NewStream(ctx, &Counter_ServiceDesc.Streams[0], "/monotonic_counter.Counter/Peer", opts...)
	if err != nil {
		return nil, err
	}
	x := &counterPeerClient{stream}
	return x, nil
}

type Counter_PeerClient interface {
	Send(*MergeRequest) error
	Recv() (*MergeResponse, error)
	grpc.ClientStream
}

type counterPeerClient struct {
	grpc.ClientStream
}

func (x *counterPeerClient) Send(m *MergeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *counterPeerClient) Recv() (*MergeResponse, error) {
	m := new(MergeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *counterClient) Value(ctx context.Context, in *ValueRequest, opts ...grpc.CallOption) (*ValueResponse, error) {
	out := new(ValueResponse)
	err := c.cc.Invoke(ctx, "/monotonic_counter.Counter/Value", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CounterServer is the server API for Counter service.
// All implementations must embed UnimplementedCounterServer
// for forward compatibility
type CounterServer interface {
	// Merges counts from another node
	Peer(Counter_PeerServer) error
	Value(context.Context, *ValueRequest) (*ValueResponse, error)
	mustEmbedUnimplementedCounterServer()
}

// UnimplementedCounterServer must be embedded to have forward compatible implementations.
type UnimplementedCounterServer struct {
}

func (UnimplementedCounterServer) Peer(Counter_PeerServer) error {
	return status.Errorf(codes.Unimplemented, "method Peer not implemented")
}
func (UnimplementedCounterServer) Value(context.Context, *ValueRequest) (*ValueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Value not implemented")
}
func (UnimplementedCounterServer) mustEmbedUnimplementedCounterServer() {}

// UnsafeCounterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CounterServer will
// result in compilation errors.
type UnsafeCounterServer interface {
	mustEmbedUnimplementedCounterServer()
}

func RegisterCounterServer(s grpc.ServiceRegistrar, srv CounterServer) {
	s.RegisterService(&Counter_ServiceDesc, srv)
}

func _Counter_Peer_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CounterServer).Peer(&counterPeerServer{stream})
}

type Counter_PeerServer interface {
	Send(*MergeResponse) error
	Recv() (*MergeRequest, error)
	grpc.ServerStream
}

type counterPeerServer struct {
	grpc.ServerStream
}

func (x *counterPeerServer) Send(m *MergeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *counterPeerServer) Recv() (*MergeRequest, error) {
	m := new(MergeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Counter_Value_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CounterServer).Value(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/monotonic_counter.Counter/Value",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CounterServer).Value(ctx, req.(*ValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Counter_ServiceDesc is the grpc.ServiceDesc for Counter service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Counter_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "monotonic_counter.Counter",
	HandlerType: (*CounterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Value",
			Handler:    _Counter_Value_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Peer",
			Handler:       _Counter_Peer_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "service.proto",
}
