#!/bin/bash -e

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
ROOTPATH=${SCRIPTPATH}/..

source ${SCRIPTPATH}/.env
dbport=$(echo ${DB_REMOTE}|cut -d ":" -f2)
dbremoteIP="host.docker.internal"
dbremote="${dbremoteIP}:${dbport}"
cd ${ROOTPATH}
docker build -t user_auth -f userAuth-Dockerfile .
docker run -p ${USERAUTH_GRPC_LISTEN_PORT}:${USERAUTH_GRPC_LISTEN_PORT} --env-file=tools/.env -e DB_REMOTE=${dbremote} --add-host=${dbremoteIP}:host-gateway user_auth

