package main

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	MongoHosts              string
	MongoPort               string
	MongoUser               string
	MongoPass               string
	MongoDatabase           string
	MongoCollection         string
	MemcacheSubscribersURLs string
	MemcacheViewsURLs       string
	MemcacheVideosURLs      string
	ServerPort              string
}

func GetConfig() *Config {
	return &Config{
		MongoHosts:              defaultIfEnvNil("MONGO_HOSTS", "localhost"),
		MongoPort:               defaultIfEnvNil("MONGO_PORT", "27017"),
		MongoUser:               defaultIfEnvNil("MONGO_USER", "youmnibus"),
		MongoPass:               defaultIfEnvNil("MONGO_PASS", "youmnibus"),
		MongoDatabase:           defaultIfEnvNil("MONGO_DATABASE", "youmnibus"),
		MongoCollection:         defaultIfEnvNil("MONGO_CAPTURES_COLLECTION", "captures"),
		MemcacheSubscribersURLs: defaultIfEnvNil("MEMCACHE_SUBSCRIBERS_URLS", "localhost:11211"),
		MemcacheViewsURLs:       defaultIfEnvNil("MEMCACHE_VIEWS_URLS", "localhost:11212"),
		MemcacheVideosURLs:      defaultIfEnvNil("MEMCACHE_VIDEOS_URLS", "localhost:11213"),
		ServerPort:              defaultIfEnvNil("SERVER_PORT", "8080"),
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
