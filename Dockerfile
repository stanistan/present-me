
FROM oven/bun:1.1 as tw
WORKDIR /app
COPY package.json bun.lockb .
RUN bun install --frozen-lockfile
COPY . .
RUN bun run tailwindcss -i ./static/input.css -o styles.css --minify


FROM golang:1.22.0 as app
ENV GOARCH=amd64
ENV GOOS=linux
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG VERSION_SHA
RUN go build -o server -ldflags="-X main.version=$VERSION_SHA" ./cmd/veun

FROM scratch as prod
COPY --from=app /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=app /app/server /app
COPY --from=app /app/static /static
COPY --from=tw /app/styles.css /static/styles.css
ENTRYPOINT ["/app", "serve"]
