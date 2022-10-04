package redis

import (
	"log"
	"strconv"

	"github.com/bzhtux/sample_apps/redis/models"
	"github.com/gomodule/redigo/redis"
)

var (
	pool = newPool()
)

func newPool() *redis.Pool {

	rc := new(models.RedisConfig)
	rc.NewConfig()

	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				rc.Host+":"+strconv.Itoa(rc.Port),
				redis.DialDatabase(rc.DB),
				redis.DialUsername(rc.Username),
				redis.DialPassword(rc.Password),
				redis.DialUseTLS(rc.SSLenabled),
			)
			if err != nil {
				log.Println(err.Error())
			}
			return c, err
		},
	}
}

func RedisConnect() {
	client := pool.Get()
	defer client.Close()
	log.Println("Setting a new key Hello with value From Tanzu Application Platform !")
	_, err := client.Do("SET", "Hello", "From Tanzu Application Platform !")
	if err != nil {
		log.Println(err.Error())
	}
	res, err := client.Do("GET", "Hello")
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("Getting key Hello => %s \n", res)
}
