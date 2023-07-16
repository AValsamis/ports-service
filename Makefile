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
	REDIS_URL=redis://redis:6379/0 PORTS_JSON_PATH=assets/ports.json ./ports-service.out

docker-build: ## Build the docker image of the application
	docker build -t ports-service -f build/Dockerfile .

docker-run: ## Spin up the application and Redis container using Docker Compose
	docker-compose -f ./build/docker-compose.yml up

docker-down: ## Bring down the application and Redis container
	docker-compose -f ./build/docker-compose.yml down

.PHONY: help lint fmt test build run-local docker-build docker-run docker-up docker-down
