#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
IMAGE=node:16.9.1-stretch-slim
PORT=8080
NAME=explorer
docker kill $NAME >/dev/null 2>&1 || true
docker run --rm \
    --name $NAME \
    -it \
    -v $(pwd):/mnt \
    -p $PORT:$PORT \
    $IMAGE \
    bash -c cd /mnt \
        && /usr/local/bin/npm i \
        && /usr/local/bin/npm run start-dev
echo "Development server is hosted at http://localhost:$PORT"
docker logs -f $NAME
