SHELL=/bin/bash

.PHONY: all build clean tidy dev test generate lint

tidy:
	@echo "Tidying up..."
	@go mod tidy

dev:
	@echo "starting dev..."
	@PORT=8081 air

test:
	@echo "Running tests..."
	@go test ./...

generate:
	@echo "Generating code..."
	@go generate ./...

lint:
	@echo "Linting ..."
	@golangci-lint run --timeout 5m
