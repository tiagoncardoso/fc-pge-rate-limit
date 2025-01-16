package main

import (
	"context"
	"fmt"
	"github.com/tiagoncardoso/fc-pge-rate-limit/config"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/cache"
)

func main() {
	ctx := context.Background()
	envConf, err := config.SetupEnvConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(envConf)

	redisClient := cache.NewRedisConfig(envConf.RedisHost, envConf.RedisPort, envConf.RedisPass, ctx)
	err = redisClient.Set("key", "{value}", envConf.BlockTime)
	if err != nil {
		panic(err)
	}

	fmt.Println("Hello, World!")
}
