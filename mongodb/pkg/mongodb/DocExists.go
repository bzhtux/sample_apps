package mongodb

import (
	"context"

	"github.com/bzhtux/sample_apps/mongodb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h Handler) DocExists(dt string) bool {
	var collection = models.MongoCollection{Database: "sampledb", Collection: "books"}
	col := h.clt.Database(collection.Database).Collection(collection.Collection)
	var res bson.M
	if err := col.FindOne(context.TODO(), bson.D{{Key: "Title", Value: dt}}).Decode(&res); err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means the query did not match any documents.
			return false
		} else {
			return true
		}
	}
	return true
}
