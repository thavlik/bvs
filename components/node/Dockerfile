FROM thavlik/cardano-base:latest AS builder
FROM debian:11-slim
COPY --from=builder /root/.local/bin/cardano-node /usr/local/bin/cardano-node
COPY configs /configs
COPY docker-entrypoint.sh /usr/local/bin
ENV CARDANO_CONFIG_PATH="/configs"
ENV CARDANO_DB_PATH="/mnt/shared/cardano-data/db"
ENV CARDANO_NODE_SOCKET_PATH="/mnt/shared/cardano-data/db/node.socket"
CMD /bin/bash

