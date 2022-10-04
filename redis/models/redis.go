package models

import (
	"log"
	"os"

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

func (rc *RedisConfig) NewConfig() *RedisConfig {
	rc.GetConfigFromFile()
	// fmt.Println("Username from file: ", rc.Username)
	// fmt.Println(rc)
	rc.GetConfigFromEnv()
	// fmt.Println("Username from env: ", os.Getenv("REDIS_USERNAME"))
	// fmt.Println("Username used for RedisConfig: ", rc.Username)
	// fmt.Println("SSL used for RedisConfig: ", rc.SSLenabled)
	// fmt.Println(reflect.TypeOf(rc.SSLenabled))
	return rc
}
