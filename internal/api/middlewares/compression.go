package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	fmt.Println("Compression middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client accepts gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
		}

		// Set the response header
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap the resposneWriter
		w = &gzipresponsewriter{ResponseWriter: w, Writer: gz}

		next.ServeHTTP(w, r)
		fmt.Println("Sent response from compression Middleware")
	})
}

// gzipResponseWriter wraps the standard http.ResponseWriter to route writes through gzip
type gzipresponsewriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipresponsewriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
