package main

import (
	"log"
	// We are referencing our own internal package here
	"github.com/arit-pal/gedit/internal/editor"
)

func main() {
	// Create a new editor instance.
	ed, err := editor.NewEditor()
	if err != nil {
		log.Fatalf("Error initializing editor: %v", err)
	}

	// Start the editor. This will be our main event loop.
	ed.Start()
}
