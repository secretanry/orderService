BINARY_NAME ?= orders-app
GO_VERSION ?= 1.24
IMAGE_NAME ?= myapp-service

COMPOSE_FILE ?= compose.yml
DEV_COMPOSE_FILE ?= docker-compose.dev.yml

GO_FILES ?= $(wildcard *.go) $(wildcard **/*.go)
GO_PKGS ?= $(shell go list ./... | grep -v /vendor/)

.DEFAULT_GOAL := run


.PHONY: run
run:
	@go run main.go

.PHONY: build
build:
	@CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/$(BINARY_NAME) main.go

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: compose-up
compose-up:
	@docker compose -f $(COMPOSE_FILE) up --build -d

.PHONY: compose-down
compose-down:
	@docker compose -f $(COMPOSE_FILE) down

.PHONY: compose-up-dev
compose-up-dev:
	@docker compose -f $(DEV_COMPOSE_FILE) up --build -d

.PHONY: compose-down-dev
compose-down-dev:
	@docker compose -f $(DEV_COMPOSE_FILE) down

.PHONY: generate-swagger
generate-swagger:
	@swag init
