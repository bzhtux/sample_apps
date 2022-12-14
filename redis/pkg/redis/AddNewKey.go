package redis

import (
	"log"
	"net/http"

	"github.com/bzhtux/sample_apps/redis/models"
	"github.com/gin-gonic/gin"
)

func (h Handler) AddNewKey(c *gin.Context) {
	var msg = models.Message{}
	err := c.BindJSON(&msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad request",
			"message": "Can not record new key: " + err.Error(),
		})
	} else {
		client := h.NewPool.Get()
		defer client.Close()
		if h.KeyExists(msg.Key) {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "Conflict",
				"message": "Key already exists: " + msg.Key,
			})
		} else {
			_, err := client.Do("SET", msg.Key, msg.Value)
			if err != nil {
				log.Println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "Server Error",
					"message": "New key " + msg.Key + " = " + msg.Value + " not recorded",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"status":  "OK",
					"message": "New key has been recorded successfuly",
					"data": gin.H{
						"key":   msg.Key,
						"value": msg.Value,
					},
				})
			}
		}
	}
}
