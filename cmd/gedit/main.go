package main

import (
	"log"
	"os"

	"github.com/arit-pal/gedit/internal/editor"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: gedit <filename>")
	}
	fileName := os.Args[1]

	ed, err := editor.NewEditor(fileName)
	if err != nil {
		log.Fatalf("Error initializing editor: %v", err)
	}

	ed.Start()
}
