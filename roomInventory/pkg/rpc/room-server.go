package rpc

import (
	proto "bookingSystem/proto/roomInventory"
	"bookingSystem/roomInventory/pkg/db"
	"bookingSystem/roomInventory/pkg/env"
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RPCServer struct {
	handle *db.RoomsTableHandle
	proto.UnimplementedRoomInventoryServer
}

func (rs *RPCServer) NewRoom(ctx context.Context, roomWithUID *proto.RoomWithUserToken) (*emptypb.Empty, error) {
	model := roomWithUID.GetRoom().GetKind().GetRoomModel().String()
	// TODO :: validate user token.
	err := rs.handle.NewRoom(ctx, roomWithUID.GetRoom().GetNumber(), roomWithUID.GetRoom().GetDescription(), model,
		roomWithUID.GetRoom().GetKind().GetDescription())
	return &emptypb.Empty{}, err
}

func (rs *RPCServer) RemoveRoom(ctx context.Context, roomWithToken *proto.RoomNumberWithUserToken) (*emptypb.Empty, error) {
	// TODO :: validate user token.
	err := rs.handle.DeleteRoom(ctx, roomWithToken.GetRoomNumber())
	return &emptypb.Empty{}, err
}

func NewRPCServer(roomTablePtr *db.RoomsTableHandle, env *env.RoomEnvValues) (*RPCServer, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.GRPCListenPort()))
	if err != nil {
		log.WithFields(log.Fields{
			"port": env.GRPCListenPort(),
		}).WithError(err).Error("failed to create rpc server, failed to listen on port!")
		return nil, err
	}
	serverInstance := grpc.NewServer()
	roomServer := &RPCServer{
		handle: roomTablePtr,
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
