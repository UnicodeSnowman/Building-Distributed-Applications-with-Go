package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"pluralsight.com/datamanager"
	"pluralsight.com/dto"
	"pluralsight.com/qutils"
)

const url = "amqp://guest:guest@localhost:5672"

func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		qutils.PersistReadingsQueue,
		"",
		false,
		true,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to get access to messages")
	}

	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := &dto.SensorMessage{}
		dec.Decode(sd)

		err := datamanager.SaveReading(sd)

		if err != nil {
			log.Printf("Failed to save reading from sensor %v. Error: %s", sd.Name, err.Error())
		} else {
			msg.Ack(false)
		}
	}
}
