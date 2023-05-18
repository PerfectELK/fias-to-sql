#!/bin/bash

set -e

SCRIPT=$(readlink -f "$0")
SCRIPT_PATH=$(dirname "$SCRIPT")
cd $SCRIPT_PATH/

cp ../.env ./.env

docker-compose up --build