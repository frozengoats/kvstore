GO_IMAGE := golang:1.24-alpine

.PHONY: test
test:
	$(GO_RUN) test

.PHONY: lint-check
lint-check:
	docker run -t --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v2.1.1 golangci-lint run