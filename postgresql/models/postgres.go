package models

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const (
	CONFIG_DIR = "/config"
)

type PGConfig struct {
	Host       string `yaml:"host" envconfig:"PG_HOST"`
	Port       int    `yaml:"port" envconnfig:"PG_PORT"`
	Username   string `yaml:"username" envconfig:"PG_USERNAME"`
	Password   string `yaml:"password" envconfig:"PG_PASSWORD"`
	DB         string `yaml:"database" envconfig:"PG_DB"`
	SSLenabled bool   `yaml:"sslenabled" envconfig:"PG_SSL"`
}

type Config interface {
	NewConfig()
}

type Books struct {
	*gorm.Model
	ID     uint   `gorm:"primaryKey"`
	Title  string `gorm:"not null" json:"title" binding:"required"`
	Author string `gorm:"not null" json:"author" binding:"required"`
}

func (cfg *PGConfig) GetConfigFromFile() {
	file, err := os.Open(CONFIG_DIR + "/postgres.yaml")
	if err != nil {
		log.Println("Error loading config file ...", err.Error())
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&cfg); err != nil {
		log.Println("Error loading config file ...", err.Error())
	}
}

func (cfg *PGConfig) GetConfigFromEnv() {
	envconfig.Process("", cfg)
}

func (cfg *PGConfig) NewConfig() *PGConfig {
	cfg.GetConfigFromFile()
	cfg.GetConfigFromEnv()
	return cfg
}
