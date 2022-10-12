package models

import (
	"log"
	"os"

	"github.com/bzhtux/servicebinding/bindings"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG_DIR = "/config"
)

type RedisConfig struct {
	Host       string `yaml:"host" envconfig:"REDIS_HOST"`
	Port       int    `yaml:"port" envconnfig:"REDIS_PORT"`
	Username   string `yaml:"username" envconfig:"REDIS_USERNAME"`
	Password   string `yaml:"password" envconfig:"REDIS_PASSWORD"`
	DB         int    `yaml:"database" envconfig:"REDIS_DB"`
	SSLenabled bool   `yaml:"sslenabled" envconfig:"REDIS_SSL"`
}

type Config interface {
	NewConfig()
}

type Messages struct {
	Key   string
	Value string
}

func (rc *RedisConfig) GetConfigFromFile() {
	file, err := os.Open(CONFIG_DIR + "/redis.yaml")
	if err != nil {
		log.Println("Error loading config file ...", err.Error())
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&rc); err != nil {
		log.Println("Error loading config file ...", err.Error())
	}
}

func (rc *RedisConfig) GetConfigFromEnv() {
	envconfig.Process("", rc)
}

func (rc *RedisConfig) GetConfigFromBindings(t string) {
	b, err := bindings.NewBinding(t)
	if err != nil {
		log.Printf("Error while getting bindings: %s", err.Error())
	}
	rc.Host = b.Host
	// fmt.Printf("Redis Host: %s\n", b.Host)
	rc.Port = int(b.Port)
	// fmt.Printf("Redis Port: %d\n", b.Port)
	rc.Username = b.Username
	// fmt.Printf("Redis Username: %s\n", b.Username)
	rc.Password = b.Password
	// fmt.Printf("Redis Password: %s\n", b.Password)
}

func (rc *RedisConfig) NewConfig() *RedisConfig {
	rc.GetConfigFromFile()
	rc.GetConfigFromBindings("redis")
	rc.GetConfigFromEnv()
	return rc
}
