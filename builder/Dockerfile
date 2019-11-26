FROM golang:latest

WORKDIR /go/

RUN go get -u github.com/rakyll/statik && \
    apt update && \
    apt install -y npm

