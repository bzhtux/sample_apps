package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	clt *mongo.Client
}

func New(clt *mongo.Client) Handler {
	client, _ := NewClient()
	return Handler{client}
}
