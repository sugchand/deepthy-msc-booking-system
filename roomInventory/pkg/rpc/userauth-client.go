package rpc

import (
	proto "bookingSystem/proto/userProto"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type userAuthClient struct {
	client proto.UserAuthClient
}

func (uac *userAuthClient) isTokenValid(ctx context.Context, tokenIn *proto.UserToken) (*proto.UserToken, error) {
	return uac.client.IsTokenValid(ctx, tokenIn)
}

func newUserAuthClient(ctx context.Context, userAuthRemoteHost string, userAuthRemotePort uint32) (*userAuthClient, error) {
	remoteAddr := fmt.Sprintf("%s:%d", userAuthRemoteHost, userAuthRemotePort)
	conn, err := grpc.Dial(remoteAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.WithFields(log.Fields{
			"userauth-host": userAuthRemoteHost,
			"userauth-port": userAuthRemotePort,
		}).WithError(err).Error("failed to connect to the userauth service!")
		return nil, err
	}
	return &userAuthClient{
		client: proto.NewUserAuthClient(conn),
	}, nil
}
