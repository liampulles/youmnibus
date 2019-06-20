package project

import (
	"encoding/json"
	"fmt"

	"github.com/liampulles/youmnibus/internal/mongo"

	"github.com/bradfitz/gomemcache/memcache"
	"google.golang.org/api/youtube/v3"
)

type SubscribersAt struct {
	SubscriberCount uint64 `json:"subscriberCount,omitempty,string"`
	At              string `json:"at,omitempty"`
}

type ViewsAt struct {
	ViewCount uint64 `json:"viewCount,omitempty,string"`
	At        string `json:"at,omitempty"`
}

type VideosAt struct {
	VideoCount uint64 `json:"videoCount,omitempty,string"`
	At         string `json:"at,omitempty"`
}

func GetMemcacheClient(memcacheURL string) *memcache.Client {
	return memcache.New(memcacheURL)
}

func MapSubscribersAt(channelID string, channelDatums []*mongo.ChannelData) ([]*SubscribersAt, error) {
	result := make([]*SubscribersAt, len(channelDatums))

	for i, channelData := range channelDatums {
		stats, err := GetStatisticsElement(channelData.Data, channelID)
		if err != nil {
			return nil, err
		}
		subsAt := &SubscribersAt{stats.SubscriberCount, channelData.Time}
		result[i] = subsAt
	}

	return result, nil
}

func MapViewsAt(channelID string, channelDatums []*mongo.ChannelData) ([]*ViewsAt, error) {
	result := make([]*ViewsAt, len(channelDatums))

	for i, channelData := range channelDatums {
		stats, err := GetStatisticsElement(channelData.Data, channelID)
		if err != nil {
			return nil, err
		}
		viewsAt := &ViewsAt{stats.ViewCount, channelData.Time}
		result[i] = viewsAt
	}

	return result, nil
}

func MapVideosAt(channelID string, channelDatums []*mongo.ChannelData) ([]*VideosAt, error) {
	result := make([]*VideosAt, len(channelDatums))

	for i, channelData := range channelDatums {
		stats, err := GetStatisticsElement(channelData.Data, channelID)
		if err != nil {
			return nil, err
		}
		videosAt := &VideosAt{stats.VideoCount, channelData.Time}
		result[i] = videosAt
	}

	return result, nil
}

func MarshalAndStore(memClient *memcache.Client, channelID string, data interface{}) ([]byte, error) {
	JSON, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return JSON, memClient.Set(&memcache.Item{Key: channelID, Value: JSON})
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
