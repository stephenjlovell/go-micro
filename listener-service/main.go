package main

import (
	"log"

	"github.com/stephenjlovell/go-micro/listener/event"
)

func main() {
	// connect to RabbitMQ
	conn, err := event.TryConnect()
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		log.Panic("unable to connect to RabbitMQ")
	}
	// create a consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		log.Panic(err)
	}
	// start listening for messages. watch the queue and consume events
	consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic(err)
	}
}
