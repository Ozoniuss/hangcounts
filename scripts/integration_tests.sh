#!/bin/bash

cleanup() {
  echo "Cleaning up Docker Compose..."
  docker compose -f docker-compose.integration.yaml down
}

trap cleanup EXIT

docker compose -f docker-compose.integration.yaml up -d
if [ $? -ne 0 ]; then
  echo "Failed to start Docker Compose. Exiting..."
  exit 1
fi

env HANGCOUNTS_RUN_INTEGRATION_TESTS=true go test -v ./...
exit $?
