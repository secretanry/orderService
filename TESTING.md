# Testing Guide

This document provides comprehensive information about testing the wb-L0 project, including unit tests, integration tests, and how to set up the testing environment.

## Overview

The project includes several types of tests:

1. **Unit Tests** - Test individual components in isolation using mocks
2. **Integration Tests** - Test the complete application flow including HTTP endpoints, database, cache, and Kafka messaging
3. **End-to-End Tests** - Test the full system with production-like external dependencies

## Prerequisites

Before running tests, ensure you have the following installed:

- Go 1.24 or later
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

## Test Structure

```
wb-L0/
├── integration_test.go          # Main integration test suite with Kafka
├── handlers/
│   └── purchases_test.go        # Unit tests for HTTP handlers
├── services/
│   ├── composer/orders/
│   │   └── purchases_test.go    # Unit tests for order services
│   └── cache/
│       └── memory_test.go       # Unit tests for cache services
├── docker-compose.test.yml      # Test services configuration (KRaft Kafka)
├── test.env                     # Test environment variables
├── scripts/
│   ├── kafka-test-init.sh
│   └── run_integration_tests.sh # Test runner script with cleanup
├── Makefile                     # Test automation commands
└── TESTING.md                   # This file
```

## Running Tests

### Unit Tests

Unit tests test individual components in isolation using mocks and stubs.

```bash
# Run all unit tests (excluding integration tests)
make test-unit

# Run unit tests for specific package
go test -v ./handlers
go test -v ./services/composer/orders
go test -v ./services/cache

# Run unit tests with coverage
go test -v -cover ./...

# Run unit tests and generate coverage report
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Integration Tests

Integration tests test the complete application flow including HTTP endpoints, database operations, cache functionality, and Kafka messaging.

#### Using Makefile (Recommended)

```bash
# Run integration tests with full setup and teardown
make test-integration


# Run all tests (unit + integration)
make test
```

#### Using Test Script

```bash
# Make script executable
chmod +x scripts/run_integration_tests.sh

# Run integration tests
./scripts/run_integration_tests.sh
```

#### Manual Setup

1. Start test services:
```bash
docker-compose -f docker-compose.test.yml up -d --build
```

2. Wait for services to be ready:
```bash
# Check service health
docker-compose -f docker-compose.test.yml ps
```

3. Set environment variables:
```bash
export $(cat test.env | xargs)
```

4. Run integration tests:
```bash
go test -v ./integration_test.go
```

5. Clean up:
```bash
docker-compose -f docker-compose.test.yml down -v
```

## Test Services

The integration tests use the following services defined in `docker-compose.test.yml`:

- **PostgreSQL** (port 5433) - Test database with full schema
- **Redis** (port 6380) - Test cache with authentication
- **Kafka** (port 9093) - Test message broker in KRaft mode (same image as production)

These services run on different ports to avoid conflicts with your development environment and use the same Docker images as production.

## Test Configuration

The test environment is configured via `test.env`:

```env
# Test Environment Configuration
APP_PORT=8080
RUN_MODE=test

# Database Configuration
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5433
DB_USER=test_user
DB_PASS=test_pass
DB_NAME=test_db

# Kafka Configuration
BROKER_TYPE=kafka
KAFKA_URL=localhost:9093
KAFKA_CONSUMER_GROUP=test-consumer-group
KAFKA_TOPIC=orders

