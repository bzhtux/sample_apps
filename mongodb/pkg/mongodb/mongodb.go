package mongodb

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/bzhtux/sample_apps/mongodb/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient() (*mongo.Client, error) {
	mc := new(config.MongoConfig)
	mc.NewConfig()
	var uri = "mongodb://" + mc.Username + ":" + mc.Password + "@" + mc.Host + ":" + strconv.Itoa(int(mc.Port)) + "/?maxPoolSize=20&retryWrites=true&w=majority"
	client, err := mongo.Connect(context.TODO(), options.Client().SetConnectTimeout(3*time.Second).SetServerSelectionTimeout(3*time.Second).ApplyURI(uri))
	if err != nil {
		log.Printf("--- Error Connecting to MongoDB instance: %s\n" + err.Error())
		return nil, err
	}
	return client, nil
}

func NewCollection(clt *mongo.Client) (*mongo.Collection, error) {
	return clt.Database("sampledb").Collection("book"), nil
}
