package main

import (
	"sync"
	"testing"

	"ratelimiter/core"
)

func TestTokenBucketConcurrency(t *testing.T) {
	limit := 100.0
	limiter := core.NewTokenBucketLimiter(limit, 0) // no refill, so total allowed = exactly limit

	var wg sync.WaitGroup
	var mu sync.Mutex
	allowedCount := 0

	totalRequests := 100000
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()
			if limiter.Allow("client-1") {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowedCount != int(limit) {
		t.Errorf("expected exactly %d allowed requests, got %d", int(limit), allowedCount)
	} else {
		t.Logf("correct: exactly %d out of %d concurrent requests were allowed", allowedCount, totalRequests)
	}
}

func TestSlidingWindowConcurrency(t *testing.T) {
	limit := 100
	limiter := core.NewSlidingWindowLimiter(limit, 60_000_000_000) // 60s window, in nanoseconds via time.Duration

	var wg sync.WaitGroup
	var mu sync.Mutex
	allowedCount := 0

	totalRequests := 100000
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()
			if limiter.Allow("client-1") {
				mu.Lock()
				allowedCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if allowedCount != limit {
		t.Errorf("expected exactly %d allowed requests, got %d", limit, allowedCount)
	} else {
		t.Logf("correct: exactly %d out of %d concurrent requests were allowed", allowedCount, totalRequests)
	}
}
