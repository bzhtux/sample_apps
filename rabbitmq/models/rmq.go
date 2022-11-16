package models

type RmqConfig struct {
	Host     string `yaml:"host" envconfig:"RMQ_HOST"`
	Port     int    `yaml:"port" envconfig:"RMQ_PORT"`
	Username string `yaml:"username" envconfig:"RMQ_USER"`
	Password string `yaml:"password" envconfig:"RMQ_PASSWORD"`
	Queue    string `yaml:"queue" envconfig:"RMQ_QUEUE"`
	TLS      bool   `yaml:"tls" envconfig:"RMQ_TLS"`
}

// type RmqClient struct {
// 	queueName       string
// 	logger          *log.Logger
// 	connection      *amqp.Connection
// 	channel         *amqp.Channel
// 	done            chan bool
// 	notifyConnClose chan *amqp.Error
// 	notifyChanClose chan *amqp.Error
// 	notifyConfirm   chan amqp.Confirmation
// 	isReady         bool
// }

type Config interface {
	NewConfig()
}

type Message struct {
	Body string
}
