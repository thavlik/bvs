FROM thavlik/otobuild:latest
ENV pkg=api
COPY pkg/definitions/definitions.go .
RUN oto -template $templates/server.go.plush \
    -pkg $pkg \
    -out server.gen.go \
    ./definitions.go \
    && gofmt -w server.gen.go
RUN oto -template $templates/client.go.plush \
    -pkg $pkg \
    -out client.gen.go \
    ./definitions.go \
    && gofmt -w client.gen.go
