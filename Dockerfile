FROM node:19.6-alpine as frontend

WORKDIR /app
COPY frontend/package.json frontend/yarn.lock ./
RUN yarn install
COPY frontend /app
RUN yarn run generate

FROM golang:1.20-alpine as backend
WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server /app
RUN go build -o server ./cmd/server-nuxt

FROM scratch as prod
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend /app/server /app
COPY --from=frontend /app/.output/public /static
ENTRYPOINT ["/app"]
