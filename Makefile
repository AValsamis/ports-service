help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

lint: ## Perform linting
	go vet ./...
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w `find . -name '*.go'`
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
	golangci-lint run

fmt: ## Format the source code
	go fmt ./...

test: ## Run tests locally with race detector and test coverage
	go test ./... -race -cover

build: ## Build the application
	go build -o ports-service.out ./cmd

run-local: ## Run the application locally
	PORTS_JSON_PATH=ports.json ./ports-service.out

docker-build: ## Build the docker image of the application
	docker build -t ports-service .

docker-run: ## Run the application in docker
	docker run -it --rm -e PORTS_JSON_PATH=ports.json ports-service

.PHONY: help lint fmt test-local build run-local docker-build docker-run
