package sender

import (
	"log"

	"github.com/streadway/amqp"

	"gh-assan/rmsq/utils"
)

const url = "amqp://rabbitmq:rabbitmq@localhost:5672/"
const queueName = "hello"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Send() {

	connection, channel := utils.GetChannel(url)

	defer connection.Close()
	defer channel.Close()

	queue := utils.GetQueue(queueName, channel)

	for {
		body := "Hello World!"
		err := channel.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
	}
}
