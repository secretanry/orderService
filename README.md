# WB-L0 - Go Microservice with Comprehensive Monitoring

A production-ready Go microservice built with Gin framework, featuring PostgreSQL, Redis caching, Kafka message processing, and comprehensive monitoring with Prometheus, Grafana, and OpenTelemetry.

## 🚀 Features

- **High Performance**: Built with Gin framework for fast HTTP handling
- **Data Persistence**: PostgreSQL for reliable data storage
- **Caching**: Redis and in-memory cache support
- **Message Processing**: Kafka integration for event-driven architecture
- **Comprehensive Monitoring**: 
  - Prometheus metrics collection
  - Grafana dashboards with automatic provisioning
  - OpenTelemetry distributed tracing
  - Structured logging with correlation IDs
  - Health checks for all components
- **Production Ready**: Graceful shutdown, config management, and error handling
- **Testing**: Integration tests and monitoring tests included

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│   Gin Server    │───▶│   PostgreSQL    │
│                 │    │   (Port 8080)   │    │   Database      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Redis Cache   │
                       │   (Optional)    │
                       └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Kafka Broker  │
                       │   (Event Bus)   │
                       └─────────────────┘
                              │
                              ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Grafana       │◀───│   Prometheus    │◀───│   Metrics       │
│   (Port 3000)   │    │   (Port 9090)   │    │   (Port 8081)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 📋 Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose**
- **PostgreSQL** (or use Docker)
- **Redis** (optional, can use in-memory cache)
- **Kafka** (or use Docker)

## 🛠️ Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd wb-L0
```

### 2. Environment Configuration

Create a `.env` file based on the example:

```bash
# Application Configuration
APP_PORT=8080
RUN_MODE=debug
LOG_LEVEL=info

# Database Configuration
DB_TYPE=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=password
DB_NAME=wb_l0

# Kafka Configuration
BROKER_TYPE=kafka
KAFKA_URL=localhost:9092
KAFKA_CONSUMER_GROUP=wb-l0-group
KAFKA_TOPIC=orders

# Cache Configuration
CACHE_TYPE=redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASS=
REDIS_DATABASE=0

# Monitoring Configuration
METRICS_PORT=8081
```

### 3. Start the Application

```bash
# Run directly
go run main.go

# Or build and run
go build -o wb-l0
./wb-l0
```

### 4. Start Monitoring Stack

```bash
# Start Prometheus and Grafana
make monitoring-up

# Or manually
docker-compose -f docker-compose.monitoring.yml up -d
```

## 📊 Monitoring Setup

The project includes a comprehensive monitoring solution with automatic setup:

### Automatic Setup

```bash
# Start monitoring stack with automatic Grafana configuration
make monitoring-up
```

This command:
- Creates `monitoring.env` from your existing environment
- Starts Prometheus and Grafana containers
- Automatically configures Grafana with:
  - Prometheus datasource
  - Pre-built dashboard with all metrics
  - Proper authentication (admin/admin)

### Manual Setup

If you prefer manual setup:

1. **Start Prometheus & Grafana:**
   ```bash
   docker-compose -f docker-compose.monitoring.yml up -d
   ```

2. **Configure Grafana:**
   - Access Grafana at http://localhost:3000 (admin/admin)
   - Add Prometheus datasource: `http://prometheus:9090`
   - Import dashboard from `grafana-dashboard.json`

### Monitoring URLs

- **Grafana Dashboard**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Application Metrics**: http://localhost:8081/metrics
- **Health Check**: http://localhost:8080/health

## 🔧 API Endpoints

### Health Checks

```bash
# Overall health
GET /health

# Readiness probe
GET /health/ready

# Liveness probe
GET /health/live
```

### Order Management

```bash
# Get order by ID
GET /api/order/{order_id}

# Example
curl http://localhost:8080/api/order/b563feb7b2b84b6test
```

### System Endpoints

```bash
# Root endpoint
GET /

# Metrics endpoint (for Prometheus)
GET /metrics
```

## 📈 Monitoring Features

### Metrics Collected

- **HTTP Metrics**: Request count, duration, status codes
- **Cache Metrics**: Hit/miss rates, operation counts
- **Database Metrics**: Query count, duration by operation/table
- **Kafka Metrics**: Messages processed by topic/status
- **Business Metrics**: Order retrieval count and duration
- **System Metrics**: Go runtime metrics, memory usage

### Dashboard Panels

The Grafana dashboard includes:

1. **HTTP Request Rate & Duration**
2. **Cache Hit Rate & Operations**
3. **Database Query Performance**
4. **Order Retrieval Metrics**
5. **Service Health Status**
6. **Error Rates (4xx/5xx)**
7. **Kafka Message Processing**

### Logging

- **Structured JSON logs** with correlation IDs
- **Request tracing** with trace/span IDs
- **Performance metrics** for each request
- **Error tracking** with stack traces

## 🧪 Testing

### Run Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# Monitoring tests
go test ./monitoring_test.go
```

### Test with Monitoring

```bash
# Start monitoring stack
make monitoring-up

