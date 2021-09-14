FROM thavlik/cardano-base:latest AS builder
FROM debian:11-slim
COPY --from=builder /root/.local/bin/cardano-cli /usr/local/bin/cardano-cli
