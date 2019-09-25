FROM golang:latest

RUN go get golang.org/x/net/websocket && \
    go get github.com/oxequa/realize

WORKDIR /go/


