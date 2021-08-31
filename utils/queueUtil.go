package utils

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

//const SensorDiscoveryExchange = "SensorDiscovery"

func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	connection, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection to message broker")
	channel, err := connection.Channel()
	failOnError(err, "Failed to get channel for connection")

	return connection, channel
}

func GetQueue(name string, channel *amqp.Channel) *amqp.Queue {
	queue, err := channel.QueueDeclare(
		name,  //name string,
		false, //durable bool,
		false, //autoDelete bool,
		false, //exclusive bool,
		false, //noWait bool,
		nil)   //args amqp.Table)

	failOnError(err, "Failed to declare queue")

	return &queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
