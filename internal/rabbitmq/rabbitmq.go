package rabbitmq

import (
	yerror "github.com/liampulles/youmnibus/internal/error"

	"github.com/streadway/amqp"
)

func GetRabbitMQConnectionOrFail(amqpUrl string) *amqp.Connection {
	conn, err := amqp.Dial(amqpUrl)
	yerror.FailOnError(err, "Could not establish connection to RabbitMQ")
	return conn
}

func GetRabbitMQChannelOrFail(conn *amqp.Connection) *amqp.Channel {
	ch, err := conn.Channel()
	yerror.FailOnError(err, "Could not establish RabbitMQ channel")
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	yerror.FailOnError(err, "Could not establish QoS for RabbitMQ channel")
	return ch
}

func GetRabbitMQQueueWithDeadLetterExchangeOrFail(ch *amqp.Channel, queueName string, deadLetterExchangeName string) amqp.Queue {
	// Configure a dead-letter-queue
	err := ch.ExchangeDeclare(
		deadLetterExchangeName, // name
		"direct",               // kind
		true,                   // durable
		false,                  // delete when usused
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	yerror.FailOnError(err, "Could not establish RabbitMQ Dead Letter Exchange "+deadLetterExchangeName)

	args := make(amqp.Table, 1)
	args["x-dead-letter-exchange"] = deadLetterExchangeName
	return getRabbitMQQueueOrFail(ch, queueName, args)
}

func GetRabbitMQQueueOrFail(ch *amqp.Channel, queueName string) amqp.Queue {
	return getRabbitMQQueueOrFail(ch, queueName, make(amqp.Table))
}

func getRabbitMQQueueOrFail(ch *amqp.Channel, queueName string, args amqp.Table) amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)
	yerror.FailOnError(err, "Could not establish RabbitMQ queue "+queueName)
	return q
}

func GetRabbitMQConsumerOrFail(ch *amqp.Channel, q amqp.Queue, name string) <-chan amqp.Delivery {
	cons, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	yerror.FailOnError(err, "Failed to register RabbitMQ consumer "+name)
	return cons
}

func PublishString(ch *amqp.Channel, q amqp.Queue, body string) error {
	return ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
}
