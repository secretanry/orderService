package models

import "wb-L0/modules/pg"

type OrderPayment struct {
	Id           int64  `gorm:"primary_key;auto_increment"`
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
