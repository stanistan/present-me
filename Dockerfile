FROM golang:1.16.3-alpine AS build

WORKDIR /go/src/app
ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/server -mod readonly ./cmd/server
ENTRYPOINT ["/go/bin/server"]
