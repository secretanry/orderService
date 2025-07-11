package convert

import (
	"time"

	"wb-L0/models/pg_models"
	"wb-L0/structs"
)

func PgToApiOrder(order *pg_models.Order) *structs.Order {
	parsedTime := time.Unix(order.DateCreated, 0).Format(time.RFC3339)
	return &structs.Order{
		OrderUid:          order.Uid,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          *PgToApiDelivery(order.Delivery),
		Payment:           *PgToApiPayment(order.Payment),
		Items:             PgToApiItems(order.Items),
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerId:        order.CustomerId,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmId:              order.SmId,
		DateCreated:       parsedTime,
		OofShard:          order.OofShard,
	}
}

func PgToApiDelivery(delivery *pg_models.OrderDelivery) *structs.Delivery {
	return &structs.Delivery{
		Name:    delivery.Name,
		Phone:   delivery.Phone,
		Zip:     delivery.Zip,
		City:    delivery.City,
		Address: delivery.Address,
		Region:  delivery.Region,
		Email:   delivery.Email,
	}
}

func PgToApiPayment(payment *pg_models.OrderPayment) *structs.Payment {
	return &structs.Payment{
		Transaction:  payment.Transaction,
		RequestId:    payment.RequestId,
		Currency:     payment.Currency,
		Provider:     payment.Provider,
		Amount:       payment.Amount,
		PaymentDt:    payment.PaymentDt,
		Bank:         payment.Bank,
		DeliveryCost: payment.DeliveryCost,
		GoodsTotal:   payment.GoodsTotal,
		CustomFee:    payment.CustomFee,
	}
}

func PgToApiItem(item *pg_models.OrderItem) *structs.Item {
	return &structs.Item{
		ChrtId:      item.ChrtId,
		TrackNumber: item.TrackNumber,
		Price:       item.Price,
		Rid:         item.Rid,
		Name:        item.Name,
		Sale:        item.Sale,
		Size:        item.Size,
		TotalPrice:  item.TotalPrice,
		NmId:        item.NmId,
		Brand:       item.Brand,
		Status:      item.Status,
	}
}

func PgToApiItems(items []*pg_models.OrderItem) []structs.Item {
	result := make([]structs.Item, len(items))
	for i, item := range items {
		result[i] = *PgToApiItem(item)
	}
	return result
}
