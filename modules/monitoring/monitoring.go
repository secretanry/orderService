package monitoring

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"

	"wb-L0/modules/config"
	"wb-L0/modules/health"
)

var (
	// Global instances
	logger             *zap.Logger
	tracer             trace.Tracer
	meter              metric.Meter
	provider           *sdkmetric.MeterProvider
	monitoringInstance *Monitoring
	initialized        bool // Flag to prevent duplicate initialization
	// Prometheus metrics
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	cacheHits = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
	)
	cacheMisses = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
	)
	databaseQueries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "database_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "table"},
	)
	databaseQueryDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)
	kafkaMessagesProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "kafka_messages_processed_total",
			Help: "Total number of Kafka messages processed",
		},
		[]string{"topic", "status"},
	)
	// OpenTelemetry metrics
	orderRetrievalCounter  metric.Int64Counter
	orderRetrievalDuration metric.Float64Histogram
)

type Monitoring struct {
	server   *http.Server
	health   *health.Checker
	Registry *prometheus.Registry // Optional custom registry
}

// InitWithRegistry allows injecting a custom Prometheus registry (for tests)
func (m *Monitoring) InitWithRegistry(errChan chan error, registry *prometheus.Registry) error {
	if initialized {
		monitoringInstance = m
		return nil
	}
	monitoringInstance = m
	m.health = health.NewChecker()
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.TimeKey = "timestamp"
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logConfig.EncoderConfig.StacktraceKey = "stacktrace"
	var err error
	logger, err = logConfig.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	if registry == nil {
		registry = prometheus.DefaultRegisterer.(*prometheus.Registry)
	}
	m.Registry = registry
	// Register Prometheus metrics
	m.Registry.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		cacheHits,
		cacheMisses,
		databaseQueries,
		databaseQueryDuration,
		kafkaMessagesProcessed,
	)
	// Initialize OpenTelemetry MeterProvider with Prometheus exporter
	exporter, err := otelprom.New()
	if err != nil {
		return fmt.Errorf("failed to create prometheus exporter: %w", err)
	}
	provider = sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter = provider.Meter("wb-L0")
	// Create OpenTelemetry metrics
	orderRetrievalCounter, err = meter.Int64Counter(
		"order_retrieval_total",
		metric.WithDescription("Total number of order retrievals"),
	)
	if err != nil {
		return fmt.Errorf("failed to create order retrieval counter: %w", err)
	}
	orderRetrievalDuration, err = meter.Float64Histogram(
		"order_retrieval_duration_seconds",
		metric.WithDescription("Order retrieval duration in seconds"),
	)
	if err != nil {
		return fmt.Errorf("failed to create order retrieval duration histogram: %w", err)
	}
	// Initialize tracer
	tracer = otel.Tracer("wb-L0")
	// Start metrics server
	config := config.GetConfig()
	metricsPort := 8081 // Default metrics port if not configured
	if config != nil && config.MetricsPort > 0 {
		metricsPort = config.MetricsPort
	}
	metricsAddr := fmt.Sprintf("0.0.0.0:%d", metricsPort)
	m.server = &http.Server{
		Addr:    metricsAddr,
		Handler: promhttp.HandlerFor(m.Registry, promhttp.HandlerOpts{}),
	}
	go func() {
		logger.Info("Starting metrics server", zap.String("address", metricsAddr))
		if err := m.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("metrics server error: %w", err)
		}
	}()

	// Mark as initialized
	initialized = true
	return nil
}

// Init uses the default registry
func (m *Monitoring) Init(errChan chan error) error {
	return m.InitWithRegistry(errChan, nil)
}

func (m *Monitoring) SuccessfulMessage() string {
	return "Monitoring successfully initialized"
}

func (m *Monitoring) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down monitoring")
	if m.server != nil {
		return m.server.Shutdown(ctx)
	}
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return logger
}

// GetTracer returns the global tracer instance
func GetTracer() trace.Tracer {
	return tracer
}

// GetMeter returns the global meter instance
func GetMeter() metric.Meter {
	return meter
}

// Metrics helpers
func IncrementHTTPRequests(method, endpoint, status string) {
	httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
}

func ObserveHTTPRequestDuration(method, endpoint string, duration time.Duration) {
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

func IncrementCacheHits() {
	cacheHits.Inc()
}

func IncrementCacheMisses() {
	cacheMisses.Inc()
}

func IncrementDatabaseQueries(operation, table string) {
	databaseQueries.WithLabelValues(operation, table).Inc()
}

func ObserveDatabaseQueryDuration(operation, table string, duration time.Duration) {
	databaseQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

func IncrementKafkaMessagesProcessed(topic, status string) {
	kafkaMessagesProcessed.WithLabelValues(topic, status).Inc()
}

// OpenTelemetry metrics helpers
func IncrementOrderRetrieval() {
	if orderRetrievalCounter != nil {
		orderRetrievalCounter.Add(context.Background(), 1)
	}
}

func ObserveOrderRetrievalDuration(duration time.Duration) {
	if orderRetrievalDuration != nil {
		orderRetrievalDuration.Record(context.Background(), duration.Seconds())
	}
}

// GetHealthChecker returns the health checker instance
func (m *Monitoring) GetHealthChecker() *health.Checker {
	return m.health
}

// GetMonitoring returns the global monitoring instance
func GetMonitoring() *Monitoring {
	return monitoringInstance
}

// ResetForTesting resets the initialization state for testing
func ResetForTesting() {
	initialized = false
	monitoringInstance = nil
	logger = nil
	tracer = nil
	meter = nil
	provider = nil
	orderRetrievalCounter = nil
	orderRetrievalDuration = nil
}
