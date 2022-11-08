package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/bzhtux/sample_apps/postgresql/models"
	"github.com/bzhtux/servicebinding/bindings"
	"github.com/kelseyhightower/envconfig"
	. "github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

const (
	DEFAULT_CONFIG_DIR = "/config"
	pgConfigDir        = "MONGO_CONFIG_DIR"
	pgConfigFile       = "mongo.yaml"
)

type PGConfigDir struct {
	root string
}

type Config interface {
	NewConfig()
}

type PGConfig models.PGSpec

func GetConfigDir() *PGConfigDir {
	CONFIG_DIR, exists := os.LookupEnv(pgConfigDir)
	pgcd := &PGConfigDir{}
	if !exists {
		pgcd.root = DEFAULT_CONFIG_DIR
	} else {
		pgcd.root = CONFIG_DIR
	}
	return pgcd
}

func (pgc *PGConfig) LoadConfigFromFile() error {
	pgcd := GetConfigDir()
	cfg, err := os.Open(filepath.Join(pgcd.root, pgConfigFile))
	if err != nil {
		log.Printf("Error opening PostgreSQL config file: %s", err.Error())
		return err
	}
	defer cfg.Close()
	d := yaml.NewDecoder(cfg)
	if err := d.Decode(&pgc); err != nil {
		return err
	}
	return nil
}

func (cfg *PGConfig) LoadConfigFromEnv() {
	envconfig.Process("", cfg)
}

func (cfg *PGConfig) LoadConfigFromBindings(t string) error {

	b, err := bindings.NewBinding(t)
	if err != nil {
		log.Printf("Error while getting bindings: %s", err.Error())
		return err
	}

	if err := Decode(b, &cfg); err != nil {
		return err
	}
	// cfg.Host = b.Host
	// log.Printf("PostgreSQL Host: %s\n", cfg.Host)
	// cfg.Port = int(b.Port)
	// log.Printf("PostgreSQL Port: %d\n", b.Port)
	// cfg.Username = b.Username
	// log.Printf("PostgreSQL Username: %s\n", b.Username)
	// cfg.Password = b.Password
	// log.Printf("PostgreSQL Password: %s\n", b.Password)
	// cfg.DB = b.Database
	log.Printf("> DB: %s\n", b.Database)
	log.Printf("> PG_DB: %S\n", cfg.Database)
	// cfg.SSLenabled = b.SSL
	// log.Printf("PostgreSQL SSL: %v\n",)
	return nil
}

func (cfg *PGConfig) NewConfig() *PGConfig {
	cfg.LoadConfigFromFile()
	cfg.LoadConfigFromBindings("postgresql")
	cfg.LoadConfigFromEnv()
	log.Printf("> Host: %s\n", cfg.Host)
	log.Printf("> Port: %s\n", cfg.Port)
	log.Printf("> User: %s\n", cfg.Username)
	log.Printf("> Pass: %s\n", cfg.Password)
	log.Printf("> DB: %s\n", cfg.Database)
	return cfg
}
