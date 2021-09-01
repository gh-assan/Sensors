package utils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
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

func GetChannelForExchange(url string, exchange string) (*amqp.Connection, *amqp.Channel) {
	connection, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection to message broker")
	channel, err := connection.Channel()
	failOnError(err, "Failed to get channel for connection")

	err = channel.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)

	failOnError(err, "Failed to declare exchange")

	return connection, channel
}

func GetQueueForExchange(exchange string, channel *amqp.Channel) *amqp.Queue {

	queue, err := channel.QueueDeclare(
		"",    //name string,
		false, //durable bool,
		true,  //autoDelete bool,
		false, //exclusive bool,
		false, //noWait bool,
		nil)   //args amqp.Table)

	failOnError(err, "Failed to declare queue")

	err = channel.QueueBind(
		queue.Name, // queue name
		"",         // routing key
		exchange,   // exchange
		false,
		nil,
	)

	failOnError(err, "Failed to bind queue to exchange")

	return &queue
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
