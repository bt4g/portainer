#!/usr/bin/env bash
set -euo pipefail


PLATFORM=${1:-"linux"}
ARCH=${2:-"amd64"}

BINARY_VERSION_FILE="./binary-version.json"

dockerVersion=$(jq -r '.docker' < "${BINARY_VERSION_FILE}")
helmVersion=$(jq -r '.helm' < "${BINARY_VERSION_FILE}")
kubectlVersion=$(jq -r '.kubectl' < "${BINARY_VERSION_FILE}")
mingitVersion=$(jq -r '.mingit' < "${BINARY_VERSION_FILE}")

mkdir -p dist

echo "Downloading binaries for docker ${dockerVersion}, helm ${helmVersion}, kubectl ${kubectlVersion} and mingit ${mingitVersion}"

./build/download_docker_binary.sh "$PLATFORM" "$ARCH" "$dockerVersion" &
./build/download_helm_binary.sh "$PLATFORM" "$ARCH" "$helmVersion" &
./build/download_kubectl_binary.sh "$PLATFORM" "$ARCH" "$kubectlVersion" &
./build/download_mingit_binary.sh "$PLATFORM" "$ARCH" "$mingitVersion" &
wait
