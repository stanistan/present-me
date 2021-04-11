FROM golang:1.16.3-alpine AS base
WORKDIR /go/src/app
ENV CGO_ENABLED=0

FROM base as dev
CMD go run -mod readonly github.com/cespare/reflex -s -- go run -mod readonly ./cmd/server

FROM base as build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go/bin/server -mod readonly ./cmd/server

FROM scratch AS prod
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/server /app
ENTRYPOINT ["/app"]
