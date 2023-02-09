COMMIT = $(shell git rev-parse HEAD)
IMAGE = gcr.io/present-me-310221/present-me:$(COMMIT)
GOLANGCI_LINT_VERSION = v1.51.1

.PHONY: watch lint docker-build docker-push

watch:
	docker compose up

lint:
	docker run --rm -it \
		-v $(PWD):/app \
		-w /app \
		golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) \
		golangci-lint run -v

docker-build:
	docker build --target prod -t $(IMAGE) .

docker-push: docker-build
	docker push $(IMAGE)
