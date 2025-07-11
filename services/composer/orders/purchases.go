package orders

import (
	"context"
	"log"

	"wb-L0/services/cache"
	"wb-L0/services/database"
	"wb-L0/structs"
)

func GetOrderById(ctx context.Context, orderId string) (*structs.Order, error) {
	order, err := cache.GetCache().GetOrder(ctx, orderId)
	if order != nil {
		log.Printf("Cache hit for: %s!", orderId)
		return order, nil
	}
	if !cache.IsErrCacheMiss(err) {
		return nil, err
	}
	log.Println("Cache miss, looking in database")
	order, err = database.GetDatabase().GetOrderById(ctx, orderId)
	if err != nil {
		return nil, err
	}
	err = cache.GetCache().PutOrder(ctx, orderId, order)
	if err != nil {
		log.Printf("Unable to cache order: %s!", orderId)
	} else {
		log.Printf("Cache updates for: %s", orderId)
	}
	return order, nil
}
