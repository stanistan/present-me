COMMIT = $(shell git rev-parse HEAD)
IMAGE = gcr.io/present-me-310221/present-me:$(COMMIT)

.PHONY: watch lint docker-build docker-push

watch:
	docker compose up

lint:
	golangci-lint run -v

docker-build:
	docker build --target prod -t $(IMAGE) .

docker-push: docker-build
	docker push $(IMAGE)
