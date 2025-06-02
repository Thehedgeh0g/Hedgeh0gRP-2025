#!/bin/bash
set -e

mkdir -p ./bin

docker build -t protocli-builder .

CONTAINER_ID=$(docker create protocli-builder)

docker cp "$CONTAINER_ID":/app/dist/protocli-linux ./bin/protocli-linux
docker cp "$CONTAINER_ID":/app/dist/protocli.exe ./bin/protocli.exe

docker rm "$CONTAINER_ID"

echo "Бинарники скопированы в ./bin"
