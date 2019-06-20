package youtube

import (
	"context"

	yerror "github.com/liampulles/youmnibus/internal/error"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetYoutubeServiceOrFail(youtubeAPIKey string) *youtube.Service {
	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))
	yerror.FailOnError(err, "Could not create Youtube Service.")
	return youtubeService
}

func RetrieveChannelStatistics(yServ *youtube.Service, channelID string) (*youtube.ChannelListResponse, error) {
	call := yServ.Channels.List("statistics")
	call = call.Id(channelID)
	return call.Do()
}
