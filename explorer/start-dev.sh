#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
IMAGE=thavlik/bvs-explorer:latest
PORT=8080
docker build -t $IMAGE .
docker kill explorer > /dev/null 2>&1 || true
docker run --rm -d --name explorer -v $(pwd):/mnt -p $PORT:$PORT $IMAGE
echo "Development server is hosted at http://localhost:$PORT"
