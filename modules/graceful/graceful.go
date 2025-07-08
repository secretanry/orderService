package graceful

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	gracefulCtx    context.Context
	gracefulCancel context.CancelFunc
)

type Graceful struct{}

func Init() error {
	ctx, cancel := context.WithCancel(context.Background())
	gracefulCtx = ctx
	gracefulCancel = cancel
	quitQueue := make(chan os.Signal, 1)
	signal.Notify(quitQueue, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-quitQueue:
			DoShutdown()
		}
	}()
	return nil
}

func GetContext() context.Context {
	return gracefulCtx
}

func DoShutdown() {
	log.Println("Shutdown signal received. Turning off.")
	gracefulCancel()
}
