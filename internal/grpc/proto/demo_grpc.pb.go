// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.1
// source: proto/demo.proto

package proto

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
	URLs_AddURL_FullMethodName   = "/demo.URLs/AddURL"
	URLs_ListURLs_FullMethodName = "/demo.URLs/ListURLs"
	URLs_GetURL_FullMethodName   = "/demo.URLs/GetURL"
	URLs_DelURL_FullMethodName   = "/demo.URLs/DelURL"
)

// URLsClient is the client API for URLs service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLsClient interface {
	AddURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error)
	ListURLs(ctx context.Context, in *ListURLsRequest, opts ...grpc.CallOption) (*ListURLsResponse, error)
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	DelURL(ctx context.Context, in *DelURLRequest, opts ...grpc.CallOption) (*DelURLResponse, error)
}

type uRLsClient struct {
	cc grpc.ClientConnInterface
}

func NewURLsClient(cc grpc.ClientConnInterface) URLsClient {
	return &uRLsClient{cc}
}

func (c *uRLsClient) AddURL(ctx context.Context, in *AddURLRequest, opts ...grpc.CallOption) (*AddURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddURLResponse)
	err := c.cc.Invoke(ctx, URLs_AddURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLsClient) ListURLs(ctx context.Context, in *ListURLsRequest, opts ...grpc.CallOption) (*ListURLsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListURLsResponse)
	err := c.cc.Invoke(ctx, URLs_ListURLs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLsClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, URLs_GetURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLsClient) DelURL(ctx context.Context, in *DelURLRequest, opts ...grpc.CallOption) (*DelURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DelURLResponse)
	err := c.cc.Invoke(ctx, URLs_DelURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLsServer is the server API for URLs service.
// All implementations must embed UnimplementedURLsServer
// for forward compatibility.
type URLsServer interface {
	AddURL(context.Context, *AddURLRequest) (*AddURLResponse, error)
	ListURLs(context.Context, *ListURLsRequest) (*ListURLsResponse, error)
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	DelURL(context.Context, *DelURLRequest) (*DelURLResponse, error)
	mustEmbedUnimplementedURLsServer()
}

// UnimplementedURLsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedURLsServer struct{}

func (UnimplementedURLsServer) AddURL(context.Context, *AddURLRequest) (*AddURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddURL not implemented")
}
func (UnimplementedURLsServer) ListURLs(context.Context, *ListURLsRequest) (*ListURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListURLs not implemented")
}
func (UnimplementedURLsServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}
func (UnimplementedURLsServer) DelURL(context.Context, *DelURLRequest) (*DelURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelURL not implemented")
}
func (UnimplementedURLsServer) mustEmbedUnimplementedURLsServer() {}
func (UnimplementedURLsServer) testEmbeddedByValue()              {}

// UnsafeURLsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLsServer will
// result in compilation errors.
type UnsafeURLsServer interface {
	mustEmbedUnimplementedURLsServer()
}

func RegisterURLsServer(s grpc.ServiceRegistrar, srv URLsServer) {
	// If the following call pancis, it indicates UnimplementedURLsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&URLs_ServiceDesc, srv)
}

func _URLs_AddURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).AddURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLs_AddURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).AddURL(ctx, req.(*AddURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLs_ListURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).ListURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLs_ListURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).ListURLs(ctx, req.(*ListURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLs_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLs_GetURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLs_DelURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLsServer).DelURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLs_DelURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLsServer).DelURL(ctx, req.(*DelURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLs_ServiceDesc is the grpc.ServiceDesc for URLs service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLs_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "demo.URLs",
	HandlerType: (*URLsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddURL",
			Handler:    _URLs_AddURL_Handler,
		},
		{
			MethodName: "ListURLs",
			Handler:    _URLs_ListURLs_Handler,
		},
		{
			MethodName: "GetURL",
			Handler:    _URLs_GetURL_Handler,
		},
		{
			MethodName: "DelURL",
			Handler:    _URLs_DelURL_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/demo.proto",
}
