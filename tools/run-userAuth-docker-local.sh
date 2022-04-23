#!/bin/bash -e

SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
ROOTPATH=${SCRIPTPATH}/..

source ${SCRIPTPATH}/.env
cd ${ROOTPATH}
docker build -t user_auth -f userAuth-Dockerfile .
docker run  --net=host --env-file=tools/.env user_auth

