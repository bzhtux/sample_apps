package config

import (
	"log"
	"os"
	"reflect"

	"github.com/bzhtux/sample_apps/redis/models"
	"github.com/bzhtux/servicebinding/bindings"
	"github.com/kelseyhightower/envconfig"
	. "github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type Entry struct {
	Value            interface{}
	AuthorizedValues []interface{} // Leave empty for "any"
	Type             reflect.Kind
	IsSlice          bool
}

type object map[string]interface{}

// type RedisConfig struct {
// 	*models.RedisSpec
// }

type RedisConfig models.RedisSpec

var config object

var configDefaults = object{
	"host":         &Entry{"localhost", []interface{}{}, reflect.String, false},
	"port":         &Entry{5432, []interface{}{}, reflect.Int16, false},
	"uri":          &Entry{"", []interface{}{}, reflect.String, false},
	"username":     &Entry{"", []interface{}{}, reflect.String, false},
	"password":     &Entry{"", []interface{}{}, reflect.String, false},
	"database":     &Entry{"", []interface{}{}, reflect.String, false},
	"ssl":          &Entry{true, []interface{}{}, reflect.Bool, false},
	"certificates": &Entry{"", []interface{}{}, reflect.Slice, true},
}

const (
	CONFIG_DIR = "/config"
)

func loadDefaultConfig(src, dest object) {
	for k, v := range src {
		if obj, ok := v.(object); ok {
			sub := make(object, len(obj))
			loadDefaultConfig(obj, sub)
			dest[k] = sub
		} else {
			entry := v.(*Entry)
			val := entry.Value
			t := reflect.TypeOf(val)
			if t != nil && t.Kind() == reflect.Slice {
				list := reflect.ValueOf(val)
				length := list.Len()
				slice := reflect.MakeSlice(reflect.SliceOf(t.Elem()), 0, length)
				for i := 0; i < length; i++ {
					slice = reflect.Append(slice, list.Index(i))
				}
				val = slice.Interface()
			}
			dest[k] = &Entry{val, entry.AuthorizedValues, entry.Type, entry.IsSlice}
		}
	}
}

func (rc *RedisConfig) GetConfigFromFile() error {
	file, err := os.Open(CONFIG_DIR + "/redis.yaml")
	if err != nil {
		log.Println("Error loading config file ...", err.Error())
		return err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&rc); err != nil {
		log.Println("Error loading config file ...", err.Error())
		return err
	}

	return nil
}

func (rc *RedisConfig) GetConfigFromEnv() {
	envconfig.Process("", rc)
}

func (rc *RedisConfig) GetConfigFromBindings(t string) error {
	b, err := bindings.NewBinding(t)
	if err != nil {
		log.Printf("Error while getting bindings: %s", err.Error())
		return err
	}

	if err := Decode(b, &rc); err != nil {
		return err
	}
	// rc.Host = b.Host
	// fmt.Printf("Redis Host: %s\n", b.Host)
	// rc.Port = int(b.Port)
	// fmt.Printf("Redis Port: %d\n", b.Port)
	// rc.Username = b.Username
	// fmt.Printf("Redis Username: %s\n", b.Username)
	// rc.Password = b.Password
	// fmt.Printf("Redis Password: %s\n", b.Password)
	// rc.SSLenabled = b.SSL
	// fmt.Printf("Redis SSL: %v\n", b.SSL)
	return nil
}

func (rc *RedisConfig) NewConfig() *RedisConfig {
	rc.GetConfigFromFile()
	rc.GetConfigFromBindings("redis")
	rc.GetConfigFromEnv()
	return rc
}
