package main

import (
	"github.com/DoItNext/DoItNext/pkg/server"
	"github.com/DoItNext/DoItNext/pkg/util/config"
	"github.com/go-mods/zerolog-rotate"
	"github.com/go-mods/zerolog-rotate/log"
	"time"
)

func main() {

	// Create the main logger
	log.Logger = logger.New(logger.Config{
		RwConfig: func(rw *logger.RotateConfig) { rw.LogPath = "logs"; rw.FileName = "server" },
		CwConfig: func(cw *logger.ConsoleConfig) { cw.TimeFormat = time.RFC3339 },
	})

	// Starting the server
	log.Info("DoItNext is starting ...")

	// load server configurations
	if err := config.Load("./config"); err != nil {
		log.Panic(err, "invalid application configuration: %s")
	}

	// Initialize and run the server
	srv := &server.Server{}
	srv.Initialize()
	srv.Run()

	log.Info("DoItNext is ending ...")
}
