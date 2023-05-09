FROM node:19.6-alpine as frontend
RUN corepack enable

WORKDIR /app
COPY frontend/.yarn .yarn/
COPY frontend/package.json frontend/yarn.lock frontend/.yarnrc.yml ./
RUN yarn 
COPY frontend /app
RUN yarn run generate

FROM golang:1.20-alpine as backend
WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server /app
RUN go build -o server ./cmd/server

FROM scratch as prod
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=backend /app/server /app
COPY --from=frontend /app/.output/public /static
ENTRYPOINT ["/app"]
