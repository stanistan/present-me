.PHONY: watch lint

watch:
	go run -mod readonly github.com/cespare/reflex -s -v -- go run -mod readonly ./cmd/server config.yaml

lint:
	docker run --rm -it -v $(PWD):/app -w /app golangci/golangci-lint:v1.39.0 golangci-lint run -v
