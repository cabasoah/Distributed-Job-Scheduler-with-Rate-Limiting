package ratelimiter

import (
	"sync"
	"time"
)

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
}

type TokenBucketLimiter struct {
	mu         sync.Mutex
	buckets    map[string]*tokenBucket
	capacity   float64
	refillRate float64 // token added per second
}

func NewTokenBucketLimiter(capacity, refillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		buckets:    make(map[string]*tokenBucket),
		capacity:   capacity,
		refillRate: refillRate,
	}
}

func (l *TokenBucketLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	b, exists := l.buckets[key]
	if !exists {
		b = &tokenBucket{tokens: l.capacity, lastRefill: time.Now()}
		l.buckets[key] = b
	}

	now := time.Now()
	b.tokens = min(l.capacity, b.tokens+now.Sub(b.lastRefill).Seconds()*l.refillRate)
	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false

}
