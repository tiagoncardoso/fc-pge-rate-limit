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
	RateLimitBy time.Duration
	MaxRequests int
	WindowTime  int
}

type CacheData struct {
	InitialTimeStamp int64  `json:"initial_timestamp"`
	Requests         int    `json:"requests"`
	Details          string `json:"details"`
	RateLimitStatus  int    `json:"rate_limit_status"`
}

type Limiter struct {
	CacheClient       cache.CacheInterface
	RateTimer         RateTimer
	CacheData         CacheData
	CacheKey          string
	Details           string
	RateLimitAchieved bool
}

func NewLimiter(cacheClient cache.CacheInterface, rateTimer RateTimer, key string, details string, rateLimitBy time.Duration) *Limiter {
	var cacheData CacheData
	requestsStr, _ := cacheClient.Get(key)

	err := json.Unmarshal([]byte(requestsStr), &cacheData)
	if err != nil {
		rllog.Info("No cache for this client request yet. It will be create now")
	}

	rateTimer.RateLimitBy = rateLimitBy

	return &Limiter{
		CacheClient:       cacheClient,
		RateTimer:         rateTimer,
		CacheData:         cacheData,
		CacheKey:          key,
		Details:           details,
		RateLimitAchieved: cacheData.RateLimitStatus == RATE_LIMIT_EXCEEDED,
	}
}

func (l *Limiter) IsRateLimited() bool {
	var cacheDataJson string

	if l.RateLimitAchieved {
		return true
	}

	l.CacheData.Requests++
	l.CacheData.Details = l.Details

	if l.rateLimitExceeded() {
		l.CacheData.RateLimitStatus = RATE_LIMIT_EXCEEDED
		cacheDataJson = helpers.ParseStructToString(l.CacheData)

		l.saveCache(cacheDataJson, true)
		return true
	}

	if l.isFirstRequest() {
		l.CacheData.InitialTimeStamp = time.Now().Unix()
	}

	l.CacheData.RateLimitStatus = RATE_LIMIT_NOT_EXCEEDED
	cacheDataJson = helpers.ParseStructToString(l.CacheData)
	l.saveCache(cacheDataJson, false)

	return false
}

func (l *Limiter) saveCache(cacheDataJson string, rateLimitAchieved bool) {
	var err error
	var ttl time.Duration = -1

	if rateLimitAchieved {
		ttl = time.Duration(l.RateTimer.WindowTime) * l.RateTimer.RateLimitBy
	}

	if l.isFirstRequest() {
		err = l.CacheClient.Set(l.CacheKey, cacheDataJson, l.RateTimer.RateLimitBy)
	} else {
		err = l.CacheClient.Update(l.CacheKey, cacheDataJson, ttl)
	}

	if err != nil {
		rllog.Error("Failed to set cache: " + err.Error())
	}
}

func (l *Limiter) isFirstRequest() bool {
	return l.CacheData.Requests <= 1
}

func (l *Limiter) rateLimitExceeded() bool {
	return l.CacheData.Requests > l.RateTimer.MaxRequests || l.CacheData.RateLimitStatus == RATE_LIMIT_EXCEEDED
}
