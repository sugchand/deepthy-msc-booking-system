package main

import (
	"bookingSystem/roomInventory/pkg/db"
	"bookingSystem/roomInventory/pkg/env"
	"bookingSystem/roomInventory/pkg/rpc"
	"context"
	"fmt"
	"os"
	"os/signal"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

func run(ctx context.Context, envValues *env.RoomEnvValues, args []string) error {
	roomTableHandle, err := db.NewRoomsTableHandle(ctx, envValues)
	if err != nil {
		os.Exit(exitCodeErr)
	}

	// finally initialize the rpc server
	_, err = rpc.NewRPCServer(roomTableHandle, envValues)
	if err != nil {
		os.Exit(exitCodeErr)
	}
	return nil
}

func main() {
	// initialize the env variables first.
	envValues := env.NewRoomInventoryEnv()
	if envValues == nil {
		os.Exit(0)
	}
	// lets initialize the DB backend

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(exitCodeInterrupt)
	}()
	if err := run(ctx, envValues, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitCodeErr)
	}
}
