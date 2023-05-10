// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.1
// source: remote.proto

package remote

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

// RemotingClient is the client API for Remoting service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RemotingClient interface {
	Receive(ctx context.Context, opts ...grpc.CallOption) (Remoting_ReceiveClient, error)
	ListProcesses(ctx context.Context, in *ListProcessesRequest, opts ...grpc.CallOption) (*ListProcessesResponse, error)
	GetProcessDiagnostics(ctx context.Context, in *GetProcessDiagnosticsRequest, opts ...grpc.CallOption) (*GetProcessDiagnosticsResponse, error)
}

type remotingClient struct {
	cc grpc.ClientConnInterface
}

func NewRemotingClient(cc grpc.ClientConnInterface) RemotingClient {
	return &remotingClient{cc}
}

func (c *remotingClient) Receive(ctx context.Context, opts ...grpc.CallOption) (Remoting_ReceiveClient, error) {
	stream, err := c.cc.NewStream(ctx, &Remoting_ServiceDesc.Streams[0], "/remote.Remoting/Receive", opts...)
	if err != nil {
		return nil, err
	}
	x := &remotingReceiveClient{stream}
	return x, nil
}

type Remoting_ReceiveClient interface {
	Send(*RemoteMessage) error
	Recv() (*RemoteMessage, error)
	grpc.ClientStream
}

type remotingReceiveClient struct {
	grpc.ClientStream
}

func (x *remotingReceiveClient) Send(m *RemoteMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *remotingReceiveClient) Recv() (*RemoteMessage, error) {
	m := new(RemoteMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *remotingClient) ListProcesses(ctx context.Context, in *ListProcessesRequest, opts ...grpc.CallOption) (*ListProcessesResponse, error) {
	out := new(ListProcessesResponse)
	err := c.cc.Invoke(ctx, "/remote.Remoting/ListProcesses", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *remotingClient) GetProcessDiagnostics(ctx context.Context, in *GetProcessDiagnosticsRequest, opts ...grpc.CallOption) (*GetProcessDiagnosticsResponse, error) {
	out := new(GetProcessDiagnosticsResponse)
	err := c.cc.Invoke(ctx, "/remote.Remoting/GetProcessDiagnostics", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RemotingServer is the server API for Remoting service.
// All implementations must embed UnimplementedRemotingServer
// for forward compatibility
type RemotingServer interface {
	Receive(Remoting_ReceiveServer) error
	ListProcesses(context.Context, *ListProcessesRequest) (*ListProcessesResponse, error)
	GetProcessDiagnostics(context.Context, *GetProcessDiagnosticsRequest) (*GetProcessDiagnosticsResponse, error)
	mustEmbedUnimplementedRemotingServer()
}

// UnimplementedRemotingServer must be embedded to have forward compatible implementations.
type UnimplementedRemotingServer struct {
}

func (UnimplementedRemotingServer) Receive(Remoting_ReceiveServer) error {
	return status.Errorf(codes.Unimplemented, "method Receive not implemented")
}
func (UnimplementedRemotingServer) ListProcesses(context.Context, *ListProcessesRequest) (*ListProcessesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProcesses not implemented")
}
func (UnimplementedRemotingServer) GetProcessDiagnostics(context.Context, *GetProcessDiagnosticsRequest) (*GetProcessDiagnosticsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProcessDiagnostics not implemented")
}
func (UnimplementedRemotingServer) mustEmbedUnimplementedRemotingServer() {}

// UnsafeRemotingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RemotingServer will
// result in compilation errors.
type UnsafeRemotingServer interface {
	mustEmbedUnimplementedRemotingServer()
}

func RegisterRemotingServer(s grpc.ServiceRegistrar, srv RemotingServer) {
	s.RegisterService(&Remoting_ServiceDesc, srv)
}

func _Remoting_Receive_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RemotingServer).Receive(&remotingReceiveServer{stream})
}

type Remoting_ReceiveServer interface {
	Send(*RemoteMessage) error
	Recv() (*RemoteMessage, error)
	grpc.ServerStream
}

type remotingReceiveServer struct {
	grpc.ServerStream
}

func (x *remotingReceiveServer) Send(m *RemoteMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *remotingReceiveServer) Recv() (*RemoteMessage, error) {
	m := new(RemoteMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Remoting_ListProcesses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProcessesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemotingServer).ListProcesses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/remote.Remoting/ListProcesses",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemotingServer).ListProcesses(ctx, req.(*ListProcessesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Remoting_GetProcessDiagnostics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProcessDiagnosticsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RemotingServer).GetProcessDiagnostics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/remote.Remoting/GetProcessDiagnostics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RemotingServer).GetProcessDiagnostics(ctx, req.(*GetProcessDiagnosticsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Remoting_ServiceDesc is the grpc.ServiceDesc for Remoting service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Remoting_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "remote.Remoting",
	HandlerType: (*RemotingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListProcesses",
			Handler:    _Remoting_ListProcesses_Handler,
		},
		{
			MethodName: "GetProcessDiagnostics",
			Handler:    _Remoting_GetProcessDiagnostics_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Receive",
			Handler:       _Remoting_Receive_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "remote.proto",
}