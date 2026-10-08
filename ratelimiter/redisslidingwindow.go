// ratelimiter/redis_sliding_window.go
package ratelimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSlidingWindowLimiter struct {
	client     *redis.Client
	windowSize time.Duration
	threshold  int
}

func NewRedisSlidingWindowLimiter(client *redis.Client, windowSize time.Duration, threshold int) *RedisSlidingWindowLimiter {
	return &RedisSlidingWindowLimiter{
		client: client,
		windowSize: windowSize,
		threshold: threshold}
}

func (l *RedisSlidingWindowLimiter) Allow(key string) bool {
	ctx := context.Background()
	windowSecs := int64(l.windowSize.Seconds())
	now := time.Now().Unix()
	nowWindow := now / windowSecs

	currentKey := fmt.Sprintf("rl:%s:%d", key, nowWindow)
	previousKey := fmt.Sprintf("rl:%s:%d", key, nowWindow-1)

	currentCount, _ := l.client.Get(ctx, currentKey).Int64()   // Int64() returns 0 if the key doesn't exist
	previousCount, _ := l.client.Get(ctx, previousKey).Int64()

	elapsed := float64(now%windowSecs) / float64(windowSecs)
	estimated := float64(previousCount)*(1-elapsed) + float64(currentCount)

	if estimated >= float64(l.threshold) {
		return false
	}

	pipe := l.client.Pipeline()
	pipe.Incr(ctx, currentKey)
	pipe.Expire(ctx, currentKey, 2*l.windowSize)
	pipe.Exec(ctx)
	return true
}