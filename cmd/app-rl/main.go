package main

import (
	"context"
	"github.com/tiagoncardoso/fc-pge-rate-limit/config"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/web"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/web/handler"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/middleware"
	"time"
)

func main() {
	ctx := context.Background()
	envConf, err := config.SetupEnvConfig()
	if err != nil {
		panic(err)
	}

	timeHandler := handler.NewTimeApiHandler(ctx)

	webServer := web.NewWebServer(envConf.AppPort)
	webServer.Router.Use(middleware.RateLimiter(
		middleware.WithIPRateLimiter(envConf.IpLimitRatePerSecond, time.Duration(envConf.IpBlockTime)*time.Second),
		middleware.WithApiKeyRateLimiter(envConf.TokenLimitRatePerSecond, time.Duration(envConf.TokenBlockTime)*time.Second),
		middleware.WithRedisCache(envConf.RedisHost, envConf.RedisPort, envConf.RedisPass, ctx),
	))

	webServer.AddHandler("/time/greetings", "GET", timeHandler.GetActualDate)

	webServer.Start()
}
