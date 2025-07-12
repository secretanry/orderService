package routing

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"

	"wb-L0/handlers"
	"wb-L0/modules/health"
	"wb-L0/modules/monitoring"
)

func MountSystemRoutes(r *gin.Engine) {
	r.Use(ApiContextMiddleware())
	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.json"),
		httpSwagger.DocExpansion("none"),
	)))
	r.GET("/docs/swagger.json", func(c *gin.Context) {
		ctx := handlers.GetApiContext(c)
		data, err := os.ReadFile("docs/swagger.json")
		if err != nil {
			ctx.ApiError(http.StatusInternalServerError, "Failed to read Swagger JSON")
			return
		}
		ctx.String(http.StatusOK, string(data))
	})

	// Health check endpoints
	r.GET("/health", healthCheck)
	r.GET("/health/ready", readinessCheck)
	r.GET("/health/live", livenessCheck)
}

func healthCheck(c *gin.Context) {
	status := health.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	// Get monitoring instance and health checker
	mon := monitoring.GetMonitoring()
	if mon == nil {
		status.Status = "unhealthy"
		status.Services["monitoring"] = "not initialized"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	healthChecker := mon.GetHealthChecker()

	// Check database health
	if err := healthChecker.CheckDatabaseHealth(); err != nil {
		status.Status = "unhealthy"
		status.Services["database"] = "error"
	} else {
		status.Services["database"] = "healthy"
	}

	// Check cache health
	if err := healthChecker.CheckCacheHealth(); err != nil {
		status.Status = "unhealthy"
		status.Services["cache"] = "error"
	} else {
		status.Services["cache"] = "healthy"
	}

	// Check kafka health
	if err := healthChecker.CheckBrokerHealth(); err != nil {
		status.Status = "unhealthy"
		status.Services["kafka"] = "error"
	} else {
		status.Services["kafka"] = "healthy"
	}

	if status.Status == "healthy" {
		c.JSON(http.StatusOK, status)
	} else {
		c.JSON(http.StatusServiceUnavailable, status)
	}
}

func readinessCheck(c *gin.Context) {
	// Readiness check - service is ready to receive traffic
	status := health.HealthStatus{
		Status:    "ready",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	// Get monitoring instance and health checker
	mon := monitoring.GetMonitoring()
	if mon == nil {
		status.Status = "not ready"
		status.Services["monitoring"] = "not initialized"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	healthChecker := mon.GetHealthChecker()

	// Check if all required services are available
	if err := healthChecker.CheckDatabaseHealth(); err != nil {
		status.Status = "not ready"
		status.Services["database"] = "not available"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	if err := healthChecker.CheckCacheHealth(); err != nil {
		status.Status = "not ready"
		status.Services["cache"] = "not available"
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}

	status.Services["database"] = "ready"
	status.Services["cache"] = "ready"
	c.JSON(http.StatusOK, status)
}

func livenessCheck(c *gin.Context) {
	// Liveness check - service is alive and running
	status := health.HealthStatus{
		Status:    "alive",
		Timestamp: time.Now().Format(time.RFC3339),
		Services:  make(map[string]string),
	}

	status.Services["application"] = "running"
	c.JSON(http.StatusOK, status)
}
