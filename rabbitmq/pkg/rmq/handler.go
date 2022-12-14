package rmq

import (
	"strconv"

	"github.com/bzhtux/sample_apps/rabbitmq/pkg/config"
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

func NewConn(cfg *config.Conf) (*amqp.Connection, error) {
	// RMQ uri schema:
	// amqp://<RMQ Username>:<RMQ Password>@<RMQ Hostname>:<RMQ Port>
	rmq_url := "amqp://" + cfg.Database.Username + ":" + cfg.Database.Password + "@" + cfg.Database.Host + ":" + strconv.Itoa(cfg.Database.Port)
	client, err := amqp.Dial(rmq_url)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewChan(clt *amqp.Connection) (*amqp.Channel, error) {
	ch, err := clt.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func NewQueue(cfg *config.Conf) string {
	return cfg.Database.Queue
}
