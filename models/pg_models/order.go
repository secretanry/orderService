package pg_models

import (
	"gorm.io/gorm"

	"wb-L0/modules/pg"
)

type Order struct {
	Id                int64          `gorm:"primaryKey;autoIncrement"`
	Uid               string         `gorm:"type:varchar(50);not null;unique"`
	TrackNumber       string         `gorm:"type:varchar(100)"`
	Entry             string         `gorm:"type:varchar(50)"`
	Delivery          *OrderDelivery `gorm:"-"`
	Payment           *OrderPayment  `gorm:"-"`
	Items             []*OrderItem   `gorm:"-"`
	Locale            string         `gorm:"type:varchar(5)"`
	InternalSignature string         `gorm:"type:text"`
	CustomerId        string         `gorm:"type:varchar(50)"`
	DeliveryService   string         `gorm:"type:varchar(50)"`
	Shardkey          string         `gorm:"type:varchar(50)"`
	SmId              int            `gorm:"type:integer"`
	DateCreated       int64          `gorm:"type:bigint"`
	OofShard          string         `gorm:"type:varchar(5)"`
}

func (*Order) TableName() string {
	return "order"
}

func init() {
	pg.RegisterModel(new(Order))
}

func InsertOrder(db *gorm.DB, order *Order) error {
	err := db.Create(order).Error
	return err
}

func GetOrderById(db *gorm.DB, uid string) (*Order, error) {
	order := new(Order)
	err := db.Where(&Order{Uid: uid}).First(order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (order *Order) loadPayment(db *gorm.DB) error {
	payment, err := GetOrderPaymentByOrderId(db, order.Id)
	if err != nil {
		return err
	}
	order.Payment = payment
	return nil
}

func (order *Order) loadDelivery(db *gorm.DB) error {
	delivery, err := GetDeliveryByOrderId(db, order.Id)
	if err != nil {
		return err
	}
	order.Delivery = delivery
	return nil
}

func (order *Order) loadItems(db *gorm.DB) error {
	items, err := GetOrderItemsByOrderId(db, order.Id)
	if err != nil {
		return err
	}
	order.Items = items
	return nil
}

func (order *Order) LoadAttributes(db *gorm.DB) error {
	var err error
	err = order.loadPayment(db)
	if err != nil {
		return err
	}
	err = order.loadDelivery(db)
	if err != nil {
		return err
	}
	err = order.loadItems(db)
	if err != nil {
		return err
	}
	return nil
}
