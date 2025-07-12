# Automatic Grafana Setup

This document describes the automatic Grafana setup that creates dashboards and visualizations for your WB-L0 service monitoring.

## Overview

The automatic setup includes:
- **Automatic Datasource Configuration**: Prometheus datasource is automatically configured
- **Pre-built Dashboard**: Comprehensive WB-L0 Service Dashboard with 11 panels
- **Provisioning**: All configurations are applied automatically on startup
- **Health Checks**: Scripts verify that everything is working correctly

## What Gets Created Automatically

### 1. Prometheus Datasource
- **Name**: Prometheus
- **URL**: http://prometheus:9090
- **Access**: Proxy
- **Default**: Yes

### 2. WB-L0 Service Dashboard
The dashboard includes the following panels:

#### **Performance Metrics**
- **HTTP Request Rate**: Requests per second by method and endpoint
- **HTTP Request Duration**: 95th and 50th percentile response times
- **Error Rate**: 4xx and 5xx error rates over time

#### **Cache Metrics**
- **Cache Hit Rate**: Percentage of cache hits with color-coded thresholds
- **Cache Operations**: Hits vs misses over time

#### **Database Metrics**
- **Database Query Rate**: Queries per second by operation and table
- **Database Query Duration**: 95th percentile query times

#### **Business Metrics**
- **Order Retrieval Rate**: Orders retrieved per second
- **Kafka Messages Processed**: Message processing rate by topic and status

#### **System Health**
- **Service Health Status**: Overall service availability
- **HTTP Status Codes**: Distribution of response codes (pie chart)

#### **Advanced Features**
- **Templating**: Filter by endpoint and HTTP method
- **Annotations**: Automatic deployment markers
- **Auto-refresh**: Dashboard updates every 5 seconds

## File Structure

```
grafana/
├── provisioning/
│   ├── datasources/
│   │   └── prometheus.yml          # Prometheus datasource config
│   └── dashboards/
│       ├── dashboard.yml           # Dashboard provider config
│       └── wb-l0-dashboard.json    # Main dashboard definition
scripts/
└── setup-grafana.sh               # Automatic setup script
```

## Usage

### Start Monitoring Stack with Auto-Setup
```bash
make monitoring-up
```

This command will:
1. Create `monitoring.env` from your existing configuration
2. Start Prometheus and Grafana containers
3. Wait for services to be ready
4. Automatically configure Grafana
5. Verify that everything is working

### Manual Setup (if needed)
```bash
make setup-grafana
```

### Check Status
```bash
make monitoring-status
```

### View Logs
```bash
make monitoring-logs
```

## Dashboard Features

### **Real-time Monitoring**
- All metrics update every 5 seconds
- Historical data available for trend analysis
- Automatic time range selection

### **Interactive Filtering**
- Filter by HTTP method (GET, POST, etc.)
- Filter by endpoint (/api/orders, /health, etc.)
- Multi-select capabilities

### **Color-coded Thresholds**
- **Cache Hit Rate**: Red (<70%), Yellow (70-90%), Green (>90%)
- **Service Health**: Red (down), Green (up)
- **Error Rates**: Visual indicators for different error types

### **Responsive Layout**
- 12-column grid system
- Optimized for different screen sizes
- Dark theme for better visibility

## Customization

### Adding New Panels
1. Edit `grafana/provisioning/dashboards/wb-l0-dashboard.json`
2. Add new panel definitions
3. Restart the monitoring stack: `make monitoring-restart`

### Modifying Queries
All panels use Prometheus queries. Common patterns:
```promql
# Rate of events
rate(metric_name[5m])

# Percentile calculations
histogram_quantile(0.95, rate(metric_bucket[5m]))

# Error filtering
rate(metric{status=~"5.."}[5m])
```

### Adding New Metrics
1. Add metrics to your application code
2. Update the dashboard JSON to include new panels
3. Restart the monitoring stack

## Troubleshooting

### Dashboard Not Appearing
```bash
# Check if provisioning worked
make monitoring-status

# Check Grafana logs
make monitoring-logs

# Manual setup
make setup-grafana
```

### No Data in Panels
1. Verify Prometheus is scraping metrics: http://localhost:9090/targets
2. Check if your application is running and generating metrics
3. Verify the Prometheus datasource is configured correctly

### Access Issues
- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Health Check**: http://localhost:8080/health

## Advanced Configuration

### Custom Dashboard Variables
The dashboard includes templating variables:
- `$endpoint`: Filter by API endpoint
- `$method`: Filter by HTTP method

### Annotations
Automatic annotations mark:
- Service deployments
- Configuration changes
- Error events

### Alerting (Future Enhancement)
The dashboard is prepared for alerting rules:
- High error rates
- Low cache hit rates
- Slow response times
- Service downtime

## Performance Considerations

### Dashboard Optimization
- Queries use 5-minute rate windows for efficiency
- Panels refresh every 5 seconds
- Historical data retention: 200 hours (configurable)

### Resource Usage
- Grafana: ~100MB RAM
- Prometheus: ~200MB RAM (varies with data volume)
- Storage: ~1GB per day (varies with metrics volume)

## Next Steps

1. **Start using the dashboard**: `make monitoring-up`
2. **Customize for your needs**: Modify the dashboard JSON
3. **Add alerting**: Configure Prometheus alerting rules
4. **Scale monitoring**: Add more services to the monitoring stack

The automatic setup provides a production-ready monitoring solution that you can start using immediately! 