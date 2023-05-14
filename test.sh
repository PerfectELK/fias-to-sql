#!/bin/bash

set -e

SCRIPT=$(readlink -f "$0")
SCRIPT_PATH=$(dirname "$SCRIPT")
cd $SCRIPT_PATH/

export $(grep -v '^#' .env | xargs -d '\r') && go test -v -count=1 ./pkg/db/pgsql