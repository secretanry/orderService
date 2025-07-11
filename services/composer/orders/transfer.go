package orders

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"wb-L0/modules/graceful"
	"wb-L0/services/broker"
	"wb-L0/services/database"
	"wb-L0/structs"
)

func StartDataTransfer() {
	messageChan := broker.GetBroker().StartConsuming(graceful.GetContext())
	go func() {
		for {
			select {
			case <-graceful.GetContext().Done():
				return
			case message := <-messageChan:
				if message.Err != nil {
					log.Println(message.Err)
					continue
				}
				order, err := unmarshall(message.Value)
				if err != nil {
					log.Printf("Data marshalling failed: %s. Message skipped!", err.Error())
					err = message.Ack()
					if err != nil {
						log.Printf("Ack failed: %s. Message skipped!", err.Error())
					}
					continue
				}
				ctx, cancel := context.WithTimeout(graceful.GetContext(), 5*time.Second)
				err = database.GetDatabase().InsertOrder(ctx, order)
				cancel()
				if err != nil {
					if database.IsErrDataInvalid(err) {
						log.Printf("Insert order failed: %s. Message skipped!", err.Error())
						err = message.Ack()
						if err != nil {
							log.Printf("Ack failed: %s. Message skipped!", err.Error())
						}
					} else {
						log.Printf("Insert order failed: %s. Retrying!", err.Error())
						err = message.Nack()
						if err != nil {
							log.Printf("Nack failed: %s. Message skipped", err.Error())
						}
					}
				}
			}
		}
	}()
}

func unmarshall(data []byte) (*structs.Order, error) {
	var order structs.Order
	err := json.Unmarshal(data, &order)
	return &order, err
}
