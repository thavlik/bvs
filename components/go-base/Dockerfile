FROM golang:1.16
ENV GO111MODULE=on
WORKDIR /go/src/github.com/thavlik/bvs
COPY go.mod .
COPY go.sum .
RUN go mod download
