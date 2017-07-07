package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	go client()
	go server()

	fmt.Println("Running...")

	// options for blocking forever

	// with `select`...
	// select {}

	// with a channel...
	<-make(chan struct{})

	// with `Scanln`...
	//var a string
	//fmt.Scanln(&a)
}

// Basic server to generate "hello world" messages every second
// Pretty much just for testing purposes
func server() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello Rabbit"),
	}

	for {
		ch.Publish("", q.Name, false, false, msg)
		time.Sleep(1000 * time.Millisecond)
	}
}

// Basic client to consume messages from our default queue. Again, mostly
// for intro/testing purposes
func client() {
	conn, ch, q := getQueue()
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name,
		"",
		false, // normally set to true, or need to manually ack, if message is used to trigger something else
		false,
		false,
		false,
		nil,
	)

	failOnError(err, "Failed to register a consumer")

	for msg := range msgs {
		log.Printf("Received message with message: %s", msg.Body)
	}
}

func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a Channel")
	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	return conn, ch, &q
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
