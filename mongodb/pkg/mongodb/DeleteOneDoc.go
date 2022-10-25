package mongodb

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/bzhtux/sample_apps/mongodb/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h Handler) DeleteOnDocByName(c *gin.Context) {
	encodedbook := c.Params.ByName("book")
	b, err := url.QueryUnescape(encodedbook)
	if err != nil {
		log.Printf("Decoding %s error: %s", encodedbook, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "Internal Error",
			"message": "Can not decode " + encodedbook,
			"error":   err.Error(),
		})
	}

	var collection = models.MongoCollection{Database: "sampledb", Collection: "books"}
	col := h.clt.Database(collection.Database).Collection(collection.Collection)

	if _, err := col.DeleteOne(context.TODO(), bson.D{{Key: "Title", Value: b}}); err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means the query did not match any documents.
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "Not Found",
				"message": "Book titled " + b + " was not found.",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "Internal Error",
				"message": "Can not get doc from mongoDB",
				"error":   err.Error(),
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"message": "Doc deleted from mongoDB",
		})
	}
}
