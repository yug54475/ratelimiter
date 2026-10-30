package core

import (
	"sync"
	"time"
)

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

type TokenBucketLimiter struct {
	capacity   float64
	refillRate float64 // tokens added per second
	mu         sync.Mutex
	buckets    map[string]*bucket
}

func NewTokenBucketLimiter(capacity float64, refillPerSecond float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:   capacity,
		refillRate: refillPerSecond,
		buckets:    make(map[string]*bucket),
	}
}

func (t *TokenBucketLimiter) Allow(clientID string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	b, ok := t.buckets[clientID]
	now := time.Now()
	if !ok {
		b = &bucket{tokens: t.capacity, lastRefill: now}
		t.buckets[clientID] = b
	}

	elapsed := now.Sub(b.lastRefill).Seconds()
	b.tokens += elapsed * t.refillRate
	if b.tokens > t.capacity {
		b.tokens = t.capacity
	}
	b.lastRefill = now

	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}
