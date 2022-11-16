package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bzhtux/sample_apps/rabbitmq/models"
	"github.com/bzhtux/servicebinding/bindings"
	"github.com/kelseyhightower/envconfig"
	ms "github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_CONFIG_DIR = "/config"
	// Export from the shell to get a local config dir
	// DEFAULT_CONFIG_DIR = "/Users/yfoeillet/go/src/github.com/bzhtux/sample_apps/rabbitmq/config"
	rmqConfigDir  = "RMQ_CONFIG_DIR"
	rmqConfigFile = "rmq.yaml"
	// If you want to enforce queueName, define the enforced queue below
	rmqQueueName = "vmwarenforced"
)

type RMQConfigDir struct {
	Root string
}

type Config interface {
	NewConfig()
}

type RMQConfig models.RmqConfig

func GetConfigDir() *RMQConfigDir {
	CONFIG_DIR, exists := os.LookupEnv(rmqConfigDir)
	rmqcd := &RMQConfigDir{}
	if !exists {
		rmqcd.Root = DEFAULT_CONFIG_DIR
	} else {
		rmqcd.Root = CONFIG_DIR
	}
	return rmqcd
}

func (rc *RMQConfig) LoadConfigFromFile() error {

	rmqcd := GetConfigDir()
	cfg, err := os.Open(filepath.Join(rmqcd.Root, rmqConfigFile))
	if err != nil {
		log.Printf("Error opening RMQ config file: %s", err.Error())
		return err
	}
	defer cfg.Close()
	d := yaml.NewDecoder(cfg)

	if err := d.Decode(&rc); err != nil {
		log.Printf("--- Error decoding YAML: %s (%v)", err.Error(), d)
		return err
	}

	return nil
}

func (rc *RMQConfig) LoadConfigFromEnv() {
	envconfig.Process("", rc)
}

func (rc *RMQConfig) LoadConfigFromBindings(t string) error {

	b, err := bindings.NewBinding(t)
	if err != nil {
		log.Printf("Error while getting bindings: %s", err.Error())
		return err
	}

	if err := ms.Decode(b, &rc); err != nil {
		return err
	}
	return nil
}

func (rc *RMQConfig) NewConfig() *RMQConfig {
	rc.LoadConfigFromFile()
	rc.LoadConfigFromBindings("rabbitmq")
	log.Printf("*** NewConfig.rc.Queue: %s\n", rc.Queue)
	if rc.Queue == "" {
		log.Printf("In if statetement ...")
		rc.Queue = rmqQueueName
	}
	rc.LoadConfigFromEnv()
	return rc
}
