ARG KINDEST_REPOSITORY=""
ARG KINDEST_TAG="latest"

FROM ${KINDEST_REPOSITORY}thavlik/go-base:${KINDEST_TAG}
COPY . .
WORKDIR /go/src/github.com/thavlik/bvs/test/suites
CMD ["bash", "-c", "go test ./... -v -timeout 3h"]
