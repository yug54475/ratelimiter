package core

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	limit  int
	window time.Duration
	mu     sync.Mutex
	hits   map[string][]time.Time
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:  limit,
		window: window,
		hits:   make(map[string][]time.Time),
	}
}

func (s *SlidingWindowLimiter) Allow(clientID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.window)

	timestamps := s.hits[clientID]
	// drop anything older than the window
	valid := timestamps[:0]
	for _, t := range timestamps {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= s.limit {
		s.hits[clientID] = valid
		return false
	}

	valid = append(valid, now)
	s.hits[clientID] = valid
	return true
}