# Run tests
go test -tags=integration ./...

# Check metrics
curl http://localhost:8081/metrics
```

## 🐳 Docker Support

### Development

```bash
# Start all services
docker-compose -f docker-compose.dev.yml up -d

# Start monitoring only
docker-compose -f docker-compose.monitoring.yml up -d
```

### Production

```bash
# Build image
docker build -t wb-l0 .

# Run container
docker run -p 8080:8080 -p 8081:8081 wb-l0
```

## 📁 Project Structure

```
wb-L0/
├── main.go                 # Application entry point
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies checksum
├── .env                    # Environment configuration
├── docker-compose.yml      # Main Docker Compose
├── docker-compose.monitoring.yml  # Monitoring stack
├── Dockerfile              # Application container
├── Makefile                # Build and deployment scripts
├── README.md               # This file
├── MONITORING.md           # Detailed monitoring guide
├── TESTING.md              # Testing documentation
├── GRAFANA-SETUP.md        # Grafana setup guide
├── handlers/               # HTTP request handlers
│   ├── purchases.go        # Order management handlers
│   └── helpers.go          # Handler utilities
├── modules/                # Core application modules
│   ├── config/             # Configuration management
│   ├── monitoring/         # Monitoring and observability
│   ├── health/             # Health check system
│   ├── graceful/           # Graceful shutdown
│   ├── initializer/        # Application initialization
│   ├── kafka/              # Kafka integration
│   ├── pg/                 # PostgreSQL integration
│   ├── redis/              # Redis integration
│   └── server/             # HTTP server setup
├── services/               # Business logic services
│   ├── broker/             # Message broker interface
│   ├── cache/              # Cache interface and implementations
│   ├── composer/           # Service orchestration
│   └── database/           # Database interface and implementation
├── models/                 # Data models
│   └── pg_models/          # PostgreSQL-specific models
├── structs/                # Shared data structures
├── routing/                # HTTP routing configuration
├── grafana/                # Grafana provisioning
│   └── provisioning/
│       ├── dashboards/     # Dashboard configurations
│       └── datasources/    # Data source configurations
├── scripts/                # Utility scripts
└── docs/                   # API documentation
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_PORT` | HTTP server port | 8080 |
| `RUN_MODE` | Application mode (debug/prod) | debug |
| `DB_TYPE` | Database type | postgres |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_USER` | Database username | postgres |
| `DB_PASS` | Database password | - |
| `DB_NAME` | Database name | wb_l0 |
| `BROKER_TYPE` | Message broker type | kafka |
| `KAFKA_URL` | Kafka broker URL | localhost:9092 |
| `KAFKA_CONSUMER_GROUP` | Kafka consumer group | wb-l0-group |
| `KAFKA_TOPIC` | Kafka topic name | orders |
| `CACHE_TYPE` | Cache type (redis/memory) | redis |
| `REDIS_HOST` | Redis host | localhost |
| `REDIS_PORT` | Redis port | 6379 |
| `REDIS_PASS` | Redis password | - |
| `REDIS_DATABASE` | Redis database number | 0 |
| `METRICS_PORT` | Metrics server port | 8081 |
| `LOG_LEVEL` | Logging level | info |

## 🚀 Deployment

### Local Development

```bash
# Start all services
make dev-up

# Run application
go run main.go

# Start monitoring
make monitoring-up
```

### Production Deployment

1. **Build the application:**
   ```bash
   make build
   ```

2. **Set production environment:**
   ```bash
   export RUN_MODE=production
   export LOG_LEVEL=warn
   ```

3. **Deploy with Docker:**
   ```bash
   docker-compose -f docker-compose.yml up -d
   ```

4. **Start monitoring:**
   ```bash
   docker-compose -f docker-compose.monitoring.yml up -d
   ```

## 📊 Monitoring Commands

```bash
# Start monitoring stack
make monitoring-up

# Stop monitoring stack
make monitoring-down

# View monitoring status
make monitoring-status

# Setup Grafana automatically
make setup-grafana
```

## 🔍 Troubleshooting

### Common Issues

1. **Port conflicts:**
   - Ensure ports 8080, 8081, 9090, 3000 are available
   - Check with `lsof -i :<port>`

2. **Database connection:**
   - Verify PostgreSQL is running
   - Check connection credentials in `.env`

3. **Kafka connection:**
   - Ensure Kafka broker is accessible
   - Check topic exists and is properly configured

4. **Monitoring issues:**
   - Verify Prometheus can reach `host.docker.internal:8081`
   - Check Grafana datasource configuration

### Debug Commands

```bash
# Check application health
curl http://localhost:8080/health

# View metrics
curl http://localhost:8081/metrics

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# View application logs
tail -f logs/app.log
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [Prometheus](https://prometheus.io/) - Monitoring system
- [Grafana](https://grafana.com/) - Visualization platform
- [OpenTelemetry](https://opentelemetry.io/) - Observability framework

---

**Need help?** Check the [MONITORING.md](MONITORING.md) for detailed monitoring setup or [TESTING.md](TESTING.md) for testing instructions. 