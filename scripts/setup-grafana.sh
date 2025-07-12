#!/bin/bash

# Setup Grafana with automatic configuration
echo "Setting up Grafana with automatic configuration..."

# Wait for Grafana to be ready
echo "Waiting for Grafana to be ready..."
until curl -s http://localhost:3000/api/health; do
    echo "Grafana not ready yet, waiting..."
    sleep 5
done

echo "Grafana is ready!"

# Create organization and user (optional - using default admin/admin)
echo "Grafana is configured with:"
echo "  - Username: admin"
echo "  - Password: admin"
echo "  - URL: http://localhost:3000"

# Check if datasource is configured
echo "Checking Prometheus datasource..."
DS_STATUS=$(curl -s -u admin:admin http://localhost:3000/api/datasources/name/Prometheus | jq -r '.name // "not_found"')

if [ "$DS_STATUS" = "Prometheus" ]; then
    echo "✅ Prometheus datasource is configured"
else
    echo "❌ Prometheus datasource not found"
fi

# Check if dashboard is imported
echo "Checking dashboard..."
DASHBOARD_STATUS=$(curl -s -u admin:admin http://localhost:3000/api/search?query=WB-L0 | jq -r '.[0].title // "not_found"')

if [ "$DASHBOARD_STATUS" = "WB-L0 Service Dashboard" ]; then
    echo "✅ WB-L0 Service Dashboard is imported"
else
    echo "❌ Dashboard not found"
fi

echo ""
echo "🎉 Grafana setup complete!"
echo ""
echo "Access your monitoring stack:"
echo "  📊 Grafana: http://localhost:3000 (admin/admin)"
echo "  📈 Prometheus: http://localhost:9090"
echo "  🔍 Your App: http://localhost:8080"
echo "  ❤️  Health Check: http://localhost:8080/health"
echo "  📊 Metrics: http://localhost:8081/metrics"
echo ""
echo "Dashboard includes:"
echo "  - HTTP Request Rate & Duration"
echo "  - Cache Hit Rate"
echo "  - Database Query Performance"
echo "  - Order Retrieval Metrics"
echo "  - Service Health Status"
echo "  - Error Rates"
echo "  - Kafka Message Processing" 