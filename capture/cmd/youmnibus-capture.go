package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	// Youtube setup
	yServ := getYoutubeServiceOrFail()

	// Generic RabbitMQ setup
	conn := getRabbitMQConnectionOrFail()
	ch := getRabbitMQChannelOrFail(conn)

	// Input setup
	inputQ := getRabbitMQQueueOrFail(ch, os.Getenv("INPUT_QUEUE_NAME"))
	inputCons := getRabbitMQConsumerOrFail(ch, inputQ, os.Getenv("CONSUMER_NAME"))

	// Run the consumer loop
	forever := make(chan bool)
	go func() {
		for d := range inputCons {
			channelID := string(d.Body)
			log.Printf("Received a message: %s", channelID)
			// Get channel data
			chData, err := retrieveChannelStatistics(yServ, channelID)
			logAndNackOnError(err, d, "Could not retrieve channel data for channel id %s: %s", channelID, err)

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func getYoutubeServiceOrFail() *youtube.Service {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(os.Getenv("YOUTUBE_API_KEY")))
	failOnError(err, "Could not create Youtube Service.")
	return youtubeService
}

func retrieveChannelStatistics(yServ *youtube.Service, channelId string) (*youtube.ChannelListResponse, error) {
	call := yServ.Channels.List("statistics")
	call = call.Id(channelId)
	return call.Do()
}
