package main

import (
	"log"
	"net/http"

	"github.com/bzhtux/sample_apps/rabbitmq/pkg/config"
	"github.com/bzhtux/sample_apps/rabbitmq/pkg/rmq"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := new(config.Conf)
	cfg.NewConfig()
	log.Printf("\033[32mLauching App %s with version %s on port %d\033[0m\n", cfg.App.Name, cfg.App.Version, cfg.App.Port)

	// Get a new RabbitMQ client
	clt, err := rmq.NewConn(cfg)
	if err != nil {
		log.Printf("--- Error Getting new RMQ client: %s\n", err.Error())
	}

	// Get a new RabbitMQ channel
	ch, err := rmq.NewChan(clt)
	if err != nil {
		log.Printf("--- Error Getting new RMQ channel: %s\n", err.Error())
	}

	// Set a new RabbitMQ handler
	// For more details see pkg/rmq/handler.go file
	rh := rmq.New(clt, ch, cfg.Database.Queue)

	// gin.SetMode(gin.ReleaseMode)
	// DebugMode should be used for dev only
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.MaxMultipartMemory = 16 << 32 // 16 MiB

	// Declare a RabbitMQ queue for sending message through
	q, err := ch.QueueDeclare(
		cfg.Database.Queue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Printf("--- Error setting RMQ Queue: %s\n", err.Error())
	}

	// Consume messages from the define defined above
	msgs, err := ch.Consume(
		q.Name, // name
		"",     // Consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("--- Error getting RMQ Queue: %s\n", err.Error())
	}

	// Healthcheck route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Ok",
			"message": "Alive",
		})
	})

	var forever chan struct{}
	router.POST("/msg", rh.SendMessage)
	go func() {
		for d := range msgs {

			log.Printf("New message received: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	router.Run(":8080")

	<-forever
}
