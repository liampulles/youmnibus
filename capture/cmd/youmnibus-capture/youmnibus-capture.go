package main

import (
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/liampulles/youmnibus/capture/internal/mongo"
	"github.com/liampulles/youmnibus/capture/internal/project"
	"github.com/liampulles/youmnibus/capture/internal/rabbitmq"
	"github.com/liampulles/youmnibus/capture/internal/youtube"
)

func main() {
	// Get config
	conf := GetConfigOrFail()

	// Youtube setup
	yServ := youtube.GetYoutubeServiceOrFail(conf.YouTubeAPIKey)

	// RabbitMQ setup
	conn := rabbitmq.GetRabbitMQConnectionOrFail(conf.AMQPURL)
	defer conn.Close()
	inputCh := rabbitmq.GetRabbitMQChannelOrFail(conn)
	defer inputCh.Close()
	inputQ := rabbitmq.GetRabbitMQQueueWithDeadLetterExchangeOrFail(inputCh, conf.InputQueueName, conf.InputDeadLetterExchange)
	inputCons := rabbitmq.GetRabbitMQConsumerOrFail(inputCh, inputQ, conf.ConsumerName)

	// Mongo setup
	mClient := mongo.GetAndConnectMongoClientOrFail(conf.MongoURL)
	mColl := mongo.GetCollection(mClient, conf.MongoDatabase, conf.MongoCollection)

	// Memcache setup
	subsClient := project.GetMemcacheClient(conf.MemcacheSubscribersURL)
	viewsClient := project.GetMemcacheClient(conf.MemcacheViewsURL)
	videosClient := project.GetMemcacheClient(conf.MemcacheVideosURL)

	// Run the consumer loop
	forever := make(chan bool)
	go func() {
		for d := range inputCons {
			channelID := string(d.Body)
			log.Printf("Received a message: %s", channelID)
			// Get channel data
			callTime := time.Now()
			chData, err := youtube.RetrieveChannelStatistics(yServ, channelID)
			if err != nil {
				log.Printf("Could not retrieve channel data for channel id %s: %s", channelID, err)
				d.Nack(false, false)
				continue
			}
			log.Printf("Fetched channel Id (%d ms): %s", (time.Since(callTime).Nanoseconds())/1000000, channelID)

			// Validate the response
			_, err = project.GetStatisticsElement(chData, channelID)
			if err != nil {
				log.Printf("Validation on response for channel id %s failed: %s", channelID, err)
				d.Nack(false, false)
				continue
			}

			// Store the channel data in mongo
			insertRes, err := mongo.StoreChannelData(mColl, chData, channelID, callTime)
			if err != nil {
				log.Printf("Could not store channel data in mongo for channel id %s: %s", channelID, err)
				d.Nack(false, false)
				continue
			}
			log.Printf("Stored channelId: %s", channelID)

			// Invalidate the cache for the newly fetched channel - so that it will be regenerated.
			err = project.InvalidateCaches([]*memcache.Client{subsClient, viewsClient, videosClient}, channelID)
			if err != nil {
				log.Printf("Could not invalidate cache for channel id %s: %s", channelID, err)
				err = mongo.RollbackInsertion(mColl, insertRes)
				d.Nack(false, false)
				if err != nil {
					log.Printf("Could not remove data in mongo for channel id %s: %s", channelID, err)
				}
				continue
			}
			log.Printf("Invalidated cache for channelId: %s", channelID)

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C / SIGKILL")
	<-forever
}
