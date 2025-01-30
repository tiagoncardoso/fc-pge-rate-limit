package main

import (
	"context"
	"github.com/tiagoncardoso/fc-pge-rate-limit/config"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/web"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/web/handler"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/middleware"
)

func main() {
	ctx := context.Background()
	envConf, err := config.SetupEnvConfig()
	if err != nil {
		panic(err)
	}

	timeHandler := handler.NewGreetingsHandler(ctx)

	webServer := web.NewWebServer(envConf.AppPort)
	webServer.Router.Use(middleware.RateLimiter(
		envConf.RateLimitBy,
		middleware.WithIPRateLimiter(envConf.IpLimitRate, envConf.IpWindowTime),
		middleware.WithApiKeyRateLimiter(envConf.TokenLimitRate, envConf.ApiTokenWindowTime),
		middleware.WithRedisCache(envConf.RedisHost, envConf.RedisPort, envConf.RedisPass, ctx),
		//middleware.WithInMemoryCache(ctx),
	))

	webServer.AddHandler("/time/greetings", "GET", timeHandler.GetDateAndGreetings)
	webServer.AddHandler("/time/japanese-greetings", "GET", timeHandler.GetJapaneseDateAndGreetings)

	webServer.Start()
}
