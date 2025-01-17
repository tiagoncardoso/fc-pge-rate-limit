package main

import (
	"context"
	"github.com/tiagoncardoso/fc-pge-rate-limit/config"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/cache"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/web"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/web/handler"
)

func main() {
	ctx := context.Background()
	envConf, err := config.SetupEnvConfig()
	if err != nil {
		panic(err)
	}

	redisClient := cache.NewRedisConfig(envConf.RedisHost, envConf.RedisPort, envConf.RedisPass, ctx)
	err = redisClient.Set("key", "value", envConf.BlockTime)
	if err != nil {
		panic(err)
	}

	timeHandler := handler.NewTimeApiHandler(ctx)

	webServer := web.NewWebServer(envConf.AppPort)
	webServer.Router.Use()

	webServer.AddHandler("/time/greetings", "GET", timeHandler.GetActualDate)

	webServer.Start()
}
