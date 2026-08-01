package middlewares

import "net/http"

/*
Basic middle ware skeleton
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

	})
	// next(w, r)
}

*/

// We are using the http next handler cause if we have multiple
// middlewares we can pass to the next middleware once the current
// middleware has done it's work
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Prevents the browser from pre-resolving DNS for links on the page (Privacy)
        w.Header().Set("X-DNS-Prefetch-Control", "off")

        // Prevents Clickjacking by forbidding the site from being loaded in an iframe
        w.Header().Set("X-Frame-Options", "DENY")

        // Blocks the page from loading if a reflected Cross-Site Scripting (XSS) attack is detected
        w.Header().Set("X-XSS-Protection", "1; mode=block")

        // Forces the browser to trust the provided Content-Type (prevents MIME-sniffing)
        w.Header().Set("X-Content-Type-Options", "nosniff")

        // Forces browsers to exclusively use HTTPS for this site for ~2 years (HSTS)
        w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

        // Restricts resource loading (scripts, images, etc.) to the site's own domain only
        w.Header().Set("Content-Security-Policy", "default-src 'self'")

        // Prevents the browser from sending the URL of this site to other sites when users click external links
        w.Header().Set("Referrer-Policy", "no-referrer")

		// Puts the frame work used in construction of this site
		// w.Header().Set("X-Powered-By", "GoLang")
        
        next.ServeHTTP(w, r)
    })
}