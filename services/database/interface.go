package database

import (
	"context"

	"wb-L0/structs"
)

var dbInstance Database

type Database interface {
	InsertOrder(context.Context, *structs.Order) error
	GetOrderById(ctx context.Context, oid string) (*structs.Order, error)
	HealthCheck(ctx context.Context) error
}

func SetDatabase(database Database) {
	dbInstance = database
}

func GetDatabase() Database {
	return dbInstance
}
