package main

import (
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stephenjlovell/go-micro/broker/event"
)

const webPort = 80

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// connect to RabbitMQ
	conn, err := event.TryConnect()
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		log.Panic("unable to connect to RabbitMQ")
	}

	log.Printf("starting broker service on port %d\n", webPort)
	app := Config{
		Rabbit: conn,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
