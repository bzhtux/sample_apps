package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bzhtux/sample_apps/rabbitmq/models"
	bindings "github.com/bzhtux/servicebinding"
	"github.com/jinzhu/copier"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

var (
	EnvConfigDir = strings.ToUpper(AppName) + "_CONFIG_DIR"
)

const (
	DEFAULT_CONFIG_DIR    = "/config"
	ConfigFile            = "config.yml"
	AppName               = "gorabbit"
	AppDesc               = "Golang RabbitMQ App for demo purpose"
	EnvServiceBindingRoot = "SERVICE_BINDING_ROOT"
	// Export from the shell to get a local config dir
	// DEFAULT_CONFIG_DIR = "/Users/yfoeillet/go/src/github.com/bzhtux/sample_apps/rabbitmq/config"
	rmqConfigDir  = "RMQ_CONFIG_DIR"
	rmqConfigFile = "rmq.yaml"
	// If you want to enforce queueName, define the enforced queue below
	rmqQueueName = "vmwarenforced"
)

type Conf models.Config

func (cfg *Conf) GetConfigDir() *Conf {
	cfg_dir, exists := os.LookupEnv(rmqConfigDir)
	if !exists {
		cfg.Dir.Root = DEFAULT_CONFIG_DIR
	} else {
		cfg.Dir.Root = cfg_dir
	}
	return cfg
}

func (cfg *Conf) LoadConfigFromFile() error {
	config := cfg.GetConfigDir()
	c, err := os.Open(filepath.Join(config.Dir.Root, ConfigFile))
	if err != nil {
		log.Printf("Error loading config file: %s", err.Error())
		return err
	}
	defer c.Close()

	d := yaml.NewDecoder(c)

	if err := d.Decode(&cfg); err != nil {
		log.Printf("--- Error decoding config file: %s ", err.Error())
		return err
	}

	return nil
}

// func (rc *RMQConfig) LoadConfigFromEnv() {
// 	envconfig.Process("", rc)
// }

func (cfg *Conf) LoadConfigFromBindings(t string) error {
	b, err := bindings.NewBinding(t)
	if err != nil {
		log.Printf("Error while getting bindings: %s", err.Error())
		return err
	}

	// if err := ms.Decode(b, &rc); err != nil {
	// 	return err
	// }
	copier.Copy(&cfg.Database, &b)
	return nil
}

func (cfg *Conf) NewConfig() {
	cfg.LoadConfigFromFile()
	cfg.LoadConfigFromBindings("rabbitmq")
	log.Printf("*** NewConfig.rc.Queue: %s\n", cfg.Database.Queue)
	if cfg.Database.Queue == "" {
		log.Printf("In if statetement ...")
		cfg.Database.Queue = rmqQueueName
	}
	if err := envconfig.Process("", &cfg.Database); err != nil {
		log.Printf("Envconfig.process error: %s", err.Error())
	}
}
