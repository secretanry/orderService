package monitoring

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	CorrelationIDKey = "correlation_id"
	TraceIDKey       = "trace_id"
	SpanIDKey        = "span_id"
)

// MonitoringMiddleware adds monitoring, tracing, and logging to HTTP requests
func MonitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Generate correlation ID
		correlationID := uuid.New().String()
		c.Set(CorrelationIDKey, correlationID)

		// Create span for tracing
		ctx, span := GetTracer().Start(
			c.Request.Context(),
			"http.request",
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.Path),
				attribute.String("http.user_agent", c.Request.UserAgent()),
				attribute.String("correlation_id", correlationID),
			),
		)
		defer span.End()

		// Set trace context
		c.Request = c.Request.WithContext(ctx)

		// Add trace IDs to response headers
		c.Header("X-Correlation-ID", correlationID)
		c.Header("X-Trace-ID", span.SpanContext().TraceID().String())
		c.Header("X-Span-ID", span.SpanContext().SpanID().String())

		// Log request start
		logger.Info("HTTP request started",
			zap.String("correlation_id", correlationID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("remote_addr", c.ClientIP()),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Record metrics
		status := c.Writer.Status()
		IncrementHTTPRequests(c.Request.Method, c.Request.URL.Path, string(rune(status)))
		ObserveHTTPRequestDuration(c.Request.Method, c.Request.URL.Path, duration)

		// Add response attributes to span
		span.SetAttributes(
			attribute.Int("http.status_code", status),
			attribute.Int64("http.response_size", int64(c.Writer.Size())),
			attribute.Int64("http.request_duration_ms", duration.Milliseconds()),
		)

		// Log request completion
		logLevel := zap.InfoLevel
		if status >= 400 {
			logLevel = zap.WarnLevel
		}
		if status >= 500 {
			logLevel = zap.ErrorLevel
		}

		logger.Check(logLevel, "HTTP request completed").Write(
			zap.String("correlation_id", correlationID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", status),
			zap.Int("response_size", c.Writer.Size()),
			zap.Duration("duration", duration),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)
	}
}

// GetCorrelationID extracts correlation ID from gin context
func GetCorrelationID(c *gin.Context) string {
	if id, exists := c.Get(CorrelationIDKey); exists {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

// GetTraceContext extracts trace context from gin context
func GetTraceContext(c *gin.Context) context.Context {
	return c.Request.Context()
}

// LogWithContext creates a logger with correlation ID and trace context
func LogWithContext(c *gin.Context) *zap.Logger {
	correlationID := GetCorrelationID(c)
	if correlationID != "" {
		return logger.With(zap.String("correlation_id", correlationID))
	}
	return logger
}

// SpanFromContext creates a span from gin context
func SpanFromContext(c *gin.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(GetTraceContext(c), name, opts...)
}
