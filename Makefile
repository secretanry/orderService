BINARY_NAME ?= orders-app
GO_VERSION ?= 1.24
IMAGE_NAME ?= myapp-service

COMPOSE_FILE ?= compose.yml
DEV_COMPOSE_FILE ?= docker-compose.dev.yml
MONITORING_COMPOSE_FILE ?= docker-compose.monitoring.yml

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

.PHONY: test-integration
test-integration:
	@chmod +x scripts/run_integration_tests.sh
	@./scripts/run_integration_tests.sh


.PHONY: test-unit
test-unit:
	@echo "Running unit tests..."
	@go test -v ./handlers ./services/... ./modules/...

.PHONY: test
test: test-unit test-integration


.PHONY: test-cleanup
test-cleanup:
	@echo "Cleaning up test environment..."
	@docker-compose -f docker-compose.test.yml down -v
	@docker volume prune -f
	@echo "Test environment cleaned up"

# Monitoring Stack Management
.PHONY: monitoring-env
monitoring-env:
	@echo "Creating monitoring.env from existing environment configuration..."
	@if [ -f .env ]; then \
		echo "Using .env file as source..."; \
		cp .env monitoring.env; \
	elif [ -f test.env ]; then \
		echo "Using test.env file as source..."; \
		cp test.env monitoring.env; \
	else \
		echo "No .env or test.env found, using monitoring.env.example..."; \
		cp monitoring.env.example monitoring.env; \
	fi
	@echo "Adding monitoring configuration to monitoring.env..."
	@if ! grep -q "METRICS_PORT" monitoring.env; then \
		echo "METRICS_PORT=8081" >> monitoring.env; \
	fi
	@if ! grep -q "LOG_LEVEL" monitoring.env; then \
		echo "LOG_LEVEL=info" >> monitoring.env; \
	fi
	@echo "monitoring.env created successfully!"

.PHONY: monitoring-up
monitoring-up: monitoring-env
	@echo "Starting monitoring stack (Prometheus + Grafana)..."
	@docker-compose -f $(MONITORING_COMPOSE_FILE) up -d
	@echo "Monitoring stack started!"
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "Setting up Grafana automatically..."
	@./scripts/setup-grafana.sh
	@echo ""
	@echo "🎉 Monitoring stack is ready!"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Dashboard: WB-L0 Service Dashboard (automatically imported)"

.PHONY: monitoring-down
monitoring-down:
	@echo "Stopping monitoring stack..."
	@docker-compose -f $(MONITORING_COMPOSE_FILE) down
	@echo "Monitoring stack stopped!"

.PHONY: monitoring-logs
monitoring-logs:
	@echo "Showing monitoring stack logs..."
	@docker-compose -f $(MONITORING_COMPOSE_FILE) logs -f

.PHONY: monitoring-status
monitoring-status:
	@echo "Monitoring stack status:"
	@docker-compose -f $(MONITORING_COMPOSE_FILE) ps

.PHONY: monitoring-restart
monitoring-restart: monitoring-down monitoring-up

.PHONY: monitoring-cleanup
monitoring-cleanup:
	@echo "Cleaning up monitoring stack and volumes..."
	@docker-compose -f $(MONITORING_COMPOSE_FILE) down -v
	@docker volume prune -f
	@echo "Monitoring stack cleaned up!"

.PHONY: setup-grafana
setup-grafana:
	@echo "Setting up Grafana with automatic configuration..."
	@./scripts/setup-grafana.sh

# Combined commands for full application with monitoring
.PHONY: run-with-monitoring
run-with-monitoring: monitoring-up
	@echo "Starting application with monitoring..."
	@echo "Make sure to set your environment variables or use monitoring.env"
	@go run main.go

.PHONY: dev-with-monitoring
dev-with-monitoring: monitoring-up compose-up-dev
	@echo "Development environment with monitoring started!"
	@echo "Application: http://localhost:8080"
	@echo "Prometheus: http://localhost:9090"
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Health check: http://localhost:8080/health"
	@echo "Metrics: http://localhost:8081/metrics"

.PHONY: stop-all
stop-all: monitoring-down compose-down-dev
	@echo "All services stopped!"
