package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"wb-L0/models/pg_models"
	"wb-L0/modules/convert"
	"wb-L0/modules/pg"
	"wb-L0/structs"
)

type PostgresDatabase struct {
	db *pg.Postgres
}

func NewPostgres(postgres *pg.Postgres) *PostgresDatabase {
	return &PostgresDatabase{
		db: postgres,
	}
}

func (p *PostgresDatabase) InsertOrder(ctx context.Context, order *structs.Order) error {
	parsedTime, err := time.Parse(time.RFC3339, order.DateCreated)
	if err != nil {
		return ErrDataInvalid{err.Error()}
	}
	toInsertOrder := &pg_models.Order{
		Uid:               order.OrderUid,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerId:        order.CustomerId,
		DeliveryService:   order.DeliveryService,
		Shardkey:          order.Shardkey,
		SmId:              order.SmId,
		DateCreated:       parsedTime.Unix(),
		OofShard:          order.OofShard,
	}

	toInsertOrderDelivery := &pg_models.OrderDelivery{
		Name:    order.Delivery.Name,
		Phone:   order.Delivery.Phone,
		Zip:     order.Delivery.Zip,
		City:    order.Delivery.City,
		Address: order.Delivery.Address,
		Region:  order.Delivery.Region,
		Email:   order.Delivery.Email,
	}

	toInsertOrderPayment := &pg_models.OrderPayment{
		Transaction:  order.Payment.Transaction,
		RequestId:    order.Payment.RequestId,
		Currency:     order.Payment.Currency,
		Provider:     order.Payment.Provider,
		Amount:       order.Payment.Amount,
		PaymentDt:    order.Payment.PaymentDt,
		Bank:         order.Payment.Bank,
		DeliveryCost: order.Payment.DeliveryCost,
		GoodsTotal:   order.Payment.GoodsTotal,
		CustomFee:    order.Payment.CustomFee,
	}

	toInsertOrderItems := make([]*pg_models.OrderItem, len(order.Items))
	for i, item := range order.Items {
		toInsertOrderItems[i] = &pg_models.OrderItem{
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

	err = p.db.GetEngine(ctx).Transaction(func(tx *gorm.DB) error {
		err := pg_models.InsertOrder(tx, toInsertOrder)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return ErrDataInvalid{err.Error()}
			}
			return err
		}
		toInsertOrderDelivery.OrderId = toInsertOrder.Id
		err = pg_models.InsertOrderDelivery(tx, toInsertOrderDelivery)
		if err != nil {
			return err
		}
		toInsertOrderPayment.OrderId = toInsertOrder.Id
		err = pg_models.InsertOrderPayment(tx, toInsertOrderPayment)
		if err != nil {
			return err
		}
		for _, orderItem := range toInsertOrderItems {
			orderItem.OrderId = toInsertOrder.Id
			err = pg_models.InsertOrderItem(tx, orderItem)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return ErrDataInvalid{err.Error()}
	}
	return nil
}

func (p *PostgresDatabase) GetOrderById(ctx context.Context, oid string) (*structs.Order, error) {
	order, err := pg_models.GetOrderById(p.db.GetEngine(ctx), oid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound{Id: oid}
		}
		return nil, err
	}
	if order == nil {
		return nil, ErrOrderNotFound{Id: oid}
	}
	err = order.LoadAttributes(p.db.GetEngine(ctx))
	if err != nil {
		return nil, ErrInternal{Err: err.Error()}
	}
	return convert.PgToApiOrder(order), nil
}
