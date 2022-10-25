package models

type MongoSpec struct {
	Host     string `yaml:"host" envconfig:"MONGO_HOST"`
	Port     uint16 `yaml:"port" envconfig:"MONGO_PORT"`
	Username string `yaml:"username" envconfig:"MONGO_USERNAME"`
	Password string `yaml:"password" envconfig:"MONGO_PASSWORD"`
}

type MongoCollection struct {
	Database   string
	Collection string
}

type Books struct {
	ID     int
	Title  string `binding:"required"`
	Author string `binding:"required"`
}
