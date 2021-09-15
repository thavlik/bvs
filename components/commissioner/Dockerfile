ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG} AS builder
WORKDIR /usr/local/cardano
RUN wget https://hydra.iohk.io/build/7639935/download/1/cardano-node-1.29.0-linux.tar.gz \
    && tar -xzf cardano-node-1.29.0-linux.tar.gz \
    && rm cardano-node-1.29.0-linux.tar.gz
RUN go install github.com/sumup-oss/gocat@latest
WORKDIR /go/src/github.com/thavlik/bvs
COPY components/commissioner components/commissioner
WORKDIR /go/src/github.com/thavlik/bvs/components/commissioner/cmd/commissioner
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM ${KINDEST_REPOSITORY}thavlik/cardano-base:${KINDEST_TAG}
COPY --from=builder /go/bin/gocat /usr/local/bin
COPY --from=builder /go/src/github.com/thavlik/bvs/components/commissioner/cmd/commissioner/commissioner /usr/local/bin
ENV CARDANO_NODE_SOCKET_PATH="/tmp/node.socket"
ENV TESTNET_MAGIC="1097911063"
CMD ["commissioner", "server"]

