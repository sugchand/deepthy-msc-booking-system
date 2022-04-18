package rpc

import (
	proto "bookingSystem/proto/go/userProto"
	"bookingSystem/userAuth/pkg/db"
	"bookingSystem/userAuth/pkg/env"
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCServer struct {
	userTablePtr *db.UserTableHandle
	proto.UnimplementedUserAuthServer
}

func (rs *RPCServer) NewUser(ctx context.Context, userData *proto.UserWithDetails) (*emptypb.Empty, error) {
	uname := userData.GetRequest().GetUsername()
	pwd := userData.GetRequest().GetPassword()
	email := userData.GetDetails().GetEmail()
	isAdmin := userData.GetDetails().GetIsAdmin()
	// TODO :: validate username, pwd and email is right before pushing to the DB
	err := rs.userTablePtr.NewUser(ctx, uname, pwd, email, isAdmin)
	return &emptypb.Empty{}, err
}

func (rs *RPCServer) DelUser(ctx context.Context, userReq *proto.UserRequest) (*emptypb.Empty, error) {
	uname := userReq.GetUsername()
	pwd := userReq.GetPassword()
	err := rs.userTablePtr.DeleteUser(ctx, uname, pwd)
	return &emptypb.Empty{}, err
}

func (rs *RPCServer) GetUserToken(ctx context.Context, userReq *proto.UserRequest) (*proto.UserToken, error) {
	uname := userReq.GetUsername()
	pwd := userReq.GetPassword()
	token, validity, isAdmin, err := rs.userTablePtr.GetUserToken(ctx, uname, pwd)

	return &proto.UserToken{
		Token:    token,
		Validity: durationpb.New(validity),
		IsAdmin:  isAdmin,
	}, err
}

func (rs *RPCServer) UpdateUserPassword(ctx context.Context, details *proto.ResetPwdMessage) (*emptypb.Empty, error) {
	uname := details.GetUsername()
	newPwd := details.GetNewPwd()
	email := details.GetEmail()
	err := rs.userTablePtr.ResetPassword(ctx, uname, email, newPwd)
	return &emptypb.Empty{}, err
}

func NewRPCServer(userTablePtr *db.UserTableHandle, env *env.UserEnvValues) (*RPCServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.GRPCListenPort()))
	if err != nil {
		log.WithFields(log.Fields{
			"port": env.GRPCListenPort(),
		}).WithError(err).Error("failed to create rpc server, failed to listen on port!")
		return nil, err
	}
	serverInstance := grpc.NewServer()
	userAuthServer := &RPCServer{
		userTablePtr: userTablePtr,
	}
	proto.RegisterUserAuthServer(serverInstance, userAuthServer)
	if err := serverInstance.Serve(lis); err != nil {
		log.WithFields(log.Fields{
			"port": env.GRPCListenPort(),
		}).WithError(err).Error("failed to create rpc server, as grpc server creation is failed!")
		return nil, err
	}
	return userAuthServer, nil
}
