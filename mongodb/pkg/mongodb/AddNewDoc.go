package mongodb

import (
	"context"
	"net/http"

	"github.com/bzhtux/sample_apps/mongodb/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h Handler) AddNewDoc(c *gin.Context) {
	var book = models.Books{}
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad request",
			"message": "Can not add a new doc :(",
			"error":   err.Error(),
		})
	} else {
		if h.DocExists(book.Title) {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "Conflict",
				"message": "Book " + book.Title + " already exists.",
			})
		} else {
			var collection = models.MongoCollection{Database: "sampledb", Collection: "books"}
			col := h.clt.Database(collection.Database).Collection(collection.Collection)
			var doc = bson.D{{Key: "Title", Value: book.Title}, {Key: "Author", Value: book.Author}}
			result, err := col.InsertOne(context.TODO(), doc)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Server Error",
					"message": "New book " + book.Title + " by " + book.Author + " was not added to " + collection.Collection + "' collection.",
					"error":   err.Error(),
				})
			} else {
				c.JSON(http.StatusCreated, gin.H{
					"status":  "OK",
					"message": "New book added to " + collection.Collection + "' collection",
					"data": gin.H{
						"Book title":  book.Title,
						"Book Author": book.Author,
						"ID":          result.InsertedID,
						"result":      result,
					},
				})
			}

		}
	}
}
