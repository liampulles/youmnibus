package main

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	MongoURL               string
	MongoDatabase          string
	MongoCollection        string
	MemcacheSubscribersURL string
	MemcacheViewsURL       string
	MemcacheVideosURL      string
	ServerPort             string
}

func GetConfig() *Config {
	return &Config{
		MongoURL:               defaultIfEnvNil("MONGO_URL", "mongodb://localhost:27017"),
		MongoDatabase:          defaultIfEnvNil("MONGO_DATABASE", "youmnibus"),
		MongoCollection:        defaultIfEnvNil("MONGO_CAPTURES_COLLECTION", "captures"),
		MemcacheSubscribersURL: defaultIfEnvNil("MEMCACHE_SUBSCRIBERS_URL", "localhost:11211"),
		MemcacheViewsURL:       defaultIfEnvNil("MEMCACHE_VIEWS_URL", "localhost:11212"),
		MemcacheVideosURL:      defaultIfEnvNil("MEMCACHE_VIDEOS_URL", "localhost:11213"),
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
