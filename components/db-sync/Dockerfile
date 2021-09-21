ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG} AS builder
WORKDIR /go/src/github.com/thavlik/bvs
COPY components/node/pkg/api components/node/pkg/api
COPY components/db-sync components/db-sync
RUN cd components/db-sync/cmd/db-sync \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM inputoutput/cardano-db-sync:master
RUN mkdir /tmp
WORKDIR /etc/db-sync
COPY components/db-sync/schema schema
RUN mkdir state
COPY components/node/configs /configs
COPY components/db-sync/configs /etc/db-sync/configs
RUN cd /etc/db-sync/configs \
    && chmod 600 pgpass-testnet \
    && chmod 600 pgpass-mainnet
COPY --from=builder /go/src/github.com/thavlik/bvs/components/db-sync/cmd/db-sync/db-sync /bin
ENV PGPASSFILE=/etc/db-sync/configs/pgpass-testnet
CMD ["db-sync"]
