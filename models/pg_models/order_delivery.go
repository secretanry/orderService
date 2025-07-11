package pg_models

import (
	"gorm.io/gorm"

	"wb-L0/modules/pg"
)

type OrderDelivery struct {
	Id      int64  `gorm:"primaryKey;autoIncrement"`
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

func InsertOrderDelivery(db *gorm.DB, delivery *OrderDelivery) error {
	err := db.Create(delivery).Error
	return err
}

func GetDeliveryByOrderId(db *gorm.DB, oid int64) (*OrderDelivery, error) {
	delivery := new(OrderDelivery)
	err := db.Where(&OrderDelivery{OrderId: oid}).First(delivery).Error
	return delivery, err
}
