package main

import (
	"errors"
	"log"
	"os"

	yerror "github.com/liampulles/youmnibus/internal/error"
)

type Config struct {
	InputQueueName          string
	InputDeadLetterExchange string
	OutputQueueName         string
	ConsumerName            string
	RabbitMQUsername        string
	RabbitMQPassword        string
	RabbitMQHost            string
	RabbitMQPort            string
	YouTubeAPIKey           string
	MongoHosts              string
	MongoPort               string
	MongoDatabase           string
	MongoCollection         string
	MemcacheSubscribersURLs string
	MemcacheViewsURLs       string
	MemcacheVideosURLs      string
}

func GetConfigOrFail() *Config {
	youtubeAPIKey, err := errIfEnvNil("YOUTUBE_API_KEY")
	if err != nil {
		yerror.FailOnError(err, "Failed to get config")
	}

	return &Config{
		InputQueueName:          defaultIfEnvNil("INPUT_QUEUE", "channelsToFetch"),
		InputDeadLetterExchange: defaultIfEnvNil("INPUT_DEAD_LETTER_EXCHANGE", "failedChannelsToFetch"),
		OutputQueueName:         defaultIfEnvNil("OUTPUT_QUEUE", "fetchedChannels"),
		ConsumerName:            defaultIfEnvNil("NAME", ""),
		RabbitMQUsername:        defaultIfEnvNil("RABBITMQ_USERNAME", "guest"),
		RabbitMQPassword:        defaultIfEnvNil("RABBITMQ_PASSWORD", "guest"),
		RabbitMQHost:            defaultIfEnvNil("RABBITMQ_HOST", "localhost"),
		RabbitMQPort:            defaultIfEnvNil("RABBITMQ_PORT", "5672"),
		YouTubeAPIKey:           youtubeAPIKey,
		MongoHosts:              defaultIfEnvNil("MONGO_HOSTS", "localhost"),
		MongoPort:               defaultIfEnvNil("MONGO_PORT", "27017"),
		MongoDatabase:           defaultIfEnvNil("MONGO_DATABASE", "youmnibus"),
		MongoCollection:         defaultIfEnvNil("MONGO_CAPTURES_COLLECTION", "captures"),
		MemcacheSubscribersURLs: defaultIfEnvNil("MEMCACHE_SUBSCRIBERS_URLS", "localhost:11211"),
		MemcacheViewsURLs:       defaultIfEnvNil("MEMCACHE_VIEWS_URLS", "localhost:11212"),
		MemcacheVideosURLs:      defaultIfEnvNil("MEMCACHE_VIDEOS_URLS", "localhost:11213"),
	}
}

func defaultIfEnvNil(env string, def string) string {
	param, set := os.LookupEnv(env)
	if !set {
		log.Printf("Environment variable %s not defined. Using default: %s", env, def)
		return def
	}
	return param
}

func errIfEnvNil(env string) (string, error) {
	param, set := os.LookupEnv(env)
	if !set {
		return "", errors.New("Environment variable " + env + " not defined.")
	}
	return param, nil
}
