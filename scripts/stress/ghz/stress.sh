#!/bin/bash
PROTO_PATH=$1
SERVICE_PORT=$2
COUNT=$3
RPS=$4
DATA=$5
METHOD=$6

ghz --insecure --proto "$PROTO_PATH" \
    --import-paths "./proto" -c "$COUNT" \
    --rps "$RPS" -d "$DATA" "localhost:$SERVICE_PORT" \
    --call "$METHOD"