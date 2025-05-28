package common

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter limits the rate of API calls
type RateLimiter struct {
	limiters   map[string]*rate.Limiter
	mu         sync.Mutex
	rate       rate.Limit
	burstLimit int
}

// NewRateLimiter creates a new rate limiter
// ratePerSecond: How many requests per second to allow
// burstLimit: How many requests can be made in a burst
func NewRateLimiter(ratePerSecond float64, burstLimit int) *RateLimiter {
	return &RateLimiter{
		limiters:   make(map[string]*rate.Limiter),
		rate:       rate.Limit(ratePerSecond),
		burstLimit: burstLimit,
	}
}

// GetLimiter gets a rate limiter for a particular key (e.g., tenant ID or API endpoint)
func (r *RateLimiter) GetLimiter(key string) *rate.Limiter {
	r.mu.Lock()
	defer r.mu.Unlock()

	limiter, exists := r.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(r.rate, r.burstLimit)
		r.limiters[key] = limiter
	}

	return limiter
}

// Allow checks if a request is allowed to proceed
// Returns true if the request can proceed, false if it exceeds the rate limit
func (r *RateLimiter) Allow(key string) bool {
	return r.GetLimiter(key).Allow()
}

// Reserve reserves a token and returns information about when the reservation can be used
func (r *RateLimiter) Reserve(key string) *rate.Reservation {
	return r.GetLimiter(key).Reserve()
}

// Wait blocks until a token is available
func (r *RateLimiter) Wait(key string) {
	limiter := r.GetLimiter(key)
	limiter.Wait(context.Background())
}
