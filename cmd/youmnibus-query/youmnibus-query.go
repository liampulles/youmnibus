package main

import (
	"log"
	"net/http"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"github.com/liampulles/youmnibus/internal/config"
	"github.com/liampulles/youmnibus/internal/mongo"
	"github.com/liampulles/youmnibus/internal/project"
)

func main() {
	// Get config
	conf := GetConfig()

	// Mongo setup
	mClient := mongo.GetAndConnectMongoClientOrFail(config.MongoURL(conf.MongoHosts, conf.MongoPort, conf.MongoUser, conf.MongoPass))
	mColl := mongo.GetCollection(mClient, conf.MongoDatabase, conf.MongoCollection)

	// Memcache setup
	subsClient := project.GetMemcacheClient(conf.MemcacheSubscribersURLs)
	viewsClient := project.GetMemcacheClient(conf.MemcacheViewsURLs)
	videosClient := project.GetMemcacheClient(conf.MemcacheVideosURLs)

	// Creates a gin router with default middleware:
	router := gin.Default()

	router.GET("/subscribers/:channelID", func(c *gin.Context) {
		channelID := c.Params.ByName("channelID")
		item, err := subsClient.Get(channelID)
		var result []byte
		if err == nil {
			result = item.Value
		} else if err == memcache.ErrCacheMiss {
			// Get data from Mongo
			channelData, err := mongo.RetrieveChannelData(mColl, channelID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				log.Printf("Error while calling mongo: %s", err)
				return
			}
			// Project the data
			subsAt, err := project.MapSubscribersAt(channelID, channelData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				log.Printf("Error while projecting data: %s", err)
				return
			}
			// Cache the projected data
			result, err = project.MarshalAndStore(subsClient, channelID, subsAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				log.Printf("Error while caching data: %s", err)
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			log.Printf("Error while calling cache: %s", err)
			return
		}
		// Return the data
		c.Data(http.StatusOK, "application/json", result)
	})

	router.GET("/views/:channelID", func(c *gin.Context) {
		channelID := c.Params.ByName("channelID")
		item, err := viewsClient.Get(channelID)
		var result []byte
		if err == nil {
			result = item.Value
		} else if err == memcache.ErrCacheMiss {
			// Get data from Mongo
			channelData, err := mongo.RetrieveChannelData(mColl, channelID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// Project the data
			viewsAt, err := project.MapViewsAt(channelID, channelData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			// Cache the projected data
			result, err = project.MarshalAndStore(viewsClient, channelID, viewsAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Return the data
		c.Data(http.StatusOK, "application/json", result)
	})

	router.GET("/videos/:channelID", func(c *gin.Context) {
		channelID := c.Params.ByName("channelID")
		item, err := videosClient.Get(channelID)
		var result []byte
		if err == nil {
			result = item.Value
		} else if err == memcache.ErrCacheMiss {
			// Get data from Mongo
			channelData, err := mongo.RetrieveChannelData(mColl, channelID)
			if err != nil {

				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			// Project the data
			videosAt, err := project.MapVideosAt(channelID, channelData)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			// Cache the projected data
			result, err = project.MarshalAndStore(videosClient, channelID, videosAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		// Return the data
		c.Data(http.StatusOK, "application/json", result)
	})

	router.Run(":" + conf.ServerPort)
}
