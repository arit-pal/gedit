package app

import (
	"os"

	"github.com/arit-pal/gedit/internal/editor"
	"github.com/arit-pal/gedit/internal/file"
	"github.com/arit-pal/gedit/internal/input"
	"github.com/arit-pal/gedit/internal/view"
	"github.com/gdamore/tcell/v2"
)

// App is the main application struct.
type App struct {
	state *editor.State
}

// NewApp creates a new application instance.
func NewApp(fileName string) (*App, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}

	screen.EnableMouse()

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	screen.SetStyle(defStyle)

	state := editor.NewState(screen, fileName)
	state.IsDirty = false

	// Load the file content into the state.
	content, err := file.Load(fileName)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}
	state.Content = content

	return &App{state: state}, nil
}

// Run starts the main event loop for the application.
func (a *App) Run() {
	defer a.finish()
	view.UpdateScrolling(a.state)

	for {
		view.Draw(a.state)
		a.state.Screen.Show()

		ev := a.state.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			a.state.Screen.Sync()
		case *tcell.EventKey:
			if quit := input.HandleKeyEvent(ev, a.state); quit {
				return
			}
		case *tcell.EventMouse:
		}
	}
}

// finish cleans up the application.
func (a *App) finish() {
	a.state.Screen.Fini()
}
