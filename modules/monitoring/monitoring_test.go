package monitoring_test

import (
	"context"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"wb-L0/modules/health"
	"wb-L0/modules/monitoring"
	"wb-L0/services/broker"
	"wb-L0/services/cache"
	"wb-L0/services/database"
	"wb-L0/structs"
)

// Mock implementations for health checks

type mockDB struct{}

func (m *mockDB) HealthCheck(ctx context.Context) error { return assert.AnError }
func (m *mockDB) GetOrderById(ctx context.Context, orderId string) (*structs.Order, error) {
	panic("not implemented")
}
func (m *mockDB) InsertOrder(ctx context.Context, order *structs.Order) error {
	panic("not implemented")
}

type mockCache struct{}

func (m *mockCache) HealthCheck(ctx context.Context) error { return assert.AnError }
func (m *mockCache) GetOrder(ctx context.Context, orderId string) (*structs.Order, error) {
	panic("not implemented")
}
func (m *mockCache) PutOrder(ctx context.Context, orderId string, order *structs.Order) error {
	panic("not implemented")
}

type mockBroker struct{}

func (m *mockBroker) HealthCheck(ctx context.Context) error                  { return assert.AnError }
func (m *mockBroker) StartConsuming(ctx context.Context) chan broker.Message { return nil }

func TestMonitoringInitialization(t *testing.T) {
	// Ensure clean state
	monitoring.ResetForTesting()

	// Set mocks for health checks
	database.SetDatabase(&mockDB{})
	cache.SetCache(&mockCache{})
	broker.SetBroker(&mockBroker{})

	// Use a new Prometheus registry for this test
	reg := prometheus.NewRegistry()

	// Test that monitoring can be initialized
	mon := &monitoring.Monitoring{}
	errChan := make(chan error, 1)

	err := mon.InitWithRegistry(errChan, reg)
	assert.NoError(t, err)

	// Test that second initialization doesn't fail (idempotent)
	err2 := mon.InitWithRegistry(errChan, reg)
	assert.NoError(t, err2)

	// Test health check functions
	healthChecker := mon.GetHealthChecker()
	assert.NotNil(t, healthChecker)

	err = healthChecker.CheckDatabaseHealth()
	// This should fail since database is not initialized in test
	assert.Error(t, err)

	err = healthChecker.CheckCacheHealth()
	// This should fail since cache is not initialized in test
	assert.Error(t, err)

	err = healthChecker.CheckBrokerHealth()
	// This should fail since kafka is not initialized in test
	assert.Error(t, err)

	// Test metrics helpers (should not panic)
	monitoring.IncrementHTTPRequests("GET", "/test", "200")
	monitoring.ObserveHTTPRequestDuration("GET", "/test", time.Millisecond*100)
	monitoring.IncrementCacheHits()
	monitoring.IncrementCacheMisses()
	monitoring.IncrementDatabaseQueries("select", "orders")
	monitoring.ObserveDatabaseQueryDuration("select", "orders", time.Millisecond*50)
	monitoring.IncrementKafkaMessagesProcessed("test-topic", "success")
	monitoring.IncrementOrderRetrieval()
	monitoring.ObserveOrderRetrievalDuration(time.Millisecond * 200)

	// Test logger
	logger := monitoring.GetLogger()
	assert.NotNil(t, logger)

	// Test tracer
	tracer := monitoring.GetTracer()
	assert.NotNil(t, tracer)

	// Test meter
	meter := monitoring.GetMeter()
	assert.NotNil(t, meter)

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = mon.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestHealthStatus(t *testing.T) {
	status := health.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Services: map[string]string{
			"database": "healthy",
			"cache":    "healthy",
			"kafka":    "healthy",
		},
	}

	assert.Equal(t, "healthy", status.Status)
	assert.Len(t, status.Services, 3)
	assert.Equal(t, "healthy", status.Services["database"])
}
