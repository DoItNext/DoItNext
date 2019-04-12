package server

import (
	"context"
	"fmt"
	"github.com/DoItNext/DoItNext/pkg/util/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-mods/zerolog-rotate/log"
	zrmiddleware "github.com/go-mods/zerolog-rotate/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server configuration
type Server struct {
}

var router = chi.NewRouter()

// Initialize the server with predefined configuration
func (server *Server) Initialize() {
	// Create routes
	server.routes()
}

func (server *Server) routes() {

	// chi middleware
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.RequestID,                          // Inject request ID into the context of each request
		zrmiddleware.ChiLogger,                        // Log API request calls
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
	)

	// Simple ping response
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hi"))
	})
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
	})
}

func (server *Server) Run() {
	// Check router initialisation
	if router == nil {
		log.Panic("Please initialize router before starting the server")
	}

	// server configuration
	cfg := config.Configuration.Server

	// Create the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		Handler:      router,
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 20 seconds.
	idleConnectionClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		// sig is a ^C, handle it
		log.Info("Shutting down..")

		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Error(err, "HTTP server Shutdown: %v")
		}

		close(idleConnectionClosed)
	}()

	// Start server
	log.Info("Started")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Error(err, "HTTP server ListenAndServe: %v")
	}

	<-idleConnectionClosed
}
