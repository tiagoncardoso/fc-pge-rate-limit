package middleware

import (
	"context"
	"fmt"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/cache"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/helpers"
	"net/http"
	"time"
)

type RateTimer struct {
	maxRequestsPerSecond int
	windowTime           time.Duration
}

type RateLimitOptions struct {
	RateTimers    map[string]RateTimer
	RequestIp     string
	RequestApiKey string
	CacheClient   cache.CacheInterface
	Opts          []Option
}

// TODO: Caso não haja parâmetros de redis, utilizar um cache em memória

type Option func(rl *RateLimitOptions)

func RateLimiter(opts ...Option) func(http.Handler) http.Handler {
	return NewRateLimitOptions(opts...).Handler
}

func NewRateLimitOptions(opts ...Option) *RateLimitOptions {
	return &RateLimitOptions{
		RateTimers: make(map[string]RateTimer),
		Opts:       opts,
	}
}

func (ro *RateLimitOptions) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RateLimiter middleware")

		ro.RequestApiKey = r.Header.Get("API_KEY")
		ro.RequestIp = helpers.GetRequestIp(r)

		for _, opt := range ro.Opts {
			opt(ro)
		}

		fmt.Println("A lógica do limiter vem aqui")

		handler.ServeHTTP(w, r)
	})
}

func WithRateLimit(ipMaxRequestsPerSecond int, ipWindowTime time.Duration, tokenMaxRequestsPerSecond int, tokenWindowTime time.Duration) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["ip"] = RateTimer{maxRequestsPerSecond: ipMaxRequestsPerSecond, windowTime: ipWindowTime}
		rl.RateTimers["token"] = RateTimer{maxRequestsPerSecond: tokenMaxRequestsPerSecond, windowTime: tokenWindowTime}

		fmt.Println("WithRateLimit")
		fmt.Println("The IP is: ", rl.RequestIp)
		fmt.Println("The API Key is: ", rl.RequestApiKey)
	}
}

func WithIPRateLimiter(ipMaxRequestsPerSecond int, ipWindowTime time.Duration) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["ip"] = RateTimer{maxRequestsPerSecond: ipMaxRequestsPerSecond, windowTime: ipWindowTime}

		fmt.Println("WithIPRateLimiter")
		fmt.Println("The IP is: ", rl.RequestIp)
	}
}

func WithApiKeyRateLimiter(tokenMaxRequestsPerSecond int, tokenWindowTime time.Duration) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["token"] = RateTimer{maxRequestsPerSecond: tokenMaxRequestsPerSecond, windowTime: tokenWindowTime}

		fmt.Println("WithApiKeyRateLimiter")
		fmt.Println("The API Key is: ", rl.RequestApiKey)
	}
}

func WithRedisCache(redisHost string, redisPort int, redisPass string, ctx context.Context) Option {
	return func(rl *RateLimitOptions) {
		fmt.Println("WithRedisCache")

		rl.CacheClient = cache.NewRedisConfig(redisHost, redisPort, redisPass, ctx)

		err := rl.CacheClient.Set("key", "value", 60)
		if err != nil {
			panic(err)
		}
	}
}
