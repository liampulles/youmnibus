package mongo

import (
	"context"
	"time"

	yerror "github.com/liampulles/youmnibus/capture/internal/error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/youtube/v3"
)

const mongoTimeout = 10 * time.Second

type ChannelData struct {
	ChannelID string
	Time      string
	Data      *youtube.ChannelListResponse
}

func GetAndConnectMongoClientOrFail(mongoURL string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	yerror.FailOnError(err, "Could not create Mongo client")
	ctx, cancel := context.WithTimeout(context.Background(), mongoTimeout)
	err = client.Connect(ctx)
	cancel()
	yerror.FailOnError(err, "Could not connect Mongo client")
	return client
}

func GetCollection(client *mongo.Client, db string, col string) *mongo.Collection {
	return client.Database(db).Collection(col)
}

func StoreChannelData(mColl *mongo.Collection, chData *youtube.ChannelListResponse, channelId string, callTime time.Time) (*mongo.InsertOneResult, error) {
	// Add some additional data to store
	toStore := ChannelData{channelId, callTime.Format(time.RFC3339Nano), chData}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	id, err := mColl.InsertOne(ctx, toStore)
	cancel()
	return id, err
}

func RollbackInsertion(mColl *mongo.Collection, mRes *mongo.InsertOneResult) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := mColl.DeleteOne(ctx, bson.D{{"_id", bson.D{{"$eq", mRes.InsertedID}}}})
	cancel()
	return err
}
