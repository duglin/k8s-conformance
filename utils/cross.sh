#!/bin/bash

EXE=${EXE:-bin/kubecon}
BUILD_OPTS=${BUILD_OPTS:-}

# List of platforms we build binaries for at this time:
PLATFORMS="darwin/amd64 windows/amd64 linux/amd64" # OSX, Windows, Linux x86_64
PLATFORMS="$PLATFORMS linux/ppc64le linux/s390x"   # IBM POWER and z Systems
PLATFORMS="$PLATFORMS linux/arm linux/arm64"       # ARM; 32bit and 64bit

for PLATFORM in $PLATFORMS; do
  GOOS=${PLATFORM%/*}
  GOARCH=${PLATFORM#*/}
  BIN_FILENAME="${EXE}-${GOOS}-${GOARCH}"
  if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
  CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILENAME} ${BUILD_OPTS}"
  # echo "${CMD}"
  echo "Building ${BIN_FILENAME}"
  eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done

# eval errors
if [[ "${FAILURES}" != "" ]]; then
  echo ""
  echo "${EXE} build failed on: ${FAILURES}"
  exit 1
fi
