FROM golang:1.17.3-alpine

WORKDIR /go/src/app
ENV CGO_ENABLED=0

CMD go run -mod readonly github.com/cespare/reflex -s -- \
    go run -mod readonly ./cmd/server
