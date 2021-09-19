#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
IMAGE=thavlik/bvs-explorer:latest
docker build -t $IMAGE .
docker run --rm -d --name explorer -v $(pwd):/mnt -p 8080:8080 $IMAGE
