package main

import (
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/liampulles/youmnibus/internal/config"
	"github.com/liampulles/youmnibus/internal/mongo"
	"github.com/liampulles/youmnibus/internal/project"
	"github.com/liampulles/youmnibus/internal/rabbitmq"
	"github.com/liampulles/youmnibus/internal/youtube"
)

func main() {
	// Get config
	conf := GetConfigOrFail()
	log.Printf("API key: %s", conf.YouTubeAPIKey)

	// Youtube setup
	yServ := youtube.GetYoutubeServiceOrFail(conf.YouTubeAPIKey)

	// RabbitMQ setup
	conn := rabbitmq.GetRabbitMQConnectionOrFail(config.ConstructAMQPURL(conf.RabbitMQUsername, conf.RabbitMQPassword, conf.RabbitMQHost, conf.RabbitMQPort))
	defer conn.Close()
	inputCh := rabbitmq.GetRabbitMQChannelOrFail(conn)
	defer inputCh.Close()
	inputQ := rabbitmq.GetRabbitMQQueueWithDeadLetterExchangeOrFail(inputCh, conf.InputQueueName, conf.InputDeadLetterExchange)
	inputCons := rabbitmq.GetRabbitMQConsumerOrFail(inputCh, inputQ, conf.ConsumerName)

	// Mongo setup
	mClient := mongo.GetAndConnectMongoClientOrFail(config.MongoURL(conf.MongoHosts, conf.MongoPort, conf.MongoUser, conf.MongoPass))
	mColl := mongo.GetCollection(mClient, conf.MongoDatabase, conf.MongoCollection)

	// Memcache setup
	subsClient := project.GetMemcacheClient(conf.MemcacheSubscribersURLs)
	viewsClient := project.GetMemcacheClient(conf.MemcacheViewsURLs)
	videosClient := project.GetMemcacheClient(conf.MemcacheVideosURLs)

	// Run the consumer loop
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C / SIGKILL")
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

		// Invalidate the cache for the newly fetched channel - so that it will be regenerated.
		err = project.InvalidateCaches([]*memcache.Client{subsClient, viewsClient, videosClient}, channelID)
		if err != nil {
			log.Printf("Could not invalidate cache for channel id %s: %s", channelID, err)
			d.Nack(false, false)
			continue
		}
		log.Printf("Invalidated cache for channelId: %s", channelID)

		// Store the channel data in mongo
		_, err = mongo.StoreChannelData(mColl, chData, channelID, callTime)
		if err != nil {
			log.Printf("Could not store channel data in mongo for channel id %s: %s", channelID, err)
			d.Nack(false, false)
			continue
		}
		log.Printf("Stored channelId: %s", channelID)

		d.Ack(false)
	}
}
