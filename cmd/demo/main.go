package main

import (
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"

	"ratelimiter/core"
	"ratelimiter/middleware"
	"ratelimiter/store"
)

func main() {
	// --- In-memory token bucket, on "/" ---
	memLimiter := core.NewTokenBucketLimiter(5, 1)

	memHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "request allowed! (in-memory)")
	})

	memRateLimited := middleware.RateLimit(memLimiter, middleware.ByIP)(memHandler)
	http.Handle("/", memRateLimited)

	// --- Redis-backed limiter, on "/redis" ---
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	redisStore := store.NewRedisStore(redisClient)

	http.HandleFunc("/redis", func(w http.ResponseWriter, r *http.Request) {
		clientID := r.RemoteAddr
		allowed, err := redisStore.Allow(r.Context(), clientID, 5, 10) // 5 requests per 10 sec
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		if !allowed {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		fmt.Fprintln(w, "request allowed! (redis)")
	})

	fmt.Println("listening on :8080")
	http.ListenAndServe(":8080", nil)
}
