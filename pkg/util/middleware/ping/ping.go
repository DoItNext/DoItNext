package ping

import (
	"net/http"
	"strings"
)

//
// ping endpoint middleware
//
// useful as a testing request before hitting any routes
func Ping(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.EqualFold(r.URL.Path, "/ping") {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("pong"))
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
