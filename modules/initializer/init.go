package initializer

import (
	"context"
	"fmt"
	"log"

	"wb-L0/modules/config"
	"wb-L0/modules/graceful"
	"wb-L0/modules/kafka"
	"wb-L0/modules/monitoring"
	"wb-L0/modules/pg"
	"wb-L0/modules/redis"
	"wb-L0/modules/server"
	"wb-L0/services/broker"
	"wb-L0/services/cache"
	"wb-L0/services/composer/orders"
	"wb-L0/services/database"
)

var (
	initializedUnitsList []Initializable
)

func Init() {
	var err error
	err = graceful.Init()
	if err != nil {
		log.Fatal(err)
	}
	unitsList := []Initializable{
		new(config.Config),
		new(monitoring.Monitoring),
		new(server.Server),
	}
	initUnits(unitsList)
	optionalUnits, err := identifyOptionalUnits()
	if err != nil {
		log.Println(err)
		graceful.DoShutdown()
	}
	initUnits(optionalUnits)
	select {
	case <-graceful.GetContext().Done():
	default:
		orders.StartDataTransfer()
	}
}

func Shutdown(ctx context.Context) {
	for _, u := range initializedUnitsList {
		err := u.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func handleErrors(errChan chan error) {
	for {
		select {
		case <-graceful.GetContext().Done():
			return
		case err := <-errChan:
			log.Printf("asyncronous task error: %v", err)
			graceful.DoShutdown()
			return
		}
	}
}

func identifyOptionalUnits() ([]Initializable, error) {
	units := make([]Initializable, 0)
	switch config.GetConfig().BrokerType {
	case "kafka":
		instance := new(kafka.Kafka)
		units = append(units, instance)
		broker.SetBroker(broker.NewKafkaBroker(instance))
	default:
		return nil, fmt.Errorf("unknown broker type: %s", config.GetConfig().BrokerType)
	}
	switch config.GetConfig().DbType {
	case "postgres":
		instance := new(pg.Postgres)
		units = append(units, instance)
		database.SetDatabase(database.NewPostgres(instance))
	default:
		return nil, fmt.Errorf("unknown db type: %s", config.GetConfig().DbType)
	}
	switch config.GetConfig().CacheType {
	case "memory":
		cache.SetCache(cache.NewMemoryCache())
	case "redis":
		instance := new(redis.Redis)
		units = append(units, instance)
		cache.SetCache(cache.NewRedisCache(instance))
	default:
		return nil, fmt.Errorf("unknown cache type: %s", config.GetConfig().CacheType)
	}
	return units, nil
}

func initUnits(units []Initializable) {
	errChan := make(chan error)
	go handleErrors(errChan)
	for i := range units {
		u := units[i]
		select {
		case <-graceful.GetContext().Done():
			return
		default:
			err := u.Init(errChan)
			if err != nil {
				log.Println(err)
				graceful.DoShutdown()
			} else {
				initializedUnitsList = append(initializedUnitsList, u)
				log.Println(u.SuccessfulMessage())
			}
		}
	}
	return
}

type Initializable interface {
	Init(chan error) error
	SuccessfulMessage() string
	Shutdown(ctx context.Context) error
}
