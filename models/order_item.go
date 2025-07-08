package models

import "wb-L0/modules/pg"

type OrderItem struct {
	Id          int64  `gorm:"primary_key;auto_increment"`
	OrderId     int64  `gorm:"type:int;not null"`
	ChrtId      int64  `gorm:"type:int;not null"`
	TrackNumber string `gorm:"type:varchar(50);not null"`
	Price       int    `gorm:"type:int;not null"`
	Rid         string `gorm:"type:varchar(50);not null"`
	Name        string `gorm:"type:varchar(50);not null"`
	Sale        int    `gorm:"type:int;not null"`
	Size        string `gorm:"type:varchar(5);not null"`
	TotalPrice  int    `gorm:"type:int;not null"`
	NmId        int64  `gorm:"type:int;not null"`
	Brand       string `gorm:"type:varchar(50);not null"`
	Status      int    `gorm:"type:int;not null"`
}

func (OrderItem) TableName() string {
	return "order_item"
}

func init() {
	pg.RegisterModel(new(OrderItem))
}
