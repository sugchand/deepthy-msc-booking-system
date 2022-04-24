package db

import (
	"bookingSystem/roomInventory/pkg/db/impl"
	"bookingSystem/roomInventory/pkg/env"
	"context"
	"errors"
)

type RoomsTableHandle struct {
	pgDBHandle *impl.PostgresDB
	envValues  *env.RoomEnvValues
}

func (rth *RoomsTableHandle) NewRoom(ctx context.Context, roomNum uint64, roomDesc, roomType string) error {

	return impl.NewRoom(ctx, rth.pgDBHandle, roomNum, roomDesc, roomType)
}

func (rth *RoomsTableHandle) DeleteRoom(ctx context.Context, roomNum uint64) error {
	return impl.DeleteRoom(ctx, rth.pgDBHandle, roomNum)
}

func NewRoomsTableHandle(ctx context.Context, env *env.RoomEnvValues) (*RoomsTableHandle, error) {
	pgHandle := impl.NewPostgresDB(env)
	pgHandle.DB(ctx) // initialize the DB Handle
	err := impl.CreateInventoryTables(ctx, pgHandle)
	if err != nil {
		// so failed to create a table and cant procced without the user table.
		return nil, errors.New("failed to create rooms-inventory tables in postgres!.")
	}
	return &RoomsTableHandle{
		pgDBHandle: pgHandle,
		envValues:  env,
	}, nil
}
