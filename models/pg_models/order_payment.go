package pg_models

import (
	"gorm.io/gorm"

	"wb-L0/modules/pg"
)

type OrderPayment struct {
	Id           int64  `gorm:"primaryKey;autoIncrement"`
	OrderId      int64  `gorm:"type:int;not null"`
	Transaction  string `gorm:"type:varchar(50);not null"`
	RequestId    string `gorm:"type:varchar(50)"`
	Currency     string `gorm:"type:varchar(5);not null"`
	Provider     string `gorm:"type:varchar(50);not null"`
	Amount       int    `gorm:"type:int;not null"`
	PaymentDt    int64  `gorm:"type:int;not null"`
	Bank         string `gorm:"type:varchar(50);not null"`
	DeliveryCost int    `gorm:"type:int;not null"`
	GoodsTotal   int    `gorm:"type:int;not null"`
	CustomFee    int    `gorm:"type:int"`
}

func (OrderPayment) TableName() string {
	return "order_payment"
}

func init() {
	pg.RegisterModel(new(OrderPayment))
}

func InsertOrderPayment(db *gorm.DB, orderPayment *OrderPayment) error {
	err := db.Create(orderPayment).Error
	return err
}

func GetOrderPaymentByOrderId(db *gorm.DB, orderId int64) (*OrderPayment, error) {
	orderPayment := new(OrderPayment)
	err := db.Where(&OrderPayment{OrderId: orderId}).First(orderPayment).Error
	return orderPayment, err
}
