package store

import "context"

// Store is used by distributed (e.g. Redis-backed) limiters.
// Returns true if the request is allowed.
type Store interface {
	Allow(ctx context.Context, clientID string, limit int, windowSeconds int) (bool, error)
}
