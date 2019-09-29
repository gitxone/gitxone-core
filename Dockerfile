FROM golang:latest

WORKDIR /go/

RUN go get -u github.com/oxequa/realize \
 && go get -u github.com/golang/dep/cmd/dep
