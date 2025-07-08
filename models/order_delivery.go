package models

import "wb-L0/modules/pg"

type OrderDelivery struct {
	Id      int64  `gorm:"primary_key;auto_increment"`
	OrderId int64  `gorm:"type:int;not null"`
	Name    string `gorm:"type:varchar(50);not null"`
	Phone   string `gorm:"type:varchar(12);not null"`
	Zip     string `gorm:"type:varchar(12);not null"`
	City    string `gorm:"type:varchar(50);not null"`
	Address string `gorm:"type:varchar(50);not null"`
	Region  string `gorm:"type:varchar(50);not null"`
	Email   string `gorm:"type:varchar(50);not null"`
}

func (OrderDelivery) TableName() string {
	return "order_delivery"
}

func init() {
	pg.RegisterModel(new(OrderDelivery))
}
