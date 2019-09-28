FROM golang:latest

RUN go get -u github.com/golang/dep/cmd/dep && \
    go get golang.org/x/net/websocket && \
    go get github.com/rakyll/statik && \
    go get github.com/oxequa/realize

WORKDIR /go/


