// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.0
// source: proto/userProto/user.proto

package userProto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// UserAuthClient is the client API for UserAuth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserAuthClient interface {
	NewUser(ctx context.Context, in *UserWithDetails, opts ...grpc.CallOption) (*emptypb.Empty, error)
	DelUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetUserToken(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserToken, error)
	// function to check if the current token is valid , the api will respond
	// the token with the session time.
	IsTokenValid(ctx context.Context, in *UserToken, opts ...grpc.CallOption) (*UserToken, error)
	RequestUserPasswordUpdate(ctx context.Context, in *ResetPwdRequest, opts ...grpc.CallOption) (*ResetToken, error)
	UserPasswordUpdate(ctx context.Context, in *ResetPwdMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type userAuthClient struct {
	cc grpc.ClientConnInterface
}

func NewUserAuthClient(cc grpc.ClientConnInterface) UserAuthClient {
	return &userAuthClient{cc}
}

func (c *userAuthClient) NewUser(ctx context.Context, in *UserWithDetails, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bookingSystem.proto.userProto.UserAuth/NewUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAuthClient) DelUser(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bookingSystem.proto.userProto.UserAuth/DelUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAuthClient) GetUserToken(ctx context.Context, in *UserRequest, opts ...grpc.CallOption) (*UserToken, error) {
	out := new(UserToken)
	err := c.cc.Invoke(ctx, "/bookingSystem.proto.userProto.UserAuth/GetUserToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAuthClient) IsTokenValid(ctx context.Context, in *UserToken, opts ...grpc.CallOption) (*UserToken, error) {
	out := new(UserToken)
	err := c.cc.Invoke(ctx, "/bookingSystem.proto.userProto.UserAuth/IsTokenValid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAuthClient) RequestUserPasswordUpdate(ctx context.Context, in *ResetPwdRequest, opts ...grpc.CallOption) (*ResetToken, error) {
	out := new(ResetToken)
	err := c.cc.Invoke(ctx, "/bookingSystem.proto.userProto.UserAuth/RequestUserPasswordUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userAuthClient) UserPasswordUpdate(ctx context.Context, in *ResetPwdMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/bookingSystem.proto.userProto.UserAuth/UserPasswordUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserAuthServer is the server API for UserAuth service.
// All implementations must embed UnimplementedUserAuthServer
// for forward compatibility
type UserAuthServer interface {
	NewUser(context.Context, *UserWithDetails) (*emptypb.Empty, error)
	DelUser(context.Context, *UserRequest) (*emptypb.Empty, error)
	GetUserToken(context.Context, *UserRequest) (*UserToken, error)
	// function to check if the current token is valid , the api will respond
	// the token with the session time.
	IsTokenValid(context.Context, *UserToken) (*UserToken, error)
	RequestUserPasswordUpdate(context.Context, *ResetPwdRequest) (*ResetToken, error)
	UserPasswordUpdate(context.Context, *ResetPwdMessage) (*emptypb.Empty, error)
	mustEmbedUnimplementedUserAuthServer()
}

// UnimplementedUserAuthServer must be embedded to have forward compatible implementations.
type UnimplementedUserAuthServer struct {
}

func (UnimplementedUserAuthServer) NewUser(context.Context, *UserWithDetails) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewUser not implemented")
}
func (UnimplementedUserAuthServer) DelUser(context.Context, *UserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelUser not implemented")
}
func (UnimplementedUserAuthServer) GetUserToken(context.Context, *UserRequest) (*UserToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserToken not implemented")
}
func (UnimplementedUserAuthServer) IsTokenValid(context.Context, *UserToken) (*UserToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsTokenValid not implemented")
}
func (UnimplementedUserAuthServer) RequestUserPasswordUpdate(context.Context, *ResetPwdRequest) (*ResetToken, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestUserPasswordUpdate not implemented")
}
func (UnimplementedUserAuthServer) UserPasswordUpdate(context.Context, *ResetPwdMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserPasswordUpdate not implemented")
}
func (UnimplementedUserAuthServer) mustEmbedUnimplementedUserAuthServer() {}

// UnsafeUserAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserAuthServer will
// result in compilation errors.
type UnsafeUserAuthServer interface {
	mustEmbedUnimplementedUserAuthServer()
}

func RegisterUserAuthServer(s grpc.ServiceRegistrar, srv UserAuthServer) {
	s.RegisterService(&UserAuth_ServiceDesc, srv)
}

func _UserAuth_NewUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserWithDetails)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAuthServer).NewUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bookingSystem.proto.userProto.UserAuth/NewUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAuthServer).NewUser(ctx, req.(*UserWithDetails))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAuth_DelUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAuthServer).DelUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bookingSystem.proto.userProto.UserAuth/DelUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAuthServer).DelUser(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAuth_GetUserToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAuthServer).GetUserToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bookingSystem.proto.userProto.UserAuth/GetUserToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAuthServer).GetUserToken(ctx, req.(*UserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAuth_IsTokenValid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAuthServer).IsTokenValid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bookingSystem.proto.userProto.UserAuth/IsTokenValid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAuthServer).IsTokenValid(ctx, req.(*UserToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAuth_RequestUserPasswordUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetPwdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAuthServer).RequestUserPasswordUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bookingSystem.proto.userProto.UserAuth/RequestUserPasswordUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAuthServer).RequestUserPasswordUpdate(ctx, req.(*ResetPwdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserAuth_UserPasswordUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetPwdMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserAuthServer).UserPasswordUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bookingSystem.proto.userProto.UserAuth/UserPasswordUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserAuthServer).UserPasswordUpdate(ctx, req.(*ResetPwdMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// UserAuth_ServiceDesc is the grpc.ServiceDesc for UserAuth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserAuth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bookingSystem.proto.userProto.UserAuth",
	HandlerType: (*UserAuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewUser",
			Handler:    _UserAuth_NewUser_Handler,
		},
		{
			MethodName: "DelUser",
			Handler:    _UserAuth_DelUser_Handler,
		},
		{
			MethodName: "GetUserToken",
			Handler:    _UserAuth_GetUserToken_Handler,
		},
		{
			MethodName: "IsTokenValid",
			Handler:    _UserAuth_IsTokenValid_Handler,
		},
		{
			MethodName: "RequestUserPasswordUpdate",
			Handler:    _UserAuth_RequestUserPasswordUpdate_Handler,
		},
		{
			MethodName: "UserPasswordUpdate",
			Handler:    _UserAuth_UserPasswordUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/userProto/user.proto",
}
