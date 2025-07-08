package main

import (
	"context"
	"time"

	"wb-L0/modules/graceful"
	"wb-L0/modules/initializer"
)

func main() {
	initializer.Init()
	<-graceful.GetContext().Done()
	ctx, cancel := context.WithTimeout(graceful.GetContext(), time.Second*5)
	defer cancel()
	initializer.Shutdown(ctx)
}
