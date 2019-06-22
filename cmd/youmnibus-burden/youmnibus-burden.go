package main

import (
	"bufio"
	"log"
	"os"
	"time"

	yerror "github.com/liampulles/youmnibus/internal/error"
	"github.com/liampulles/youmnibus/internal/rabbitmq"
)

func main() {
	start := time.Now()

	// Get config
	conf := GetConfig()

	// RabbitMQ setup
	conn := rabbitmq.GetRabbitMQConnectionOrFail(conf.AMQPURL)
	defer conn.Close()
	ch := rabbitmq.GetRabbitMQChannelOrFail(conn)
	defer ch.Close()
	q := rabbitmq.GetRabbitMQQueueWithDeadLetterExchangeOrFail(ch, conf.QueueName, conf.DeadLetterExchange)

	file, err := os.Open(conf.ChannelListFile)
	yerror.FailOnError(err, "Failed to open file: "+conf.ChannelListFile)
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	// Just in case channel ids become enormous ? ;)
	const maxCapacity = 2048 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)
	for scanner.Scan() {
		channelID := scanner.Text()
		rabbitmq.PublishString(ch, q, channelID)
		count++
		log.Printf("Published channel ID: %s", channelID)
	}

	if err := scanner.Err(); err != nil {
		yerror.FailOnError(err, "Failed to read from scanner for file: "+conf.ChannelListFile)
	}
	elapsed := time.Since(start)
	log.Printf("Published %d channelIDs in %s. Now exiting.", count, elapsed)
}
