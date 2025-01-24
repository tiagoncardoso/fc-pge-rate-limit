package middleware

import (
	"context"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/cache"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/helpers"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/rllog"
	"net/http"
	"time"
)

type RateLimitOptions struct {
	RateTimers    map[string]fcrl.RateTimer
	RequestIp     string
	RequestApiKey string
	CacheClient   cache.CacheInterface
	Opts          []Option
}

// TODO: RateLimiter baseado em API Key deve sobrepor o RateLimiter baseado em IP
// TODO: Caso não haja parâmetros de redis, utilizar um cache em memória

type Option func(rl *RateLimitOptions)

func RateLimiter(opts ...Option) func(http.Handler) http.Handler {
	return NewRateLimitOptions(opts...).Handler
}

func NewRateLimitOptions(opts ...Option) *RateLimitOptions {
	return &RateLimitOptions{
		RateTimers: make(map[string]fcrl.RateTimer),
		Opts:       opts,
	}
}

func (ro *RateLimitOptions) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rllog.Info("RateLimiter middleware")

		ro.RequestApiKey = r.Header.Get("API_KEY")
		ro.RequestIp = helpers.GetRequestIp(r)

		for _, opt := range ro.Opts {
			opt(ro)
		}

		if ro.RequestApiKey == "" || ro.RateTimers["token"].MaxRequestsPerSecond == 0 {
			limiter := fcrl.NewLimiter(ro.CacheClient, ro.RateTimers["ip"], "ip")
			limiter.IsRateLimited(ro.RequestIp)
		} else {
			limiter := fcrl.NewLimiter(ro.CacheClient, ro.RateTimers["token"], ro.RequestApiKey)
			limiter.IsRateLimited(ro.RequestApiKey)
		}

		handler.ServeHTTP(w, r)
	})
}

func WithRateLimit(ipMaxRequestsPerSecond int, ipWindowTime time.Duration, tokenMaxRequestsPerSecond int, tokenWindowTime time.Duration) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["ip"] = fcrl.RateTimer{MaxRequestsPerSecond: ipMaxRequestsPerSecond, WindowTime: ipWindowTime}
		rl.RateTimers["token"] = fcrl.RateTimer{MaxRequestsPerSecond: tokenMaxRequestsPerSecond, WindowTime: tokenWindowTime}

		rllog.Info("WithRateLimit")
		rllog.Info("TIP", rl.RequestIp)
		rllog.Info("API Key", rl.RequestApiKey)
	}
}

func WithIPRateLimiter(ipMaxRequestsPerSecond int, ipWindowTime time.Duration) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["ip"] = fcrl.RateTimer{MaxRequestsPerSecond: ipMaxRequestsPerSecond, WindowTime: ipWindowTime}

		rllog.Info("WithIPRateLimiter")
		rllog.Info("IP", rl.RequestIp)
	}
}

func WithApiKeyRateLimiter(tokenMaxRequestsPerSecond int, tokenWindowTime time.Duration) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["token"] = fcrl.RateTimer{MaxRequestsPerSecond: tokenMaxRequestsPerSecond, WindowTime: tokenWindowTime}

		rllog.Info("WithApiKeyRateLimiter")
		rllog.Info("API Key", rl.RequestApiKey)
	}
}

func WithRedisCache(redisHost string, redisPort int, redisPass string, ctx context.Context) Option {
	return func(rl *RateLimitOptions) {
		rllog.Info("WithRedisCache")

		rl.CacheClient = cache.NewRedisConfig(redisHost, redisPort, redisPass, ctx)
	}
}
