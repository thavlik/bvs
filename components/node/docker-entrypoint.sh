#!/bin/bash
set -euo pipefail
echo "Starting cardano-node..."
cardano-node run \
    --config $CARDANO_CONFIG_PATH/testnet-config.json \
    --database-path $CARDANO_DB_PATH \
    --socket-path $CARDANO_DB_PATH/node.socket \
    --host-addr 172.17.0.2 \
    --port 1337 \
    --topology $CARDANO_CONFIG_PATH/testnet-topology.json