# Cache Configuration
CACHE_TYPE=redis
REDIS_HOST=localhost
REDIS_PORT=6380
REDIS_PASS=test_pass
REDIS_DATABASE=1
```

## Test Coverage

The integration tests cover comprehensive scenarios:

### HTTP Layer
- ✅ Successful order retrieval
- ✅ Order not found scenarios
- ✅ Invalid request handling
- ✅ Error response formats
- ✅ Parameter validation
- ✅ Concurrent request handling

### Database Layer
- ✅ Order creation and retrieval
- ✅ Complex nested data handling (delivery, payment, items)
- ✅ Data integrity verification
- ✅ Connection failure handling
- ✅ Multiple orders handling

### Cache Layer
- ✅ Cache hit/miss scenarios
- ✅ Cache update verification
- ✅ Cache failure handling
- ✅ Performance comparison (cache vs database)
- ✅ Cache data integrity

### Kafka Messaging Layer
- ✅ Message production to Kafka topics
- ✅ Message consumption and validation
- ✅ Message format verification (JSON serialization)
- ✅ Error handling scenarios
- ✅ Performance testing (348+ messages/sec)
- ✅ Topic management and metadata

### Integration Scenarios
- ✅ Complete request flow (HTTP → Cache → Database)
- ✅ Kafka producer/consumer integration
- ✅ Concurrent request handling
- ✅ Service failure scenarios
- ✅ Performance characteristics
- ✅ End-to-end message flow

## Test Data

The integration tests create realistic test data including:

- Orders with complete delivery information
- Payment details with all required fields
- Multiple items per order with various data types
- Complex nested structures
- Edge cases and error scenarios

## Kafka Integration Testing

The integration tests now include comprehensive Kafka functionality testing:

### Kafka Tests Included

1. **`TestKafkaIntegration`**
   - Produces messages to Kafka topics
   - Verifies topic metadata and configuration
   - Tests batch message production
   - Validates partition and offset information

2. **`TestKafkaConsumerIntegration`**
   - Creates Kafka consumer instances
   - Sends test messages and consumes them
   - Validates message content and format
   - Tests end-to-end message flow

3. **`TestKafkaMessageFormat`**
   - Tests complex order JSON serialization
   - Validates message format integrity
   - Ensures data consistency across serialization/deserialization

4. **`TestKafkaErrorHandling`**
   - Tests error scenarios (non-existent topics)
   - Validates graceful error handling
   - Tests invalid message formats

5. **`TestKafkaPerformance`**
   - Measures Kafka throughput and latency
   - Tests performance with multiple messages
   - Validates performance requirements

### Kafka Configuration

- **Mode**: KRaft (no Zookeeper dependency)
- **Image**: Same custom build as production
- **Topics**: `orders` (auto-created if not exists)
- **Port**: 9093 (external), 9092 (internal)
- **Authentication**: None (test environment)

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure test services are not running on the same ports as your development environment
2. **Service startup**: Kafka in KRaft mode may take 30-60 seconds to start. The test script includes health checks
3. **Database connection**: Verify PostgreSQL is accessible on port 5433
4. **Redis connection**: Check Redis authentication and port 6380 configuration
5. **Kafka KRaft issues**: If you encounter cluster ID inconsistencies, the test script automatically cleans up volumes

### Debug Mode

To run tests with more verbose output:

```bash
# Run with debug logging
go test -v -run TestIntegrationSuite ./integration_test.go

# Run specific test
go test -v -run TestKafkaIntegration ./integration_test.go
go test -v -run TestGetOrderSuccess ./integration_test.go

# Run with timeout
go test -v -timeout 5m ./integration_test.go
```

### Manual Service Verification

```bash
# Check PostgreSQL
psql -h localhost -p 5433 -U test_user -d test_db

# Check Redis
redis-cli -h localhost -p 6380 -a test_pass ping

# Check Kafka (KRaft mode)
kafka-console-consumer.sh --bootstrap-server localhost:9093 --topic orders --from-beginning

# Check Kafka topics
kafka-topics.sh --bootstrap-server localhost:9093 --list
```

### Service Logs

```bash
# View service logs
docker-compose -f docker-compose.test.yml logs test-kafka
docker-compose -f docker-compose.test.yml logs test-postgres
docker-compose -f docker-compose.test.yml logs test-redis
```

## Continuous Integration

For CI/CD pipelines, you can use the following commands:

```yaml
# Example GitHub Actions step
- name: Run Tests
  run: |
    make test-integration

# Example with specific test types
- name: Run Unit Tests
  run: make test-unit

- name: Run Integration Tests
  run: make test-integration
```

## Performance Testing

The integration tests include comprehensive performance checks:

- **Response time measurements** for HTTP endpoints
- **Cache vs database performance** comparison
- **Concurrent request handling** (10+ concurrent requests)
- **Kafka throughput testing** (348+ messages/sec)
- **Memory usage monitoring**
- **Connection pool efficiency**

### Performance Benchmarks

- **Cache hit**: < 50ms response time
- **Database query**: < 200ms response time
- **Kafka throughput**: > 300 messages/sec
- **Concurrent requests**: 10+ simultaneous requests handled correctly

## Test Results Summary

The current test suite provides:

- **16 comprehensive integration tests** covering all system components
- **Full Kafka integration** with producer/consumer testing
- **Production-like environment** using same Docker images
- **Automatic cleanup** and health checks
- **Performance validation** for all critical paths
- **Error scenario coverage** including service failures

All tests run successfully in a production-like environment with KRaft Kafka, PostgreSQL, and Redis, ensuring your application works correctly in real-world conditions. 