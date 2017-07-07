package datamanager

import (
	"errors"
	"pluralsight.com/dto"
)

var sensors map[string]int

func SaveReading(reading *dto.SensorMessage) error {
	if sensors[reading.Name] == 0 {
		getSensors()
	}

	if sensors[reading.Name] == 0 {
		return errors.New("Unable to find sensor for name '" + reading.Name + "'")
	}

	q := `
		INSERT INTO sensor_readings
			(value, sensor_id, taken_on)
		VALUES
		($1, $2, $3)
	`

	_, err := db.Exec(q, reading.Value, sensors[reading.Name], reading.Timestamp)

	return err
}

func getSensors() {
	sensors = make(map[string]int)

	q := `
		SELECT id, name
		FROM sensors
	`

	rows, _ := db.Query(q)

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)

		sensors[name] = id
	}
}
