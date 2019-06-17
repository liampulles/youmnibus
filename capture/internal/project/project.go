package project

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"google.golang.org/api/youtube/v3"
)

type SubscribersAt struct {
	SubscriberCount uint64 `json:"subscriberCount,omitempty,string"`
	At              string `json:"at,omitempty,string"`
}

type ViewsAt struct {
	ViewCount uint64 `json:"viewCount,omitempty,string"`
	At        string `json:"at,omitempty,string"`
}

type VideosAt struct {
	VideoCount uint64 `json:"videoCount,omitempty,string"`
	At         string `json:"at,omitempty,string"`
}

func GetMemcacheClient(memcacheURL string) *memcache.Client {
	return memcache.New(memcacheURL)
}

func ProjectSubscriberData(memClient *memcache.Client, channelId string, at time.Time, stats *youtube.ChannelStatistics) ([]*SubscribersAt, error) {
	// Get the statistics element in the channel data
	subsAt := getSubscribersAtTime(stats, at)

	// Retrieve the curently stored projection, or new.
	var storedStats []*SubscribersAt
	err := getArray(memClient, channelId, storedStats)
	if err != nil {
		return nil, err
	}

	// Add to the projection
	newStats := append(storedStats, subsAt)

	// Store the projection
	return storedStats, marshalAndStore(memClient, channelId, newStats)
}

func ProjectViewsData(memClient *memcache.Client, channelId string, at time.Time, stats *youtube.ChannelStatistics) ([]*ViewsAt, error) {
	// Get the statistics element in the channel data
	viewsAt := getViewsAtTime(stats, at)

	// Retrieve the curently stored projection, or new.
	var storedStats []*ViewsAt
	err := getArray(memClient, channelId, storedStats)
	if err != nil {
		return nil, err
	}

	// Add to the projection
	newStats := append(storedStats, viewsAt)

	// Store the projection
	return storedStats, marshalAndStore(memClient, channelId, newStats)
}

func ProjectVideosData(memClient *memcache.Client, channelId string, at time.Time, stats *youtube.ChannelStatistics) ([]*VideosAt, error) {
	// Get the statistics element in the channel data
	videosAt := getVideosAtTime(stats, at)

	// Retrieve the curently stored projection, or new.
	var storedStats []*VideosAt
	err := getArray(memClient, channelId, storedStats)
	if err != nil {
		return nil, err
	}

	// Add to the projection
	newStats := append(storedStats, videosAt)

	// Store the projection
	return storedStats, marshalAndStore(memClient, channelId, newStats)
}

func subscribersAtEqual(subsAt1 *SubscribersAt, subsAt2 *SubscribersAt) bool {
	return subsAt1 == subsAt2 ||
		(subsAt1 != nil &&
			subsAt2 != nil &&
			subsAt1.SubscriberCount == subsAt2.SubscriberCount &&
			subsAt1.At == subsAt2.At)
}

func marshalAndStore(memClient *memcache.Client, channelID string, data interface{}) error {
	JSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return memClient.Set(&memcache.Item{Key: channelID, Value: JSON})
}

func getArray(memClient *memcache.Client, channelId string, arrayPointer interface{}) error {
	// Retrieve the curently stored projection, or new.
	item, err := memClient.Get(channelId)
	if err != nil {
		if err != memcache.ErrCacheMiss {
			return err
		}
		log.Printf("Got cache miss for channel Id " + channelId)
		return nil
	}
	err = json.Unmarshal(item.Value, arrayPointer)
	if err != nil {
		return err
	}
	return nil
}

func InvalidateCaches(memClients []*memcache.Client, identifier string) error {
	for _, memClient := range memClients {
		err := memClient.Delete(identifier)
		if err != nil && err != memcache.ErrCacheMiss {
			return err
		}
	}
	return nil
}

func GetStatisticsElement(chData *youtube.ChannelListResponse, channelID string) (*youtube.ChannelStatistics, error) {
	channels := chData.Items
	if len(channels) != 1 {
		return nil, fmt.Errorf("Expected 1 channel for %s but got %d", channelID, len(channels))
	}
	return channels[0].Statistics, nil
}

func getSubscribersAtTime(stats *youtube.ChannelStatistics, at time.Time) *SubscribersAt {
	return &SubscribersAt{stats.SubscriberCount, at.Format(time.RFC3339)} // Profile of ISO8601
}

func getViewsAtTime(stats *youtube.ChannelStatistics, at time.Time) *ViewsAt {
	return &ViewsAt{stats.ViewCount, at.Format(time.RFC3339)} // Profile of ISO8601
}

func getVideosAtTime(stats *youtube.ChannelStatistics, at time.Time) *VideosAt {
	return &VideosAt{stats.ViewCount, at.Format(time.RFC3339)} // Profile of ISO8601
}
