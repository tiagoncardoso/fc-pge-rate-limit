package e2e

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"github.com/tiagoncardoso/fc-pge-rate-limit/internal/infra/web/handler"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/middleware"
	"net/http/httptest"
)

type RequestGreetingsTestSuite struct {
	suite.Suite
	server             *httptest.Server
	redisMock          *miniredis.Miniredis
	RateLimiterControl RateLimiterControl
}

type RateLimiterControl struct {
	RateLimitBy       string
	IpMaxRequests     int
	IpWindowTime      int
	ApiKeyMaxRequests int
	ApiKeyWindowTime  int
}

func (suite *RequestGreetingsTestSuite) SetupSuite() {
	ctx := context.Background()
	suite.RateLimiterControl = RateLimiterControl{
		RateLimitBy:       "Second",
		IpMaxRequests:     3,
		IpWindowTime:      5,
		ApiKeyMaxRequests: 4,
		ApiKeyWindowTime:  5,
	}

	mr, err := miniredis.Run()
	suite.NoError(err)
	suite.redisMock = mr

	redisClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	timeHandler := handler.NewGreetingsHandler(ctx)
	r := chi.NewRouter()
	r.Use(middleware.RateLimiter(
		suite.RateLimiterControl.RateLimitBy,
		middleware.WithIPRateLimiter(suite.RateLimiterControl.IpMaxRequests, suite.RateLimiterControl.IpWindowTime),
		middleware.WithApiKeyRateLimiter(suite.RateLimiterControl.ApiKeyMaxRequests, suite.RateLimiterControl.ApiKeyWindowTime),
		middleware.WithRedisCacheClient(redisClient, ctx),
	))
	r.Get("/time/greetings", timeHandler.GetDateAndGreetings)

	suite.server = httptest.NewServer(r)
}

func (suite *RequestGreetingsTestSuite) TearDownSuite() {
	suite.server.Close()
	suite.redisMock.Close()
}
