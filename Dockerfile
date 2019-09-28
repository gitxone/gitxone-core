FROM golang:latest

RUN go get golang.org/x/net/websocket && \
    go get github.com/oxequa/realize \
    go get github.com/rakyll/statik

WORKDIR /go/


