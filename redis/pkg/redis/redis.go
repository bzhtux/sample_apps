package redis

import (
	"log"
	"strconv"

	"github.com/bzhtux/sample_apps/redis/models"
	"github.com/gomodule/redigo/redis"
)

func NewPool() *redis.Pool {

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
