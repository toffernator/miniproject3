// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// RMClient is the client API for RM service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RMClient interface {
	Bid(ctx context.Context, in *BidMsg, opts ...grpc.CallOption) (*Ack, error)
	Result(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Outcome, error)
	ForceBid(ctx context.Context, in *BidMsg, opts ...grpc.CallOption) (*Ack, error)
}

type rMClient struct {
	cc grpc.ClientConnInterface
}

func NewRMClient(cc grpc.ClientConnInterface) RMClient {
	return &rMClient{cc}
}

func (c *rMClient) Bid(ctx context.Context, in *BidMsg, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/RM/Bid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rMClient) Result(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Outcome, error) {
	out := new(Outcome)
	err := c.cc.Invoke(ctx, "/RM/Result", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rMClient) ForceBid(ctx context.Context, in *BidMsg, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/RM/ForceBid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RMServer is the server API for RM service.
// All implementations must embed UnimplementedRMServer
// for forward compatibility
type RMServer interface {
	Bid(context.Context, *BidMsg) (*Ack, error)
	Result(context.Context, *Empty) (*Outcome, error)
	ForceBid(context.Context, *BidMsg) (*Ack, error)
	mustEmbedUnimplementedRMServer()
}

// UnimplementedRMServer must be embedded to have forward compatible implementations.
type UnimplementedRMServer struct {
}

func (UnimplementedRMServer) Bid(context.Context, *BidMsg) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Bid not implemented")
}
func (UnimplementedRMServer) Result(context.Context, *Empty) (*Outcome, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Result not implemented")
}
func (UnimplementedRMServer) ForceBid(context.Context, *BidMsg) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ForceBid not implemented")
}
func (UnimplementedRMServer) mustEmbedUnimplementedRMServer() {}

// UnsafeRMServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RMServer will
// result in compilation errors.
type UnsafeRMServer interface {
	mustEmbedUnimplementedRMServer()
}

func RegisterRMServer(s grpc.ServiceRegistrar, srv RMServer) {
	s.RegisterService(&RM_ServiceDesc, srv)
}

func _RM_Bid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RMServer).Bid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RM/Bid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RMServer).Bid(ctx, req.(*BidMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _RM_Result_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RMServer).Result(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RM/Result",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RMServer).Result(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RM_ForceBid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BidMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RMServer).ForceBid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RM/ForceBid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RMServer).ForceBid(ctx, req.(*BidMsg))
	}
	return interceptor(ctx, in, info, handler)
}

// RM_ServiceDesc is the grpc.ServiceDesc for RM service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RM_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RM",
	HandlerType: (*RMServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Bid",
			Handler:    _RM_Bid_Handler,
		},
		{
			MethodName: "Result",
			Handler:    _RM_Result_Handler,
		},
		{
			MethodName: "ForceBid",
			Handler:    _RM_ForceBid_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auction.proto",
}
