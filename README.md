# Go Rate Limiter

A thread-safe Go rate limiter with in-memory and Redis-backed implementations.

## Features

- Token bucket rate limiter with configurable capacity and refill rate.
- Sliding window rate limiter with per-client tracking.
- HTTP middleware for applying rate limits to handlers.
- Redis-backed fixed-window limiting using an atomic Lua script.
- Concurrent access tests and benchmarks.

## Requirements

- Go 1.25 or newer
- Redis 6 or newer for the Redis-backed limiter

## Installation

```bash
go get github.com/YOUR_USERNAME/ratelimiter
```

Replace `YOUR_USERNAME` with your GitHub username.

## In-memory usage

```go
limiter := core.NewTokenBucketLimiter(5, 1)

if limiter.Allow("client-1") {
    // Request is allowed.
}
```

The first argument is the bucket capacity. The second is the number of
tokens added per second.

## HTTP middleware

```go
limiter := core.NewTokenBucketLimiter(5, 1)
handler := middleware.RateLimit(limiter, middleware.ByIP)(yourHandler)

http.ListenAndServe(":8080", handler)
```

Requests that exceed the limit receive HTTP status `429 Too Many Requests`.

## Redis usage

Start Redis locally on its default port, then create a store:

```go
client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})
redisStore := store.NewRedisStore(client)

allowed, err := redisStore.Allow(
    context.Background(),
    "client-1",
    5,
    10,
)
```

The arguments after the client ID specify the request limit and window
duration in seconds.

## Testing

```bash
go test ./...
```

Run benchmarks with:

```bash
go test -bench=. ./...
```

The Redis benchmark requires a Redis server running at
`localhost:6379`.

## Project structure

```text
core/       In-memory rate limiters
middleware/ HTTP middleware
store/      Redis-backed storage
cmd/demo/   Example HTTP server
```
