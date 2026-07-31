package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a concurrency-safe per-key token bucket limiter.
type Limiter struct {
	mu          sync.Mutex
	capacity    int
	refillEvery time.Duration
	now         func() time.Time
	buckets     map[string]*bucket
}

type bucket struct {
	tokens int
	last   time.Time
}

// NewLimiter creates a limiter with a fixed burst capacity and refill interval.
func NewLimiter(capacity int, refillEvery time.Duration) *Limiter {
	return newLimiter(capacity, refillEvery, time.Now)
}

func newLimiter(capacity int, refillEvery time.Duration, now func() time.Time) *Limiter {
	return &Limiter{
		capacity:    capacity,
		refillEvery: refillEvery,
		now:         now,
		buckets:     make(map[string]*bucket),
	}
}

// Allow reports whether the key can consume a token right now.
func (l *Limiter) Allow(key string) bool {
	if l.capacity <= 0 {
		return false
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	current := l.now()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.capacity, last: current}
		l.buckets[key] = b
	}
	l.refill(b, current)
	if b.tokens == 0 {
		return false
	}
	b.tokens--
	return true
}

func (l *Limiter) refill(b *bucket, current time.Time) {
	if l.refillEvery <= 0 || current.Before(b.last) {
		b.last = current
		return
	}
	elapsed := current.Sub(b.last)
	refills := int(elapsed / l.refillEvery)
	if refills == 0 {
		return
	}
	b.tokens += refills
	if b.tokens > l.capacity {
		b.tokens = l.capacity
	}
	b.last = b.last.Add(time.Duration(refills) * l.refillEvery)
}
