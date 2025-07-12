#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SKIP_KAFKA=false
while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-kafka)
            SKIP_KAFKA=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--skip-kafka]"
            exit 1
            ;;
    esac
done

print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
}

check_docker_compose() {
    if ! command -v docker-compose > /dev/null 2>&1; then
        print_error "Docker Compose is not installed. Please install Docker Compose and try again."
        exit 1
    fi
}

wait_for_service() {
    local service_name=$1
    local max_attempts=${2:-30}
    local attempt=1
    
    print_status "Waiting for $service_name to be healthy..."
    
    while [ $attempt -le $max_attempts ]; do
        if docker-compose -f docker-compose.test.yml ps $service_name | grep -q "healthy"; then
            print_status "$service_name is healthy!"
            return 0
        fi
        
        print_status "Attempt $attempt/$max_attempts: $service_name is not ready yet..."
        sleep 2
        attempt=$((attempt + 1))
    done
    
    print_error "$service_name failed to become healthy after $max_attempts attempts"
    return 1
}

cleanup() {
    print_status "Cleaning up test services..."
    docker-compose -f docker-compose.test.yml down -v
    print_status "Cleanup completed"
}

trap cleanup EXIT

main() {
    print_status "Starting integration tests..."
    
    check_docker
    check_docker_compose
    
    export COMPOSE_PROJECT_NAME=wb-l0

    echo "[INFO] Stopping any existing test services..."
    docker-compose -f docker-compose.test.yml down -v

    if docker volume ls | grep -q wb-l0_test_kafka_data; then
      echo "[INFO] Removing old Kafka data volume..."
      docker volume rm wb-l0_test_kafka_data || true
    fi

    echo "[INFO] Starting test services..."
    docker-compose -f docker-compose.test.yml up -d --build
    
    wait_for_service test-postgres
    wait_for_service test-redis
    wait_for_service test-kafka 60
    
    print_status "Waiting for services to fully initialize..."
    sleep 10
    
    print_status "Setting up test environment..."
    export $(grep -v '^#' test.env | xargs)
    
    print_status "Running integration tests..."
    if [ "$SKIP_KAFKA" = true ]; then
        print_warning "Running tests without Kafka support"
        SKIP_KAFKA_TESTS=true go test -v -tags=integration ./integration_test.go
    else
        go test -v -tags=integration ./integration_test.go
    fi
    
    if [ $? -eq 0 ]; then
        print_status "Integration tests completed successfully!"
    else
        print_error "Integration tests failed!"
        exit 1
    fi
}

main "$@" 