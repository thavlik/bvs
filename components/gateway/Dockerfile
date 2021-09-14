ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"
FROM ${KINDEST_REPOSITORY}thavlik/cardano-base:${KINDEST_TAG} AS builder
FROM debian:11-slim
COPY --from=builder /root/.local/bin/cardano-cli /usr/local/bin/cardano-cli
