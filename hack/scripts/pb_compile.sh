#!/usr/bin/env sh

REPO_ROOT="${REPO_ROOT:-$(cd "$(dirname "$0")/../.." && pwd)}"
pwd
PB_PATH="${REPO_ROOT}/api/v1/pb/watermark/"
PROTO_FILE=${1:-"watermarksvc.proto"}

echo "Generating pb files for ${PROTO_FILE} service"
protoc -I="${PB_PATH}"  "${PB_PATH}/${PROTO_FILE}" --go_out=plugins=grpc:"${PB_PATH}"
