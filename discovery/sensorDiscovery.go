package discovery

import (
	"log"
	"time"

	"github.com/streadway/amqp"

	"gh-assan/rmsq/utils"
)

const SensorDiscoveryActions = "SensorDiscoveryActions"
const SensorDiscoveryData = "SensorDiscoveryData"

type SensorDiscovery struct {
}

func NewSensorDiscovery() *SensorDiscovery {
	return &SensorDiscovery{}
}

func (sensorDiscovery *SensorDiscovery) Start() {

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

		listenForDiscoverRespone(channelDiscovery, discoveryQueue.Name)
		select {}
	}()
	select {}

}

func listenForSensorData(sensorName string) {
	sensorListener := NewSensorListener(sensorName)

	sensorListener.Start()
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
		go listenForSensorData(string(msg.Body))
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
