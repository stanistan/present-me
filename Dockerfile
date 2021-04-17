FROM golang:1.16.3-alpine AS build
WORKDIR /go/src/app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go/bin/server -mod readonly ./cmd/server

FROM scratch AS prod
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/bin/server /app
ENTRYPOINT ["/app"]
