package main

import (
	"log"
	"os"

	"github.com/arit-pal/gedit/app"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: gedit <filename>")
	}
	fileName := os.Args[1]

	// Create a new application instance.
	geditApp, err := app.NewApp(fileName)
	if err != nil {
		log.Fatalf("Error initializing editor: %v", err)
	}

	// Start the application's main loop.
	geditApp.Run()
}
