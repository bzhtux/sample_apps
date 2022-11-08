package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bzhtux/sample_apps/mongodb/models"
	"github.com/bzhtux/servicebinding/bindings"
	"github.com/kelseyhightower/envconfig"
	. "github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_CONFIG_DIR = "/config"
	mongoConfigDir     = "MONGO_CONFIG_DIR"
	mongoConfigFile    = "mongo.yaml"
)

type MongoConfigDir struct {
	root string
}

type Config interface {
	NewConfig()
}

type MongoConfig models.MongoSpec

func GetConfigDir() *MongoConfigDir {
	CONFIG_DIR, exists := os.LookupEnv(mongoConfigDir)
	mcd := &MongoConfigDir{}
	if !exists {
		mcd.root = DEFAULT_CONFIG_DIR
	} else {
		mcd.root = CONFIG_DIR
	}
	return mcd
}

func (mc *MongoConfig) LoadConfigFromFile() error {
	mcd := GetConfigDir()
	// fmt.Printf("Mongo Config Dir: %s\n", mcd.root)
	cfg, err := os.Open(filepath.Join(mcd.root, mongoConfigFile))
	if err != nil {
		log.Printf("Error opening Mongo config file: %s", err.Error())
		return err
	}
	defer cfg.Close()
	d := yaml.NewDecoder(cfg)
	if err := d.Decode(&mc); err != nil {
		return err
	}
	return nil
}

func (mc *MongoConfig) LoadConfigFromEnv() {
	envconfig.Process("", mc)
}

func (mc *MongoConfig) LoadConfigFromBinding(t string) error {
	b, err := bindings.NewBinding(t)

	if err != nil {
		log.Printf("Error while getting bindings: %s", err.Error())
		return err
	}

	if err := Decode(b, &mc); err != nil {
		return err
	}

	return nil
}

func (mc *MongoConfig) NewConfig() {
	mc.LoadConfigFromFile()
	mc.LoadConfigFromBinding("mongodb")
	mc.LoadConfigFromEnv()
}
