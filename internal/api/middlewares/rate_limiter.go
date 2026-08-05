package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]int
	limit     int
	resetTime time.Duration
}

// NewRateLimiter initializes a rateLimiter and starts the background goroutine
// to periodically reset the visitor counts.
func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors:  make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}

	// Launch the reset loop in a background goroutine
	// so it doesn't block the rest of the application.
	go rl.resetVisitorCount()

	return rl
}

// resetVisitorCount runs in the background and clears the visitors map
// at the specified resetTime intervals to refresh their request quotas.
func (rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		
		// Lock the mutex to prevent race conditions (server crashes) if an HTTP
		// request tries to update the map at the exact same time we are clearing it.
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}


func (rl *rateLimiter) Middlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.mu.Lock()

		defer rl.mu.Unlock()

		// Accessing the IP of the visitor
		visitorIp := r.RemoteAddr
		rl.visitors[visitorIp]++
		fmt.Printf("Visitor Count from %v is %v \n", visitorIp, rl.visitors[visitorIp])

		if rl.visitors[visitorIp] > rl.limit {
			http.Error(w, "Too many Requests", http.StatusTooManyRequests)
			return
		}


		next.ServeHTTP(w, r)
	})
}