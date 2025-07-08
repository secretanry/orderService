package pg

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wb-L0/modules/config"
)

var (
	models   []interface{}
	globalDb *gorm.DB
)

type Postgres struct {
	Db *gorm.DB
}

func (p *Postgres) Init(_ chan error) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		config.GetConfig().DbHost, config.GetConfig().DbUser, config.GetConfig().DbPass,
		config.GetConfig().DbName, config.GetConfig().DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %v", err)
	}
	p.Db = db
	globalDb = db
	err = db.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("failed to migrate models: %v", err)
	}
	return nil
}

func (p *Postgres) SuccessfulMessage() string {
	return "Postgres successfully initialized"
}

func (p *Postgres) Shutdown(_ context.Context) error {
	db, err := globalDb.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB connection: %v", err)
	}
	return db.Close()
}

func GetEngine(ctx context.Context) *gorm.DB {
	return globalDb.WithContext(ctx)
}

func RegisterModel(model interface{}) {
	models = append(models, model)
}
