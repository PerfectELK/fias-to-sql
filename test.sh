#!/bin/bash

set -e

SCRIPT=$(readlink -f "$0")
SCRIPT_PATH=$(dirname "$SCRIPT")
cd $SCRIPT_PATH/

export $(grep -v '^#' .env | xargs -d '\r')
export APP_ROOT=$(pwd)
go test -v -count=1 ./internal/services/fias