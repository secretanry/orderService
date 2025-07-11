package pg_models

import (
	"gorm.io/gorm"

	"wb-L0/modules/pg"
)

type OrderItem struct {
	Id          int64  `gorm:"primaryKey;autoIncrement"`
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

func InsertOrderItem(db *gorm.DB, orderItem *OrderItem) error {
	err := db.Create(orderItem).Error
	return err
}

func GetOrderItemsByOrderId(db *gorm.DB, oid int64) ([]*OrderItem, error) {
	items := make([]*OrderItem, 0)
	err := db.Where(&OrderItem{OrderId: oid}).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
