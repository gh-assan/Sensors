package main

import (
	"log"
	"time"

	"github.com/streadway/amqp"

	"gh-assan/rmsq/utils"
)

var url = "amqp://rabbitmq:rabbitmq@localhost:5672/"

const SensorDiscoveryActions = "SensorDiscoveryActions"
const SensorDiscoveryData = "SensorDiscoveryData"

func main() {

	go func() {
		connection, channel := utils.GetChannelForExchange(url, SensorDiscoveryActions)
		defer connection.Close()
		defer channel.Close()

		signal := time.Tick(5 * time.Second)

		for range signal {
			sendDiscoveryRequest(channel)
		}
	}()

	go func() {

		connectionDiscovery, channelDiscovery := utils.GetChannel(url)
		defer connectionDiscovery.Close()
		defer channelDiscovery.Close()

		discoveryQueue := utils.GetQueue(SensorDiscoveryData, channelDiscovery)
		// channel.QueueBind(
		// 	discoveryQueue.Name, //name string,
		// 	"",                  //key string,
		// 	SensorDiscoveryQueueName, //exchange string,
		// 	false, //noWait bool,
		// 	nil)   //args amqp.Table)

		listenForDiscoverRespone(channelDiscovery, discoveryQueue.Name)
		select {}
	}()
	select {}

}

func listenForDiscoverRespone(channel *amqp.Channel, queueName string) {
	msgs, _ := channel.Consume(
		queueName, //queue string,
		"",        //consumer string,
		true,      //autoAck bool,
		false,     //exclusive bool,
		false,     //noLocal bool,
		false,     //noWait bool,
		nil)       //args amqp.Table)

	for msg := range msgs {
		log.Println("received discovery request", string(msg.Body))
		// sendQueueName(channel)
	}
}

func sendDiscoveryRequest(channel *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(`{action:"status"}`)}
	log.Println("sending discovery request")
	channel.Publish(
		SensorDiscoveryActions, //exchange string,
		"",                     //key string,
		false,                  //mandatory bool,
		false,                  //immediate bool,
		msg)                    //msg amqp.Publishing)
	log.Println("sent discovery request")
}
