// Package ratelimit is an in-memory fixed-window limiter keyed by string.
// Single-instance only; a shared store is required for horizontal scaling.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter allows Limit events per Window per key.
type Limiter struct {
	Limit  int
	Window time.Duration
	mu     sync.Mutex
	hits   map[string]*bucket
}

type bucket struct {
	start time.Time
	n     int
}

// New creates a limiter.
func New(limit int, window time.Duration) *Limiter {
	return &Limiter{Limit: limit, Window: window, hits: map[string]*bucket{}}
}

// Allow records an event for key and reports whether it is permitted.
func (l *Limiter) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.hits) > 100000 { // bound memory
		for k, b := range l.hits {
			if now.Sub(b.start) >= l.Window {
				delete(l.hits, k)
			}
		}
	}
	b, ok := l.hits[key]
	if !ok || now.Sub(b.start) >= l.Window {
		l.hits[key] = &bucket{start: now, n: 1}
		return true
	}
	b.n++
	return b.n <= l.Limit
}
