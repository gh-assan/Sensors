package sensor

import (
	"log"

	"github.com/streadway/amqp"

	"gh-assan/rmsq/utils"
)

const SensorDiscoveryActions = "SensorDiscoveryActions"
const SensorDiscoveryData = "SensorDiscoveryData"

func (sensor *Sensor) HandleDiscovery() {
	connection, channel := utils.GetChannel(url)
	defer connection.Close()
	defer channel.Close()

	sensor.sendQueueName(channel)

	connectionExchange, channelExchange := utils.GetChannelForExchange(url, SensorDiscoveryActions)
	defer connectionExchange.Close()
	defer channelExchange.Close()

	go sensor.listenForDiscoverRequests(SensorDiscoveryActions, channelExchange, channel)

	select {}
}

func (sensor *Sensor) listenForDiscoverRequests(exchange string, exchangeChannel *amqp.Channel, channel *amqp.Channel) {

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
		sensor.sendQueueName(channel)
	}
}

func (sensor *Sensor) sendQueueName(channel *amqp.Channel) {
	msg := amqp.Publishing{Body: []byte(sensor.Name)}
	queue := utils.GetQueue(SensorDiscoveryData, channel)
	channel.Publish(
		"",         //exchange string,
		queue.Name, //key string,
		false,      //mandatory bool,
		false,      //immediate bool,
		msg)        //msg amqp.Publishing)
}
