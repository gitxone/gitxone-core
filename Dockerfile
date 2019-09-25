FROM golang:latest

RUN useradd gitxone --create-home --shell /bin/bash && \
    go get golang.org/x/net/websocket

WORKDIR /go/

USER gitxone

