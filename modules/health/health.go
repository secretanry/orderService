package health

import (
	"context"
	"fmt"
)

// HealthChecker interface for service health checks
type HealthChecker interface {
	HealthCheck(context.Context) error
}

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// Checker holds references to service health checkers
type Checker struct {
	database HealthChecker
	cache    HealthChecker
	broker   HealthChecker
}

// NewChecker creates a new health checker
func NewChecker() *Checker {
	return &Checker{}
}

// SetDatabase sets the database health checker
func (c *Checker) SetDatabase(checker HealthChecker) {
	c.database = checker
}

// SetCache sets the cache health checker
func (c *Checker) SetCache(checker HealthChecker) {
	c.cache = checker
}

// SetBroker sets the broker health checker
func (c *Checker) SetBroker(checker HealthChecker) {
	c.broker = checker
}

// CheckDatabaseHealth checks database health
func (c *Checker) CheckDatabaseHealth() error {
	if c.database == nil {
		return fmt.Errorf("database not initialized")
	}
	return c.database.HealthCheck(context.Background())
}

// CheckCacheHealth checks cache health
func (c *Checker) CheckCacheHealth() error {
	if c.cache == nil {
		return fmt.Errorf("cache not initialized")
	}
	return c.cache.HealthCheck(context.Background())
}

// CheckBrokerHealth checks broker health
func (c *Checker) CheckBrokerHealth() error {
	if c.broker == nil {
		return fmt.Errorf("broker not initialized")
	}
	return c.broker.HealthCheck(context.Background())
}
