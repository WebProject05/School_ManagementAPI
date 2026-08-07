package utils

import "net/http"

// Middleware is a function that wraps an http.Handler with addional functionality
type Middleware func(http.Handler) http.Handler

func ApplyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middlware := range middlewares {
		handler = middlware(handler)
	}
	return handler
}
