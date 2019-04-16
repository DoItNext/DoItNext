package server

import (
	"context"
	"fmt"
	"github.com/DoItNext/DoItNext/pkg/util/config"
	"github.com/DoItNext/DoItNext/pkg/util/database/mysql"
	"github.com/DoItNext/DoItNext/pkg/util/middleware/logger"
	"github.com/DoItNext/DoItNext/pkg/util/middleware/ping"
	"github.com/DoItNext/DoItNext/pkg/util/middleware/secure"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-mods/zerolog-rotate/log"
	"github.com/jinzhu/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server configuration
type Server struct {
	DB *gorm.DB
}

var router = chi.NewRouter()

// Initialize the server with predefined configuration
func (server *Server) Initialize() {
	// Create routes
	server.routes()
}

// Create api routes
func (server *Server) routes() {
	// database configuration
	cfg := config.Configuration.Database

	// Establish database connection
	switch cfg.Type {
	case "mysql":
		if db, err := mysql.New(); err != nil {
			log.Error(err, "While trying to initialise the database")
			log.Fatal("Could not connect database")
		} else {
			server.DB = db
		}
	default:
		log.Panic("Only mysql database is supported for the moment")
	}

	// chi middleware
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.RequestID,                          // Inject request ID into the context of each request
		middleware.DefaultCompress,                    // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes,                    // Redirect slashes to no slash URL versions
		middleware.Recoverer,                          // Recover from panics without crashing server
		logger.Chi,                                    // Log API request calls
		secure.Headers,                                // General security headers for basic security measures
		secure.CORS,                                   // Cross-Origin Resource Sharing support
		ping.Ping,                                     // Ping heartbeat
	)

	// Simple ping response
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hi"))
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

	// Wait for interrupt signal to gracefully shutdown the server
	// with a timeout of 20 seconds.
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
