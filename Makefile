VERSION?=$(shell git describe --tags --always)
PROJECT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
PROJECT_NAME=crypto-calculator
CONTAINER_REGISTRY?=localhost:5000/

.PHONY: build
build:
	mkdir -p $(PROJECT_DIR)/bin
	go build -mod vendor -o $(PROJECT_DIR)/bin/ $(PROJECT_DIR)/cmd/app

.PHONY: mod
mod:
	go mod vendor && go mod tidy

.PHONY: deploy-local
deploy-local:
	cd $(PROJECT_DIR) && go mod vendor
	cd $(PROJECT_DIR)/deployment && docker-compose -f docker-compose.yaml up --build -d

.PHONY: lint
lint:
	cd $(PROJECT_DIR) && golangci-lint run ./...

.PHONY: swagger
swagger:
	swag init -d cmd/app,api/handlers
