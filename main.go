package main

import (
	"fmt"
	"github.com/DoItNext/DoItNext/src/config"
	"github.com/DoItNext/DoItNext/src/server"
	"github.com/go-mods/zerolog-rotate"
	"github.com/go-mods/zerolog-rotate/log"
	"github.com/iancoleman/strcase"
	"time"
)

func main() {

	// Create the main logger
	log.Logger = logger.New(logger.Config{
		RwConfig: func(rw *logger.ZrRotateConfig) {
			rw.LogPath = "logs"
			rw.FileName = "server"
			rw.TimeTagFormat = time.RFC3339
		},
		CwConfig: func(cw *logger.ZrConsoleConfig) {
			cw.TimeFormat = time.RFC3339
			cw.NoColor = true
			cw.FormatLevel = func(i interface{}) string {
				return fmt.Sprintf("| %-6s|", strcase.ToCamel(i.(string)))
			}
		},
	})

	// Starting the server
	log.Info("Starting ...")

	// load server configurations
	if err := config.Load("./config"); err != nil {
		log.Panic(err, "invalid application configuration: %s")
	}

	// Initialize and run the server
	srv := &server.Server{}
	srv.Initialize()
	srv.Run()

	// That's the end of DoItNext
	log.Info("Stopped")
}
