package main

import (
	"flag"
	"log"

	"github.com/streadway/amqp"

	"gh-assan/rmsq/utils"
)

var url = "amqp://rabbitmq:rabbitmq@localhost:5672/"

const SensorDiscoveryActions = "SensorDiscoveryActions"
const SensorDiscoveryData = "SensorDiscoveryData"

var name = flag.String("name", "sensor", "name of the sensor")

// func start()

func main() {

	flag.Parse()

	connection, channel := utils.GetChannel(url)
	defer connection.Close()
	defer channel.Close()

	sendQueueName(channel)

	connectionExchange, channelExchange := utils.GetChannelForExchange(url, SensorDiscoveryActions)
	defer connectionExchange.Close()
	defer channelExchange.Close()

	go listenForDiscoverRequests(SensorDiscoveryActions, channelExchange, channel)

	select {}

}

func listenForDiscoverRequests(exchange string, exchangeChannel *amqp.Channel, channel *amqp.Channel) {

	queue := utils.GetQueueForExchange(exchange, exchangeChannel)

	msgs, _ := exchangeChannel.Consume(
		queue.Name, //queue string,
		"",         //consumer string,
		true,       //autoAck bool,
		false,      //exclusive bool,
		false,      //noLocal bool,
		false,      //noWait bool,
		nil)        //args amqp.Table)

	for range msgs {
		log.Println("received discovery request")
		sendQueueName(channel)
	}
}

func sendQueueName(channel *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(*name)}
	queue := utils.GetQueue(SensorDiscoveryData, channel)
	channel.Publish(
		"",         //exchange string,
		queue.Name, //key string,
		false,      //mandatory bool,
		false,      //immediate bool,
		msg)        //msg amqp.Publishing)
}
