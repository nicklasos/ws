package main

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

var Connection *amqp.Connection
var Channel *amqp.Channel

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func QueueInit() {
	var err error

	Connection, err = amqp.Dial(os.Getenv("RABBITMQ_DSN"))
	failOnError(err, "Failed to connect to RabbitMQ")

	Channel, err = Connection.Channel()
	failOnError(err, "Failed to open a channel")
}

func QueueShutdown() {
	Channel.Close()
	Connection.Close()
}

func QueueRun(hub *Hub) {
	q, err := Channel.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE"), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	log.Println("Start queue")

	for {
		for d := range msgs {
			StackPush("ws.stack", 10, d.Body)
			hub.broadcast <- d.Body
		}
	}
}

func QueueSend() {
	q, err := Channel.QueueDeclare(
		os.Getenv("RABBITMQ_QUEUE"), // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "hello"
	err = Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	log.Printf("Queue Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
