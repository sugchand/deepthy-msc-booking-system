#!/bin/bash

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
export PATH="$PATH:$(go env GOPATH)/bin"

# check if the binary exists before compiling.
if ! command -v protoc &> /dev/null
then
    echo "protoc could not be found, install it using 'https://grpc.io/docs/languages/go/quickstart/'"
    exit 1
fi

cd ${SCRIPTPATH}/..

# lets compile the proto files.
protoc -I. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/userProto/user.proto proto/roomInventory/rooms.proto