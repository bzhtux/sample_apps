package main

import (
	"fmt"

	"github.com/bzhtux/sample_apps/redis/pkg/redis"
)

func main() {
	fmt.Println("\033[32mTesting Golang Redis ...\033[0m")
	redis.RedisConnect()
}
