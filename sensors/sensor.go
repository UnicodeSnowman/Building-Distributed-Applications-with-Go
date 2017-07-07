package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"pluralsight.com/dto"
	"pluralsight.com/qutils"
	"strconv"
	"time"
)

var url = "amqp://guest:guest@localhost:5672"
var name = flag.String("name", "sensor", "name of the sensor")
var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")
var max = flag.Float64("max", 5., "maximum value for generated readings")
var min = flag.Float64("min", 1., "minimum value for generated readings")
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var value = r.Float64()*(*max-*min) + *min
var nom = (*max-*min)/2 + *min

func publishMsg(ch *amqp.Channel, exchange string, queueName string, data amqp.Publishing) {
	ch.Publish(exchange, queueName, false, false, data)
}

func main() {
	flag.Parse()
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	dataQueue := qutils.GetQueue(*name, ch, false)

	// Publish the name of this sensor's queue to the "registry" queue
	//sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)
	//publishMsg(ch, sensorQueue.Name, amqp.Publishing{Body: []byte(*name)})
	publishMsg(ch, "amq.fanout", "", amqp.Publishing{Body: []byte(*name)})

	discoveryQueue := qutils.GetQueue("", ch, true)
	ch.QueueBind(
		discoveryQueue.Name,
		"",
		qutils.SensorDiscoveryExchange,
		false,
		nil,
	)

	go listenForDiscoverRequests(discoveryQueue.Name, ch)

	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")
	signal := time.Tick(dur)
	buf := new(bytes.Buffer)

	for range signal {
		log.Printf("Read sent. Value: %v\n", value)
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}

		buf.Reset()
		gob.NewEncoder(buf).Encode(reading)
		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		publishMsg(ch, "", dataQueue.Name, msg)
		calcValue()
	}
}

func listenForDiscoverRequests(name string, ch *amqp.Channel) {
	msgs, _ := qutils.ConsumeQueue(ch, name)

	for range msgs {
		publishMsg(ch, "amq.fanout", "", amqp.Publishing{Body: []byte(name)})
	}
}

func calcValue() {
	var maxStep, minStep float64
	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	value += r.Float64()*(maxStep-minStep) + minStep
}
