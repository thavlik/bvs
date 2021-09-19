#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
IMAGE=node:16.9.1-stretch-slim
PORT=8080
docker kill explorer > /dev/null 2>&1 || true
docker run --rm -d --name explorer -v $(pwd):/mnt -p $PORT:$PORT $IMAGE bash -c cd /mnt && npm run start-dev
echo "Development server is hosted at http://localhost:$PORT"
