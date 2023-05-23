FROM node:19.6-alpine as frontend
RUN corepack enable

WORKDIR /app
COPY frontend/.yarn .yarn/
COPY frontend/package.json frontend/yarn.lock frontend/.yarnrc.yml ./
RUN yarn 
COPY frontend /app

ARG VERSION_SHA
RUN echo "{ \"rev\": \"$VERSION_SHA\" }" > /app/version.json
RUN yarn run generate

FROM golang:1.20-alpine as server 
WORKDIR /app
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server /app

ARG VERSION_SHA
RUN go build -o server -ldflags="-X main.version=$VERSION_SHA" ./cmd/server

FROM scratch as prod
COPY --from=server /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=server /app/server /app
COPY --from=frontend /app/.output/public /static
ENTRYPOINT ["/app"]
