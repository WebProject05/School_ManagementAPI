package middlewares

import (
	"fmt"
	"net/http"
)

// If we have a list of allowed origins
// var allowedOrigins = []string{
// 	"myFrontEnd.com", This is the url of the frontend hosting url where the user will request and it comes to the backend
// 	"",
// }

// Cors handles Cross-Origin Resource Sharing (CORS) configurations
func Cors(next http.Handler) http.Handler {
	fmt.Println("CORS Middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		fmt.Println("Origin:", origin)

		// Check if the request is cross-origin (has an Origin header)
		if origin != "" {
			// Note: Usually local dev uses http://localhost:3000, adjust if needed
			if origin == "https://localhost:3000" || origin == "http://localhost:3000" {
				// CRITICAL: Tells the browser this specific origin is allowed to read the response
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				http.Error(w, "Not allowed by CORS", http.StatusForbidden)
				return
			}
		}

		// Tells the browser which HTTP headers the frontend is allowed to send
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		
		// Tells the browser which HTTP methods are permitted
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		
		// Allows the browser to send cookies or authorization headers with the request
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Intercept preflight OPTIONS requests and stop here
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}