package editor

import (
	"github.com/gdamore/tcell/v2"
)

// Editor holds the entire state of the text editor.
type Editor struct {
	screen tcell.Screen
}

// NewEditor creates and initializes a new Editor instance.
func NewEditor() (*Editor, error) {
	// Initialize the screen.
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}

	// Set default text style.
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)

	return &Editor{screen: screen}, nil
}

// Start begins the main event loop for the editor.
func (e *Editor) Start() {
	// The defer call ensures that our Finish function is called when Start() returns.
	// This is crucial for restoring the terminal to its normal state.
	defer e.Finish()

	// Main event loop.
	for {
		// Show the current state of the screen.
		e.screen.Show()

		// Poll for the next event (e.g., key press).
		ev := e.screen.PollEvent()

		// Use a type switch to handle different event types.
		switch ev := ev.(type) {
		case *tcell.EventResize:
			// If the screen is resized, sync to ensure the new size is used.
			e.screen.Sync()
		case *tcell.EventKey:
			// For now, we only care about one key: 'q'.
			// If the 'q' rune is pressed, we break the loop to exit.
			if ev.Rune() == 'q' {
				return // Exit the loop and the function.
			}
		}
	}
}

// Finish cleans up the editor, restoring the terminal.
func (e *Editor) Finish() {
	e.screen.Fini()
}
