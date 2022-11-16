package rmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler struct {
	NewConn  *amqp.Connection
	NewChan  *amqp.Channel
	NewQueue string
}

func New(newconn *amqp.Connection, newchan *amqp.Channel, q string) Handler {
	return Handler{newconn, newchan, q}
}

func (h Handler) Close() error {
	if h.NewConn == nil {
		return nil
	}
	return h.NewConn.Close()
}
