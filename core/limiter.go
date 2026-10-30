package core

// Limiter decides whether a request from clientID is allowed right now.
type Limiter interface {
	Allow(clientID string) bool
}
