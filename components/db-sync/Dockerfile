ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG} AS builder
WORKDIR /go/src/github.com/thavlik/bvs
COPY components/db-sync components/db-sync
RUN cd components/db-sync/cmd/db-sync \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM inputoutput/cardano-db-sync:master
RUN mkdir /tmp
WORKDIR /etc/db-sync
COPY components/db-sync/schema schema
RUN mkdir state
COPY components/node/configs /configs/node
COPY components/db-sync/configs /configs/db-sync
COPY --from=builder /go/src/github.com/thavlik/bvs/components/db-sync/cmd/db-sync/db-sync /bin
ENV PGPASSFILE=/configs/db-sync/pgpass-testnet
CMD ["db-sync"]
