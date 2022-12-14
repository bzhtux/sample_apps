package models

const (
	CONFIG_DIR = "/config"
)

type RedisSpec struct {
	Host       string `yaml:"host" envconfig:"REDIS_HOST"`
	Port       int    `yaml:"port" envconnfig:"REDIS_PORT"`
	Username   string `yaml:"username" envconfig:"REDIS_USERNAME"`
	Password   string `yaml:"password" envconfig:"REDIS_PASSWORD"`
	DB         int    `yaml:"database" envconfig:"REDIS_DB"`
	SSLenabled bool   `yaml:"sslenabled" envconfig:"REDIS_SSL"`
}

// type Config interface {
// 	NewConfig()
// }

type Messages struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// func (rc *RedisConfig) GetConfigFromFile() error {
// 	file, err := os.Open(CONFIG_DIR + "/redis.yaml")
// 	if err != nil {
// 		log.Println("Error loading config file ...", err.Error())
// 		return err
// 	}
// 	defer file.Close()

// 	d := yaml.NewDecoder(file)

// 	if err := d.Decode(&rc); err != nil {
// 		log.Println("Error loading config file ...", err.Error())
// 		return err
// 	}

// 	return nil
// }

// func (rc *RedisConfig) GetConfigFromEnv() {
// 	envconfig.Process("", rc)
// }

// func (rc *RedisConfig) GetConfigFromBindings(t string) error {
// 	b, err := bindings.NewBinding(t)
// 	if err != nil {
// 		log.Printf("Error while getting bindings: %s", err.Error())
// 		return err
// 	}
// 	rc.Host = b.Host
// 	fmt.Printf("Redis Host: %s\n", b.Host)
// 	rc.Port = int(b.Port)
// 	fmt.Printf("Redis Port: %d\n", b.Port)
// 	rc.Username = b.Username
// 	fmt.Printf("Redis Username: %s\n", b.Username)
// 	rc.Password = b.Password
// 	fmt.Printf("Redis Password: %s\n", b.Password)
// 	rc.SSLenabled = b.SSL
// 	fmt.Printf("Redis SSL: %v\n", b.SSL)
// 	return nil
// }

// func (rc *RedisConfig) NewConfig() *RedisConfig {
// 	rc.GetConfigFromFile()
// 	rc.GetConfigFromBindings("redis")
// 	rc.GetConfigFromEnv()
// 	return rc
// }
