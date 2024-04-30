package facades

import (
	"log"

	"goravel/packages/queue"
	"goravel/packages/queue/contracts"
)

func Queue() contracts.Queue {
	instance, err := queue.App.Make(queue.Binding + ".client")
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.Queue)
}
