package event

import (
	"encoding/json"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	helpers "github.com/stephenjlovell/json-helpers"
)

const (
	loggerServiceUrl = "http://logger-service/log"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

const (
	retryCount          = 10
	maxBackOffInSeconds = 60
)

func TryConnect() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	backOffInSeconds := 1
	for i := 0; i < retryCount; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			log.Printf("unable to connect: %s\n", err.Error())
			sleepFor := math.Min(float64(backOffInSeconds), float64(maxBackOffInSeconds))
			time.Sleep(time.Second * time.Duration(sleepFor))
			log.Printf("%fs - retrying connection...\n", sleepFor)
			backOffInSeconds *= 2
		} else {
			log.Println("connected to RabbitMQ...")
			return conn, nil
		}
	}
	return nil, err
}

func NewConsumer(conn *amqp.Connection) (*Consumer, error) {
	consumer := &Consumer{
		conn: conn,
	}
	if err := consumer.setup(); err != nil {
		return nil, err
	}
	return consumer, nil
}

func (c *Consumer) setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(ch)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}
	// bind to the given topics
	for _, topic := range topics {
		ch.QueueBind(
			queue.Name,
			topic,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	// FIXME this smells a little funky
	forever := make(chan bool)

	go func() {
		for d := range messages {
			payload := Payload{}
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)
		}
	}()

	log.Println("Waiting for messages...")
	<-forever
	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log":
		if err := doLogging(payload); err != nil {
			log.Println(err)
		}

		// add other cases as needed

	default:
		if err := doLogging(payload); err != nil {
			log.Println(err)
		}
	}
}

// TODO: duplicated in broker-service
func doLogging(payload Payload) (err error) {
	response, err := helpers.DoRequest("POST", loggerServiceUrl, payload)

	defer func() {
		closeErr := response.Body.Close()
		if err == nil {
			err = closeErr
		}
	}()
	return
}
