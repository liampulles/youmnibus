package youtube

import (
	"context"

	yerror "github.com/liampulles/youmnibus/capture/internal/error"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetYoutubeServiceOrFail(youtubeApiKey string) *youtube.Service {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(youtubeApiKey))
	yerror.FailOnError(err, "Could not create Youtube Service.")
	return youtubeService
}

func RetrieveChannelStatistics(yServ *youtube.Service, channelId string) (*youtube.ChannelListResponse, error) {
	call := yServ.Channels.List("statistics")
	call = call.Id(channelId)
	return call.Do()
}
