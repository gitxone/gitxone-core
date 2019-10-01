FROM golang:latest

WORKDIR /go/

RUN go get -u github.com/golang/dep/cmd/dep && \
    go get -u github.com/rakyll/statik && \
    apt update && \
    apt install -y npm

