ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG} AS builder
WORKDIR /go/src/github.com/thavlik/bvs
COPY components/node components/node
RUN cd components/node/cmd/node \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM ${KINDEST_REPOSITORY}thavlik/cardano-base:${KINDEST_TAG}
COPY --from=builder /go/src/github.com/thavlik/bvs/components/node/cmd/node/node /usr/local/bin
COPY components/node/configs /configs
ENV CARDANO_CONFIG_PATH="/configs"
ENV CARDANO_DB_PATH="/mnt/db"
ENV CARDANO_NODE_SOCKET_PATH="/mnt/db/node.socket"
ENV TESTNET_MAGIC="1097911063"
CMD ["node", "server"]

