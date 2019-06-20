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
	AMQPURL                 string
	YouTubeAPIKey           string
	MongoURL                string
	MongoDatabase           string
	MongoCollection         string
	MemcacheSubscribersURL  string
	MemcacheViewsURL        string
	MemcacheVideosURL       string
}

func GetConfigOrFail() *Config {
	defaultName, err := os.Hostname()
	if err != nil {
		log.Printf("Could not determine hostname: %s", err)
		defaultName = "worker"
	}

	youtubeAPIKey, err := errIfEnvNil("YOUTUBE_API_KEY")
	if err != nil {
		yerror.FailOnError(err, "Failed to get config")
	}

	return &Config{
		InputQueueName:          defaultIfEnvNil("INPUT_QUEUE", "channelsToFetch"),
		InputDeadLetterExchange: defaultIfEnvNil("INPUT_DEAD_LETTER_EXCHANGE", "failedChannelsToFetch"),
		OutputQueueName:         defaultIfEnvNil("OUTPUT_QUEUE", "fetchedChannels"),
		ConsumerName:            defaultIfEnvNil("NAME", defaultName),
		AMQPURL:                 defaultIfEnvNil("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
		YouTubeAPIKey:           youtubeAPIKey,
		MongoURL:                defaultIfEnvNil("MONGO_URL", "mongodb://localhost:27017"),
		MongoDatabase:           defaultIfEnvNil("MONGO_DATABASE", "youmnibus"),
		MongoCollection:         defaultIfEnvNil("MONGO_CAPTURES_COLLECTION", "captures"),
		MemcacheSubscribersURL:  defaultIfEnvNil("MEMCACHE_SUBSCRIBERS_URL", "localhost:11211"),
		MemcacheViewsURL:        defaultIfEnvNil("MEMCACHE_VIEWS_URL", "localhost:11212"),
		MemcacheVideosURL:       defaultIfEnvNil("MEMCACHE_VIDEOS_URL", "localhost:11213"),
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
