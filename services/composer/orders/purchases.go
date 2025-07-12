package orders

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"wb-L0/modules/monitoring"
	"wb-L0/services/cache"
	"wb-L0/services/database"
	"wb-L0/structs"

	"go.uber.org/zap"
)

func GetOrderById(ctx context.Context, orderId string) (*structs.Order, error) {
	start := time.Now()
	defer func() {
		monitoring.ObserveOrderRetrievalDuration(time.Since(start))
		monitoring.IncrementOrderRetrieval()
	}()

	// Create span for tracing
	ctx, _ = monitoring.GetTracer().Start(ctx, "order.retrieval",
		trace.WithAttributes(
			attribute.String("order_id", orderId),
		),
	)

	order, err := cache.GetCache().GetOrder(ctx, orderId)
	if order != nil {
		monitoring.IncrementCacheHits()
		monitoring.GetLogger().Info("Cache hit",
			zap.String("order_id", orderId))
		return order, nil
	}
	if !cache.IsErrCacheMiss(err) {
		return nil, err
	}

	monitoring.IncrementCacheMisses()
	monitoring.GetLogger().Info("Cache miss, looking in database",
		zap.String("order_id", orderId))

	order, err = database.GetDatabase().GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}

	err = cache.GetCache().PutOrder(ctx, orderId, order)
	if err != nil {
		monitoring.GetLogger().Warn("Unable to cache order",
			zap.String("order_id", orderId), zap.Error(err))
	} else {
		monitoring.GetLogger().Info("Cache updated",
			zap.String("order_id", orderId))
	}

	return order, nil
}
