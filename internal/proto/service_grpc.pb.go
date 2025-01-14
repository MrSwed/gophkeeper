// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: service.proto

package proto

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

// DataClient is the client API for Data service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataClient interface {
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	SyncItem(ctx context.Context, in *ItemSync, opts ...grpc.CallOption) (*ItemSync, error)
}

type dataClient struct {
	cc grpc.ClientConnInterface
}

func NewDataClient(cc grpc.ClientConnInterface) DataClient {
	return &dataClient{cc}
}

func (c *dataClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/service.Data/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataClient) SyncItem(ctx context.Context, in *ItemSync, opts ...grpc.CallOption) (*ItemSync, error) {
	out := new(ItemSync)
	err := c.cc.Invoke(ctx, "/service.Data/SyncItem", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataServer is the server API for Data service.
// All implementations must embed UnimplementedDataServer
// for forward compatibility
type DataServer interface {
	List(context.Context, *ListRequest) (*ListResponse, error)
	SyncItem(context.Context, *ItemSync) (*ItemSync, error)
	mustEmbedUnimplementedDataServer()
}

// UnimplementedDataServer must be embedded to have forward compatible implementations.
type UnimplementedDataServer struct {
}

func (UnimplementedDataServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDataServer) SyncItem(context.Context, *ItemSync) (*ItemSync, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncItem not implemented")
}
func (UnimplementedDataServer) mustEmbedUnimplementedDataServer() {}

// UnsafeDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataServer will
// result in compilation errors.
type UnsafeDataServer interface {
	mustEmbedUnimplementedDataServer()
}

func RegisterDataServer(s grpc.ServiceRegistrar, srv DataServer) {
	s.RegisterService(&Data_ServiceDesc, srv)
}

func _Data_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Data/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Data_SyncItem_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ItemSync)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).SyncItem(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Data/SyncItem",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).SyncItem(ctx, req.(*ItemSync))
	}
	return interceptor(ctx, in, info, handler)
}

// Data_ServiceDesc is the grpc.ServiceDesc for Data service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Data_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Data",
	HandlerType: (*DataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Data_List_Handler,
		},
		{
			MethodName: "SyncItem",
			Handler:    _Data_SyncItem_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	RegisterClient(ctx context.Context, in *RegisterClientRequest, opts ...grpc.CallOption) (*ClientToken, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) RegisterClient(ctx context.Context, in *RegisterClientRequest, opts ...grpc.CallOption) (*ClientToken, error) {
	out := new(ClientToken)
	err := c.cc.Invoke(ctx, "/service.Auth/RegisterClient", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility
type AuthServer interface {
	RegisterClient(context.Context, *RegisterClientRequest) (*ClientToken, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (UnimplementedAuthServer) RegisterClient(context.Context, *RegisterClientRequest) (*ClientToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterClient not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_RegisterClient_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).RegisterClient(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.Auth/RegisterClient",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).RegisterClient(ctx, req.(*RegisterClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterClient",
			Handler:    _Auth_RegisterClient_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	SyncUser(ctx context.Context, in *UserSync, opts ...grpc.CallOption) (*UserSync, error)
	DeleteUser(ctx context.Context, in *NoMessage, opts ...grpc.CallOption) (*OkResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) SyncUser(ctx context.Context, in *UserSync, opts ...grpc.CallOption) (*UserSync, error) {
	out := new(UserSync)
	err := c.cc.Invoke(ctx, "/service.User/SyncUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) DeleteUser(ctx context.Context, in *NoMessage, opts ...grpc.CallOption) (*OkResponse, error) {
	out := new(OkResponse)
	err := c.cc.Invoke(ctx, "/service.User/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility
type UserServer interface {
	SyncUser(context.Context, *UserSync) (*UserSync, error)
	DeleteUser(context.Context, *NoMessage) (*OkResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have forward compatible implementations.
type UnimplementedUserServer struct {
}

func (UnimplementedUserServer) SyncUser(context.Context, *UserSync) (*UserSync, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncUser not implemented")
}
func (UnimplementedUserServer) DeleteUser(context.Context, *NoMessage) (*OkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUser not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_SyncUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserSync)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).SyncUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.User/SyncUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).SyncUser(ctx, req.(*UserSync))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_DeleteUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).DeleteUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.User/DeleteUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).DeleteUser(ctx, req.(*NoMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SyncUser",
			Handler:    _User_SyncUser_Handler,
		},
		{
			MethodName: "DeleteUser",
			Handler:    _User_DeleteUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service.proto",
}
