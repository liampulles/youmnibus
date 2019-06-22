package main

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	ChannelListFile    string
	QueueName          string
	DeadLetterExchange string
	AMQPURL            string
}

func GetConfig() *Config {
	return &Config{
		ChannelListFile:    defaultIfEnvNil("CHANNEL_IDS_FILE", "./channelIDs.txt"),
		QueueName:          defaultIfEnvNil("QUEUE", "channelsToFetch"),
		DeadLetterExchange: defaultIfEnvNil("QUEUE_DEAD_LETTER_EXCHANGE", "failedChannelsToFetch"),
		AMQPURL:            defaultIfEnvNil("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
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
