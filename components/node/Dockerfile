ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
FROM ${KINDEST_REPOSITORY}thavlik/cardano-base:${KINDEST_TAG} AS builder
#FROM debian:11-slim
#RUN apt update && apt install -y libsodium23
#COPY --from=builder /root/.local/bin/cardano-node /root/.local/bin/cardano-node
#ENV PATH="~/.local/bin:$PATH"
COPY configs /configs
COPY docker-entrypoint.sh /usr/local/bin
RUN chmod +x /usr/local/bin/docker-entrypoint.sh
ENV CARDANO_CONFIG_PATH="/configs"
ENV CARDANO_DB_PATH="/db"
CMD ["docker-entrypoint.sh"]

