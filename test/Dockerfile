#ARG KINDEST_REPOSITORY=""
#ARG KINDEST_TAG="latest"
#FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG}
FROM golang:1.16
ENV GO111MODULE=on
WORKDIR /go/src/github.com/thavlik/bvs
COPY go.mod .
COPY go.sum .
COPY components components
COPY test test
WORKDIR /go/src/github.com/thavlik/bvs/test/suites
CMD ["bash", "-c", "go test ./... -v -timeout 3h"]

