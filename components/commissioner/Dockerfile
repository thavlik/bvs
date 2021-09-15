ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG} AS builder
COPY components/commissioner components/commissioner
WORKDIR /go/src/github.com/thavlik/bvs/components/commissioner/cmd/commissioner
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM inputoutput/cardano-node:latest
COPY --from=builder /go/src/github.com/thavlik/bvs/components/commissioner/cmd/commissioner/commissioner /usr/local/bin
COPY components/commissioner/configs /configs
ENV CARDANO_CONFIG_PATH="/configs"
ENV CARDANO_DB_PATH="/mnt/db"
CMD ["commissioner", "server"]

