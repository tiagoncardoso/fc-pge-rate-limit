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
	TimeStamp int64  `json:"timestamp"`
	Requests  int    `json:"requests"`
	Details   string `json:"details"`
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
		rllog.Info("No cache for this client: " + err.Error())
	}

	setCacheDataPayload(&cacheData, details)

	return &Limiter{
		CacheClient: cacheClient,
		RateTimer:   rateTimer,
		CacheData:   cacheData,
		CacheKey:    key,
		Details:     details,
	}
}

func (l *Limiter) IsRateLimited() bool {
	if l.CacheData.Requests > l.RateTimer.MaxRequestsPerSecond {
		// TODO: Update cache TTL to window time and rate limited
		return true
	}

	cacheDataJson, err := json.Marshal(l.CacheData)
	if err != nil {
		rllog.Error("Failed to marshal CacheData: " + err.Error())
		return false
	}

	l.updateCache(string(cacheDataJson))

	return false
}

func (l *Limiter) updateCache(cacheDataJson string) {
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

func setCacheDataPayload(cacheData *CacheData, details string) *CacheData {
	if cacheData.Requests == 0 {
		cacheData.TimeStamp = time.Now().Unix()
	}

	cacheData.Details = details
	cacheData.Requests++

	return cacheData
}
