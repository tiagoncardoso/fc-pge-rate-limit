package fcrl

import (
	"encoding/json"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/cache"
	"github.com/tiagoncardoso/fc-pge-rate-limit/pkg/fcrl/rllog"
	"time"
)

const (
	SECOND       = 1
	MINUTE       = 60 * SECOND
	FIVE_MINUTES = 5 * MINUTE
)

type RateTimer struct {
	MaxRequestsPerSecond int
	WindowTime           time.Duration
}

type CacheData struct {
	TimeStamp int64 `json:"timestamp"`
	Requests  int   `json:"requests"`
	TTL       int   `json:"ttl"`
}

type Limiter struct {
	CacheClient cache.CacheInterface
	RateTimer   RateTimer
	CacheData   CacheData
}

func NewLimiter(cacheClient cache.CacheInterface, rateTimer RateTimer, key string) *Limiter {
	var cacheData CacheData
	requestsStr, _ := cacheClient.Get(key)

	err := json.Unmarshal([]byte(requestsStr), &cacheData)
	if err != nil {
		rllog.Info("No cache for this client: " + err.Error())
	}

	if cacheData.Requests >= 1 {
		cacheData.Requests++
	} else {
		cacheData.Requests = 1
		cacheData.TimeStamp = time.Now().Unix()
	}

	return &Limiter{
		CacheClient: cacheClient,
		RateTimer:   rateTimer,
		CacheData:   cacheData,
	}
}

func (l *Limiter) IsRateLimited(key string) bool {
	if l.CacheData.Requests >= l.RateTimer.MaxRequestsPerSecond {
		return true
	}

	cacheDataJson, err := json.Marshal(l.CacheData)
	if err != nil {
		rllog.Error("Failed to marshal CacheData: " + err.Error())
		return false
	}

	if l.CacheData.Requests == 1 {
		if err := l.CacheClient.Set(key, string(cacheDataJson), MINUTE); err != nil {
			rllog.Error(err.Error())
		}
	} else {
		if err := l.CacheClient.Update(key, string(cacheDataJson)); err != nil {
			rllog.Error(err.Error())
		}
	}

	return false
}
