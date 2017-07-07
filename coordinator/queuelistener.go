package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/streadway/amqp"
	"pluralsight.com/dto"
	"pluralsight.com/qutils"
)

var url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener(ea *EventAggregator) *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
		ea:      ea,
	}

	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) DiscoverSensors() {
	ql.ch.ExchangeDeclare(
		qutils.SensorDiscoveryExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)

	ql.ch.Publish(
		qutils.SensorDiscoveryExchange,
		"",
		false,
		false,
		amqp.Publishing{},
	)
}

func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch, true)
	ql.ch.QueueBind(
		q.Name,
		"",
		"amq.fanout",
		false,
		nil,
	)

	msgs, _ := qutils.ConsumeQueue(ql.ch, q.Name)

	ql.DiscoverSensors()

	for msg := range msgs {
		ql.ea.PublishEvent("DataSourceDiscovered", string(msg.Body))
		name := string(msg.Body)
		sourceChan, _ := qutils.ConsumeQueue(ql.ch, name)

		if ql.sources[name] == nil {
			ql.sources[name] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	for msg := range msgs {
		sd := new(dto.SensorMessage)
		gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(sd)

		// using `%v` for "default" format, works well for structs
		fmt.Printf("Received message: %v\n", sd)

		ed := EventData{
			Name:      sd.Name,
			Value:     sd.Value,
			Timestamp: sd.Timestamp,
		}

		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}
