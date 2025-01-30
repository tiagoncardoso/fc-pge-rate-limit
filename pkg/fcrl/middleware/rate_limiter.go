package middleware

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/cache"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/helpers"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/rllog"
	"net/http"
	"time"
)

type RateLimitOptions struct {
	RateLimitBy   time.Duration
	RateTimers    map[string]fcrl.RateTimer
	RequestIp     string
	RequestApiKey string
	RequestRoute  string
	CacheClient   cache.CacheInterface
	Opts          []Option
}

type Option func(rl *RateLimitOptions)

func RateLimiter(rateLimitBy string, opts ...Option) func(http.Handler) http.Handler {
	return NewRateLimitOptions(rateLimitBy, opts...).Handler
}

func NewRateLimitOptions(rateLimitBy string, opts ...Option) *RateLimitOptions {
	return &RateLimitOptions{
		RateTimers:  make(map[string]fcrl.RateTimer),
		RateLimitBy: helpers.ParseStringToTime(rateLimitBy),
		Opts:        opts,
	}
}

func (ro *RateLimitOptions) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rllog.Info("RateLimiter middleware")

		ro.RequestApiKey = r.Header.Get("API_KEY")
		ro.RequestIp = helpers.GetRequestIp(r)
		ro.RequestRoute = r.URL.String()

		for _, opt := range ro.Opts {
			opt(ro)
		}

		if ro.tooManyRequests() {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (ro *RateLimitOptions) tooManyRequests() bool {
	var limiter *fcrl.Limiter
	var cacheKey string

	cacheDetails := fmt.Sprintf("Route: %s :: IP: %s :: API Key: %s", ro.RequestRoute, ro.RequestIp, ro.RequestApiKey)

	if ro.noApiToken() {
		cacheKey = helpers.GenerateMD5Hash(ro.RequestIp + ro.RequestRoute)
		limiter = fcrl.NewLimiter(ro.CacheClient, ro.RateTimers["ip"], cacheKey, cacheDetails, ro.RateLimitBy)
	} else {
		cacheKey = helpers.GenerateMD5Hash(ro.RequestApiKey + ro.RequestRoute)
		limiter = fcrl.NewLimiter(ro.CacheClient, ro.RateTimers["token"], cacheKey, cacheDetails, ro.RateLimitBy)
	}

	return limiter.IsRateLimited()
}

func (ro *RateLimitOptions) noApiToken() bool {
	return ro.RequestApiKey == "" || ro.RateTimers["token"].MaxRequests == 0
}

func WithRateLimit(ipMaxRequests int, ipWindowTime int, tokenMaxRequests int, tokenWindowTime int) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["ip"] = fcrl.RateTimer{MaxRequests: ipMaxRequests, WindowTime: ipWindowTime}
		rl.RateTimers["token"] = fcrl.RateTimer{MaxRequests: tokenMaxRequests, WindowTime: tokenWindowTime}

		rllog.Info("WithRateLimit")
	}
}

func WithIPRateLimiter(ipMaxRequests int, ipWindowTime int) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["ip"] = fcrl.RateTimer{MaxRequests: ipMaxRequests, WindowTime: ipWindowTime}

		rllog.Info("WithIPRateLimiter")
	}
}

func WithApiKeyRateLimiter(tokenMaxRequests int, tokenWindowTime int) Option {
	return func(rl *RateLimitOptions) {
		rl.RateTimers["token"] = fcrl.RateTimer{MaxRequests: tokenMaxRequests, WindowTime: tokenWindowTime}

		rllog.Info("WithApiKeyRateLimiter")
	}
}

func WithRedisCache(redisHost string, redisPort int, redisPass string, ctx context.Context) Option {
	return func(rl *RateLimitOptions) {
		rllog.Info("WithRedisCache")

		rl.CacheClient = cache.NewRedisCache(redisHost, redisPort, redisPass, ctx)
	}
}

func WithRedisCacheClient(redisClient *redis.Client, ctx context.Context) Option {
	return func(rl *RateLimitOptions) {
		rl.CacheClient = cache.NewRedisClientCache(redisClient, ctx)
	}
}

func WithInMemoryCache(ctx context.Context) Option {
	return func(rl *RateLimitOptions) {
		rllog.Info("WithInMemoryCache")

		rl.CacheClient = cache.NewMemoryCache(ctx)
	}
}
