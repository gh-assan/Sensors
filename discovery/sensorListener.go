package discovery

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/streadway/amqp"

	"gh-assan/rmsq/utils"
)

var url = "amqp://rabbitmq:rabbitmq@localhost:5672/"

type SensorListener struct {
	Name string
}

type SensorMessage struct {
	Name      string
	Value     float64
	Timestamp time.Time
}

func NewSensorListener(sensorName string) *SensorListener {
	return &SensorListener{Name: sensorName}
}

func (sensorListener *SensorListener) Start() {

	connection, channel := utils.GetChannel(url)
	defer connection.Close()
	defer channel.Close()

	queue := utils.GetQueue(sensorListener.Name, channel)

	listenForData(channel, queue.Name)
}

func listenForData(channel *amqp.Channel, queueName string) {
	msgs, _ := channel.Consume(
		queueName, //queue string,
		"",        //consumer string,
		true,      //autoAck bool,
		false,     //exclusive bool,
		false,     //noLocal bool,
		false,     //noWait bool,
		nil)       //args amqp.Table)

	for msg := range msgs {

		reader := bytes.NewReader(msg.Body)
		decoder := gob.NewDecoder(reader)
		sensorReading := new(SensorMessage)
		decoder.Decode(sensorReading)

		log.Println("received sensor data ", sensorReading)
	}
}
