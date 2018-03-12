package test

import (
	"testing"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {

	startTime := time.Now().UTC()

	requestsPerSecond := 3
	window := 2
	iterations := 10

	rateLimiter := common.NewRateLimiter(requestsPerSecond, window)

	for i := 0; i < iterations; i++ {
		rateLimiter.RespectRateLimit()
		//DUMP(rateLimiter)
	}

	elapsed := time.Now().UTC().Sub(startTime)
	expectedRateLimit := time.Duration(iterations/requestsPerSecond) * time.Second

	//fmt.Printf("Elapsed time: %s\n", elapsed)
	//fmt.Printf("expectedRateLimit: %s\n", expectedRateLimit)

	assert.NotNil(t, rateLimiter)
	assert.Equal(t, true, elapsed > expectedRateLimit)
}
