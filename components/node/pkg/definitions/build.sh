#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
goout=../api
name=node-defs
docker build -t $name -f Dockerfile ../../
if [[ -n $(docker ps | grep $name) ]]; then
    echo "Stopping currently running container \"$name\""
    set +e
    docker stop $name
    docker container wait $name
    docker container rm $name --force
    set -e
fi
docker run -d --name $name --rm $name tail -f /dev/null &>/dev/null
docker exec $name cat server.gen.go > $goout/server.gen.go
docker exec $name cat client.gen.go > $goout/client.gen.go
docker stop $name &>/dev/null
echo "Successfully generated definitions"
