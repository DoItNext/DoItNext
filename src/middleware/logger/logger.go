package logger

import (
	"github.com/go-mods/zerolog-rotate/middleware"
	"net/http"
)

//
// zerolog-rotate chi middleware
//
func Chi(next http.Handler) http.Handler {
	return middleware.ChiLogger(next)
}
