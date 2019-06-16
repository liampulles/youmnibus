package rabbitmq

import (
	"log"
	"os"
	yerror "youmnibus-capture/internal/error"

	"github.com/streadway/amqp"
)

func GetRabbitMQConnectionOrFail() *amqp.Connection {
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	yerror.FailOnError(err, "Could not establish connection to RabbitMQ")
	return conn
}

func GetRabbitMQChannelOrFail(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	yerror.FailOnError(err, "Could not establish RabbitMQ channel")
	return ch
}

func GetRabbitMQQueueOrFail(ch *amqp.Channel, name string) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	yerror.FailOnError(err, "Could not establish RabbitMQ queue "+name)
	return q
}

func GetRabbitMQConsumerOrFail(ch *amqp.Channel, q amqp.Queue, name string) <-chan amqp.Delivery {
	cons, err := ch.Consume(
		q.Name, // queue
		name,   // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	yerror.FailOnError(err, "Failed to register RabbitMQ consumer "+name)
	return cons
}

func PublishJSON(ch *amqp.Channel, q amqp.Queue, jsonBytes []byte) error {
	return ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBytes,
		})
}

func LogAndNackOnError(err error, d amqp.Delivery, msgF string, args ...interface{}) {
	if err != nil {
		d.Nack(false, false)
		log.Printf(msgF, args)
	}
}
