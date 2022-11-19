#!/bin/bash

set -e

defaultImage="hello-world"

image="${1:-$defaultImage}"
container=$(docker create "$image")

docker export "$container" -o "./assets/${image}.tar.gz" > /dev/null
docker rm "$container" > /dev/null

docker inspect -f '{{.Config.Cmd}}' "$image:latest" | tr -d '[]\n' > "./assets/${image}-cmd"

echo "Image content stored in assets/${image}.tar.gz"