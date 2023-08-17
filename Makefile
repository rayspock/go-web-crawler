.PHONY: test
test:
	@echo "Running tests..."
	@go test -cover -race ./...

.PHONY: generate
generate:
	@echo "Generating mocks..."
	@go generate ./...
