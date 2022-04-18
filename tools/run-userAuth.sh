#!/bin/bash -e

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
ROOTPATH=${SCRIPTPATH}/..
BINPATH=${ROOTPATH}/bin

export DB_REMOTE="localhost:2022"
export DB_NAME="booking"
export USERAUTH_GRPC_LISTEN_PORT=7070
export USERAUTH_TOKEN_VALIDITY_SECONDS=1800

cd ${ROOTPATH}
make build

echo "running the user auth service to listen on ${USERAUTH_GRPC_LISTEN_PORT}"
./bin/user-auth

