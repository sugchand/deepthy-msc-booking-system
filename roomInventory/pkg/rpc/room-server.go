package rpc

import (
	proto "bookingSystem/proto/roomInventory"
	"bookingSystem/roomInventory/pkg/db"
	"bookingSystem/roomInventory/pkg/env"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var emptyMsg = &emptypb.Empty{}

type RPCServer struct {
	handle            *db.RoomsTableHandle
	userServiceClient *userAuthClient
	proto.UnimplementedRoomInventoryServer
}

func (rs *RPCServer) NewRoom(ctx context.Context, roomWithUID *proto.RoomWithUserToken) (*emptypb.Empty, error) {
	model := roomWithUID.GetRoom().GetKind().GetRoomModel().String()

	token, err := rs.userServiceClient.isTokenValid(ctx, roomWithUID.GetToken())
	if err != nil {
		log.WithFields(log.Fields{
			"token": roomWithUID.GetToken().String(),
		}).WithError(err).Error("failed to validate user token with userauth service")
		return emptyMsg, err
	}
	if token == nil {
		// no valid token.
		return emptyMsg, errors.New("empty token from userauth, cannot proceed!")
	}
	if !token.GetIsAdmin() {
		// the user is not admin, hence cannot proceed with the room creation
		return emptyMsg, errors.New("the user is not admin, cannot create a new room!")
	}
	if token.GetValidity().AsDuration() <= time.Duration(0) {
		return emptyMsg, errors.New("the token is expired, cannot create new room")
	}
	err = rs.handle.NewRoom(ctx, roomWithUID.GetRoom().GetNumber(), roomWithUID.GetRoom().GetDescription(), model)
	return emptyMsg, err
}

func (rs *RPCServer) RemoveRoom(ctx context.Context, roomWithToken *proto.RoomNumberWithUserToken) (*emptypb.Empty, error) {
	token, err := rs.userServiceClient.isTokenValid(ctx, roomWithToken.GetToken())
	if err != nil {
		log.WithFields(log.Fields{
			"token": roomWithToken.GetToken().String(),
		}).WithError(err).Error("failed to validate user token with userauth service")
		return emptyMsg, err
	}
	if token == nil {
		// no valid token.
		return emptyMsg, errors.New("empty token from userauth, cannot proceed!")
	}
	if !token.GetIsAdmin() {
		// the user is not admin, hence cannot proceed with the room creation
		return emptyMsg, errors.New("the user is not admin, cannot delete room entry!")
	}
	if token.GetValidity().AsDuration() <= time.Duration(0) {
		return emptyMsg, errors.New("the token is expired, cannotremove a room")
	}
	err = rs.handle.DeleteRoom(ctx, roomWithToken.GetRoomNumber())
	return emptyMsg, err
}

func NewRPCServer(ctx context.Context, roomTablePtr *db.RoomsTableHandle, env *env.RoomEnvValues) (*RPCServer, error) {
	userAuthHost, userAuthPort := env.UserAuthHostAndPort()
	// lets initialize the userAuth client
	client, err := newUserAuthClient(ctx, userAuthHost, userAuthPort)
	if err != nil {
		log.WithFields(log.Fields{
			"userauth-host": userAuthHost,
			"userauth-port": userAuthPort,
		}).WithError(err).Error("failed to create client to userauth service!")
		return nil, err
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.GRPCListenPort()))
	if err != nil {
		log.WithFields(log.Fields{
			"port": env.GRPCListenPort(),
		}).WithError(err).Error("failed to create rpc server, failed to listen on port!")
		return nil, err
	}
	serverInstance := grpc.NewServer()
	roomServer := &RPCServer{
		handle:            roomTablePtr,
		userServiceClient: client,
	}
	proto.RegisterRoomInventoryServer(serverInstance, roomServer)
	if err := serverInstance.Serve(lis); err != nil {
		log.WithFields(log.Fields{
			"port": env.GRPCListenPort(),
		}).WithError(err).Error("failed to create rpc server, as grpc server creation is failed!")
		return nil, err
	}
	return roomServer, nil
}
