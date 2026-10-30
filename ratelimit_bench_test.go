package main

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"

	"ratelimiter/core"
	"ratelimiter/store"
)

func BenchmarkTokenBucket(b *testing.B) {
	limiter := core.NewTokenBucketLimiter(1000000, 1000000)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow("client-1")
		}
	})
}

func BenchmarkSlidingWindow(b *testing.B) {
	limiter := core.NewSlidingWindowLimiter(1000000, 60)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			limiter.Allow("client-1")
		}
	})
}

func BenchmarkRedisStore(b *testing.B) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	redisStore := store.NewRedisStore(client)
	ctx := context.Background()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			redisStore.Allow(ctx, "client-1", 1000000, 60)
		}
	})
}
