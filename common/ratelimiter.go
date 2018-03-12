package common

import (
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	maxRequests     int
	interval        int
	currentRequests int
	lastRequest     time.Time
	lock            sync.Mutex
}

func NewRateLimiter(maxRequests int, perSecond int) *RateLimiter {
	return &RateLimiter{
		maxRequests: maxRequests,
		interval:    perSecond,
		lastRequest: time.Now()}
}

func (rateLimiter *RateLimiter) RespectRateLimit() {
	var resetCounter bool
	duration := time.Now().UTC().Sub(rateLimiter.lastRequest)
	coolOffPeriod := time.Duration(rateLimiter.interval) * time.Second
	if rateLimiter.currentRequests >= rateLimiter.maxRequests && duration <= coolOffPeriod {
		//rateLimiter.ctx.GetLogger().Debugf("[RateLimiter] Waiting %s for rate limit cool-off: \n", coolOffPeriod)
		fmt.Printf("[RateLimiter] Waiting %s for rate limit cool-off: \n", coolOffPeriod)
		time.Sleep(coolOffPeriod)
		resetCounter = true
	}
	rateLimiter.lock.Lock()
	if resetCounter {
		rateLimiter.currentRequests = 1
	} else {
		rateLimiter.currentRequests++
	}
	rateLimiter.lastRequest = time.Now().UTC()
	rateLimiter.lock.Unlock()
}

func (rateLimiter *RateLimiter) String() string {
	return fmt.Sprintf("maxRequests: %d, interval: %ds, lastRequest: %s, currentRequests: %d",
		rateLimiter.maxRequests, rateLimiter.interval, rateLimiter.lastRequest, rateLimiter.currentRequests)
}
