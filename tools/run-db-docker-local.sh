#!/bin/bash -e
SCRIPTPATH="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
source ${SCRIPTPATH}/.env

dbport=$(echo ${DB_REMOTE}|cut -d ":" -f2)
docker rm -f postgresql
echo "postgresql is listening on :5432"
docker run --name postgresql --net=host -e POSTGRES_USER=${DB_UNAME} -e POSTGRES_PASSWORD=${DB_PWD} -e POSTGRES_DB=${DB_NAME} -d postgres