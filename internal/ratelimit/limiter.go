package ratelimit

import (
	"sync"
	"time"
)

type bucket struct {
	tokens   float64
	lastTime time.Time
}

// Limiter implements a token bucket rate limiter keyed by string (webhook ID or IP).
type Limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	rate    float64 // tokens per second
	max     float64 // max burst (= tokens per minute from config)
}

// NewLimiter creates a new Limiter with the given requests per minute limit.
func NewLimiter(requestsPerMinute int) *Limiter {
	l := &Limiter{
		buckets: make(map[string]*bucket),
		rate:    float64(requestsPerMinute) / 60.0,
		max:     float64(requestsPerMinute),
	}
	go l.cleanup()
	return l
}

// Allow checks if the given key is allowed to make a request.
// It refills tokens based on elapsed time, then consumes one token if available.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[key]
	if !ok {
		// Lazy initialization: start full
		l.buckets[key] = &bucket{tokens: l.max - 1, lastTime: now}
		return true
	}

	// Refill tokens based on elapsed time
	elapsed := now.Sub(b.lastTime).Seconds()
	b.tokens += elapsed * l.rate
	if b.tokens > l.max {
		b.tokens = l.max
	}
	b.lastTime = now

	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// cleanup runs every 5 minutes and removes buckets that haven't been accessed in > 10 minutes.
func (l *Limiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-10 * time.Minute)
		l.mu.Lock()
		for key, b := range l.buckets {
			if b.lastTime.Before(cutoff) {
				delete(l.buckets, key)
			}
		}
		l.mu.Unlock()
	}
}
