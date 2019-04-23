package secure

import (
	"github.com/go-chi/cors"
	"net/http"
)

// Headers adds general security headers for basic security measures
// see: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers
func Headers(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Protects from MimeType Sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		// Prevents browser from prefetching DNS
		w.Header().Set("X-DNS-Prefetch-Control", "off")
		// Denies website content to be served in an iframe
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Strict-Transport-Security", "max-age=5184000; includeSubDomains")
		// Prevents Internet Explorer from executing downloads in site's context
		w.Header().Set("X-Download-Options", "noopen")
		// Minimal XSS protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Serve
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

var aa = cors.New(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	ExposedHeaders:   []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           300,
})

// CORS adds Cross-Origin Resource Sharing support
func CORS(h http.Handler) http.Handler {
	return aa.Handler(h)
}
