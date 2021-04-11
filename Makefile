COMMIT = $(shell git rev-parse HEAD)
IMAGE = gcr.io/present-me-310221/present-me:$(COMMIT)

.PHONY: watch lint docker-build docker-push

watch:
	docker compose up

lint:
	docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.39.0 golangci-lint run -v

docker-build:
	docker build --target prod -t $(IMAGE) .

docker-push: docker-build
	docker push $(IMAGE)
