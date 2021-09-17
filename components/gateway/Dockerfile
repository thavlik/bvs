ARG BASE_IMAGE="debian:11-slim"
ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG} AS builder
WORKDIR /go/src/github.com/thavlik/bvs
COPY components/gateway components/gateway
WORKDIR /go/src/github.com/thavlik/bvs/components/gateway/cmd/gateway
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM ${BASE_IMAGE}
COPY --from=builder /go/src/github.com/thavlik/bvs/components/gateway/cmd/gateway/gateway /usr/local/bin
CMD ["gateway", "server"]

