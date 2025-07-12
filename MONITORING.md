# WB-L0 Service Monitoring

This document describes the comprehensive monitoring solution implemented for the WB-L0 service.

## Overview

The monitoring solution includes:

- **Prometheus Metrics**: Application and business metrics
- **OpenTelemetry Tracing**: Distributed tracing for request flows
- **Structured Logging**: JSON-formatted logs with correlation IDs
- **Health Checks**: Service health monitoring
- **Grafana Dashboards**: Visualization of metrics

## Architecture

### Monitoring Components

1. **Metrics Server** (Port 8081): Prometheus metrics endpoint
2. **Health Endpoints**:
   - `/health`: Overall service health
   - `/health/ready`: Readiness probe
   - `/health/live`: Liveness probe
3. **Application Server** (Port 8080): Main application with monitoring middleware

### Metrics Collected

#### HTTP Metrics
- `http_requests_total`: Total HTTP requests by method, endpoint, and status
- `http_request_duration_seconds`: Request duration histograms

#### Cache Metrics
- `cache_hits_total`: Total cache hits
- `cache_misses_total`: Total cache misses

#### Database Metrics
- `database_queries_total`: Database queries by operation and table
- `database_query_duration_seconds`: Query duration histograms

#### Business Metrics
- `order_retrieval_total`: Total order retrievals
- `order_retrieval_duration_seconds`: Order retrieval duration

#### Kafka Metrics
- `kafka_messages_processed_total`: Kafka messages by topic and status

## Setup Instructions

### 1. Environment Configuration

Add monitoring configuration to your `.env` file:

```bash
# Monitoring Configuration
METRICS_PORT=8081
LOG_LEVEL=info
```

### 2. Start Monitoring Stack

```bash
# Start Prometheus and Grafana
docker-compose -f docker-compose.monitoring.yml up -d
```

### 3. Start Application

```bash
# Start your application
go run main.go
```

## Monitoring Endpoints

### Metrics Endpoint
- **URL**: `http://localhost:8081/metrics`
- **Description**: Prometheus metrics in text format

### Health Endpoints
- **Health Check**: `http://localhost:8080/health`
- **Readiness Probe**: `http://localhost:8080/health/ready`
- **Liveness Probe**: `http://localhost:8080/health/live`

### Grafana Dashboard
- **URL**: `http://localhost:3000`
- **Credentials**: admin/admin

## Key Features

### 1. Request Tracing

Every HTTP request gets:
- Unique correlation ID
- Trace ID and Span ID
- Request/response logging
- Performance metrics

### 2. Structured Logging

Logs include:
- Correlation ID for request tracking
- Structured JSON format
- Log levels (INFO, WARN, ERROR)
- Contextual information

### 3. Health Monitoring

Health checks verify:
- Database connectivity
- Cache availability
- Kafka connectivity
- Application status

### 4. Performance Monitoring

Track:
- Request rates and latencies
- Cache hit rates
- Database query performance
- Business metrics

## Usage Examples

### Viewing Metrics

```bash
# View all metrics
curl http://localhost:8081/metrics

# View health status
curl http://localhost:8080/health

# View specific metrics
curl http://localhost:8081/metrics | grep http_requests_total
```

### Prometheus Queries

```promql
# Request rate by endpoint
rate(http_requests_total[5m])

# 95th percentile response time
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Cache hit rate
rate(cache_hits_total[5m]) / (rate(cache_hits_total[5m]) + rate(cache_misses_total[5m])) * 100
```

### Log Analysis

```bash
# View logs with correlation ID
grep "correlation_id" logs/app.log

# View error logs
grep '"level":"error"' logs/app.log
```

## Alerting

### Recommended Alerts

1. **High Error Rate**
   ```promql
   rate(http_requests_total{status=~"5.."}[5m]) > 0.1
   ```

2. **High Response Time**
   ```promql
   histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
   ```

3. **Low Cache Hit Rate**
   ```promql
   rate(cache_hits_total[5m]) / (rate(cache_hits_total[5m]) + rate(cache_misses_total[5m])) < 0.8
   ```

4. **Service Down**
   ```promql
   up{job="wb-l0"} == 0
   ```

## Troubleshooting

### Common Issues

1. **Metrics not appearing**
   - Check if metrics server is running on port 8081
   - Verify Prometheus configuration
   - Check application logs for errors

2. **Health checks failing**
   - Verify database connectivity
   - Check cache service status
   - Ensure Kafka is accessible

3. **High memory usage**
   - Monitor application metrics
   - Check for memory leaks
   - Review cache configuration

### Debug Commands

```bash
# Check service status
curl -f http://localhost:8080/health

# View application logs
tail -f logs/app.log

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Test metrics endpoint
curl http://localhost:8081/metrics | head -20
```

## Performance Considerations

1. **Metrics Cardinality**: Avoid high cardinality labels
2. **Log Volume**: Configure log rotation
3. **Storage**: Monitor Prometheus storage usage
4. **Network**: Use appropriate scrape intervals

## Security

1. **Metrics Access**: Restrict access to metrics endpoint
2. **Health Checks**: Don't expose sensitive information
3. **Logs**: Avoid logging sensitive data
4. **Network**: Use internal networks for monitoring

## Future Enhancements

1. **Custom Metrics**: Add business-specific metrics
2. **Alerting**: Configure alerting rules
3. **Log Aggregation**: Centralized log management
4. **APM Integration**: Application Performance Monitoring
5. **Distributed Tracing**: Jaeger integration 