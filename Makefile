.PHONY: lint test test-cov

lint:
	@golangci-lint run --exclude-use-default=false --enable-all \
		--disable exhaustivestruct

test:
	@go test -race ./...

test-cov:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func coverage.out
