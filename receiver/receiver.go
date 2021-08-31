package receiver

import (
	"log"

	"gh-assan/rmsq/utils"
)

const url = "amqp://rabbitmq:rabbitmq@localhost:5672/"
const queueName = "hello"

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Receive() {

	connection, channel := utils.GetChannel(url)

	defer connection.Close()
	defer channel.Close()

	queue := utils.GetQueue(queueName, channel)

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
