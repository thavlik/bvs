FROM debian:11-slim as builder
RUN apt-get update \
    && apt-get install -y \
        wget \
        tar \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /usr/local/cardano
RUN wget https://hydra.iohk.io/build/7639935/download/1/cardano-node-1.29.0-linux.tar.gz \
    && tar -xzf cardano-node-1.29.0-linux.tar.gz \
    && rm cardano-node-1.29.0-linux.tar.gz
FROM debian:11-slim
RUN apt-get update \
    && apt-get install -y \
        netbase \
        socat \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cardano /usr/local/cardano
ENV PATH="/usr/local/cardano:$PATH"
