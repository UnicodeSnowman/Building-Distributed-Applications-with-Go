package qutils

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const SensorDiscoveryExchange = "SensorDiscovery"
const SensorListQueue = "SensorList"
const PersistReadingsQueue = "PersistReadings"

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection")
	ch, err := conn.Channel()
	failOnError(err, "Failed to get channel for connection")

	return conn, ch
}

func GetQueue(name string, ch *amqp.Channel, autoDelete bool) *amqp.Queue {
	declaration := QueueDeclaration{
		Name:       name,
		Durable:    false,
		AutoDelete: autoDelete,
		Exclusive:  false,
		NoWait:     false,
	}

	q, err := declareQueue(declaration, ch)

	failOnError(err, "Failed to get Queue from Channel")

	return q
}

type Consumable interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

func ConsumeQueue(ch Consumable, name string) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareQueue(declaration QueueDeclaration, ch *amqp.Channel) (*amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		declaration.Name,
		declaration.Durable,
		declaration.AutoDelete,
		declaration.Exclusive,
		declaration.NoWait,
		nil,
	)

	return &q, err
}

type QueueDeclaration struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
