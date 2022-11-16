package rmq

import (
	"context"
	"net/http"
	"time"

	"github.com/bzhtux/sample_apps/rabbitmq/models"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (h Handler) SendMessage(c *gin.Context) {
	var msg = models.Message{}
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "Bad request",
			"message": err.Error(),
		})
	} else {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		h.NewChan.PublishWithContext(ctx,
			"",         // exchange
			h.NewQueue, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg.Body),
			},
		)
		c.JSON(http.StatusAccepted, gin.H{
			"status":  "Accepted",
			"message": "New message sent",
			"data": gin.H{
				"body": msg.Body,
			},
		})
	}
}
