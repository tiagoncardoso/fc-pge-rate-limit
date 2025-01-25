package fcrl

import (
	"encoding/json"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/cache"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/helpers"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/rllog"
	"time"
)

const (
	SECOND       = 1
	MINUTE       = 60 * SECOND
	FIVE_MINUTES = 5 * MINUTE
)

const (
	RATE_LIMIT_NOT_EXCEEDED = iota
	RATE_LIMIT_EXCEEDED
)

type RateTimer struct {
	MaxRequestsPerSecond int
	WindowTime           time.Duration
}

type CacheData struct {
	TimeStamp       int64  `json:"timestamp"`
	Requests        int    `json:"requests"`
	Details         string `json:"details"`
	RateLimitStatus int    `json:"rate_limit_status"`
}

type Limiter struct {
	CacheClient cache.CacheInterface
	RateTimer   RateTimer
	CacheData   CacheData
	CacheKey    string
	Details     string
}

func NewLimiter(cacheClient cache.CacheInterface, rateTimer RateTimer, key string, details string) *Limiter {
	var cacheData CacheData
	requestsStr, _ := cacheClient.Get(key)

	err := json.Unmarshal([]byte(requestsStr), &cacheData)
	if err != nil {
		rllog.Info("No cache for this client request yet. It will be create now")
	}

	return &Limiter{
		CacheClient: cacheClient,
		RateTimer:   rateTimer,
		CacheData:   cacheData,
		CacheKey:    key,
		Details:     details,
	}
}

func (l *Limiter) IsRateLimited() bool {
	var cacheDataJson string

	l.CacheData.Requests++
	l.CacheData.Details = l.Details
	l.CacheData.RateLimitStatus = RATE_LIMIT_NOT_EXCEEDED

	if l.rateLimitExceeded() {
		// TODO: Update cache TTL to window time
		l.CacheData.RateLimitStatus = RATE_LIMIT_EXCEEDED
		cacheDataJson = helpers.ParseStructToString(l.CacheData)

		l.saveCache(cacheDataJson)
		return true
	}

	if l.isFirstRequest() {
		l.CacheData.TimeStamp = time.Now().Unix()
	}

	cacheDataJson = helpers.ParseStructToString(l.CacheData)
	l.saveCache(cacheDataJson)

	return false
}

func (l *Limiter) saveCache(cacheDataJson string) {
	var err error

	if l.isFirstRequest() {
		err = l.CacheClient.Set(l.CacheKey, cacheDataJson, FIVE_MINUTES)
	} else {
		err = l.CacheClient.Update(l.CacheKey, cacheDataJson)
	}

	if err != nil {
		rllog.Error("Failed to set cache: " + err.Error())
	}
}

func (l *Limiter) isFirstRequest() bool {
	return l.CacheData.Requests <= 1
}

func (l *Limiter) rateLimitExceeded() bool {
	return l.CacheData.Requests > l.RateTimer.MaxRequestsPerSecond || l.CacheData.RateLimitStatus == RATE_LIMIT_EXCEEDED
}
