package postgresql

import (
	"net/http"

	"github.com/bzhtux/sample_apps/postgresql/models"
	"github.com/gin-gonic/gin"
)

func (h Handler) AddNewBook(c *gin.Context) {
	var book = models.Books{}
	err := c.BindJSON(&book)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad request",
			"message": "Can not record a new book: " + err.Error(),
		})
	} else {
		r := h.DB.Where("Title = ?", book.Title).First(&book)
		if r.RowsAffected == 0 {
			h.DB.Create(&book)
			c.JSON(http.StatusAccepted, gin.H{
				"status":  "Accepted",
				"message": "New book has been recorded",
				"data": gin.H{
					"ID":     book.ID,
					"Title":  book.Title,
					"Author": book.Author,
				},
			})
		} else {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "Conflict",
				"message": "A Book already exists with title: " + book.Title,
				"data": gin.H{
					"ID": book.ID,
				},
			})
		}
	}
}
