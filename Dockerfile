ARG VERSION_SHA

FROM golang:1.22.0-alpine as server
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o server -ldflags="-X main.version=$VERSION_SHA" ./cmd/veun

FROM scratch as prod
COPY --from=server /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server /app/server /app
COPY --from=frontend /app/.output/public /static
ENTRYPOINT ["/app"]
