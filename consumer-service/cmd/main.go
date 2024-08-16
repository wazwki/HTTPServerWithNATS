package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("updates", func(msg *nats.Msg) {
		log.Printf("Получено сообщение: %s", string(msg.Data))
	})
	if err != nil {
		log.Fatal(err)
	}

	select {}
}
