package initializer

import (
	"context"
	"log"

	"wb-L0/modules/config"
	"wb-L0/modules/graceful"
	"wb-L0/modules/pg"
	"wb-L0/modules/server"
)

var units []Initializable

func init() {
	units = []Initializable{
		new(config.Config),
		new(server.Server),
		new(pg.Postgres),
	}
}

func Init() {
	var err error
	err = graceful.Init()
	if err != nil {
		log.Fatal(err)
	}
	errChan := make(chan error)
	go handleErrors(errChan)
	for _, u := range units {
		select {
		case <-graceful.GetContext().Done():
			return
		default:
			err = u.Init(errChan)
			if err != nil {
				log.Println(err)
				graceful.DoShutdown()
			}
			log.Println(u.SuccessfulMessage())
		}
	}
}

func Shutdown(ctx context.Context) {
	for _, u := range units {
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

type Initializable interface {
	Init(chan error) error
	SuccessfulMessage() string
	Shutdown(ctx context.Context) error
}
