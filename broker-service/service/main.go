package main

import (
	"fmt"
	"log"
	"net/http"
)

// webPort is the port on which the HTTP server will listen.
const webPort = "80"

// Config holds the configuration for the application.
type Config struct{}

func main() {
	// Create a new instance of Config.
	app := Config{}

	// Log a message indicating that the broker service is starting.
	log.Printf("Starting broker service on port %s\n", webPort)

	// Create an HTTP server.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(), // Use the Config's routes method as the server's handler.
	}

	// Start the HTTP server.
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
