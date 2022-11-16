package rmq

import (
	"strconv"

	"github.com/bzhtux/sample_apps/rabbitmq/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewConn() (*amqp.Connection, error) {
	// Define a new connection to RabbitMQ server

	rc := new(config.RMQConfig)
	rc.NewConfig()
	// RMQ uri schema:
	// amqp://<RMQ Username>:<RMQ Password>@<RMQ Hostname>:<RMQ Port>
	rmq_url := "amqp://" + rc.Username + ":" + rc.Password + "@" + rc.Host + ":" + strconv.Itoa(rc.Port)
	client, err := amqp.Dial(rmq_url)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewChan(clt *amqp.Connection) (*amqp.Channel, error) {
	// Define a new channel for the RabbitMQ connection to server

	ch, err := clt.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func NewQueue() string {
	// Define a new queue to use for sending and receiving messages

	rc := new(config.RMQConfig)
	rc.NewConfig()
	return rc.Queue
}
