package queue

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

func Connect() *amqp.Channel {
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")
	conn, err := amqp.Dial(dsn)
	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func StartConsuming(in chan []byte, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_CONSUMER_QUEUE"), // name
		true,                                 // durable
		false,                                // delete when unused
		false,                                // exclusive
		false,                                // no-wait
		nil,                                  // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,      // queue
		"go-worker", // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			in <- []byte(d.Body)
		}
		close(in)
	}()
}

func Notify(payload string, ch *amqp.Channel) {
	err := ch.Publish(
		os.Getenv("RABBITMQ_DESTINATION_POSITIONS_EX"), // exchange
		os.Getenv("RABBITMQ_DESTINATION_ROUTING_KEY"),  // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Println("Message sent: ", payload)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %v", msg, err)
	}
}
