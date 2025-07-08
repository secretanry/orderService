package models

import "wb-L0/modules/pg"

type Order struct {
	Id                int64          `gorm:"primary_key;auto_increment"`
	Uid               string         `gorm:"type:varchar(50);not null"`
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
	SmId              string         `gorm:"type:varchar(50)"`
	DateCreated       int64          `gorm:"type:int(11)"`
	OofShard          string         `gorm:"type:varchar(5)"`
}

func (Order) TableName() string {
	return "order"
}

func init() {
	pg.RegisterModel(new(Order))
}
