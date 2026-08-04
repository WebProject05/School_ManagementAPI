package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

// ResponseTimeMiddleware logs the method, status, URL, and duration of each request.
func ResponseTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received Request in Response Writer")
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		duration := time.Since(start)
		wrappedWriter.Header().Set("X-Response-Time", duration.String())
		next.ServeHTTP(wrappedWriter, r)
		
		// Calculate the duration, we are gonna recalculate the time for duration.
		// There might be small changes but it does not matter
		
		duration = time.Since(start)

		fmt.Printf("Method: %s, Status: %d, URL: %s, Duration: %v\n", 
			r.Method, 
			wrappedWriter.status, 
			r.URL.Path, // Using Path is usually cleaner than String() for URLs, but either works!
			duration,
		)
		fmt.Println("Sent Response from Response Time Middleware")
	})
}

// responseWriter wraps http.ResponseWriter to capture the HTTP status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader intercepts the status code before passing it to the real ResponseWriter
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}